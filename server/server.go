package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/internal/blockinternal"
	"github.com/df-mc/dragonfly/server/internal/iteminternal"
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	_ "github.com/df-mc/dragonfly/server/item" // Imported for maintaining correct initialisation order.
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/maps"
	"iter"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Server implements a Dragonfly server. It runs the main server loop and
// handles the connections of players trying to join the server.
type Server struct {
	conf Config

	once    sync.Once
	started atomic.Pointer[time.Time]

	world, nether, end *world.World

	customBlocks []protocol.BlockEntry
	customItems  []protocol.ItemComponentEntry

	listeners []Listener
	incoming  chan incoming

	pmu sync.RWMutex
	// p holds a map of all players currently connected to the server. When they
	// leave, they are removed from the map.
	p map[uuid.UUID]*onlinePlayer
	// pwg is a sync.WaitGroup used to wait for all players to be disconnected
	// before server shutdown, so that their data is saved properly.
	pwg sync.WaitGroup
	// wg is used to wait for all Listeners to be closed and their respective
	// goroutines to be finished.
	wg sync.WaitGroup
}

type incoming struct {
	s *session.Session
	p *onlinePlayer
	w *world.World
}

type onlinePlayer struct {
	handle *world.EntityHandle
	xuid   string
	name   string
}

// HandleFunc is a function that may be passed to Server.Accept(). It can be
// used to prepare the session of a player before it can do anything.
type HandleFunc func(p *player.Player)

// New creates a Server using a default Config. The Server's worlds are created
// and connections from the Server's listeners may be accepted by calling
// Server.Listen() and Server.Accept() afterwards.
func New() *Server {
	var conf Config
	return conf.New()
}

// Listen starts running the server's listeners but does not block, unlike Run.
// Connections will be accepted on a different goroutine until the listeners
// are closed using a call to Close. Once started, players may be accepted
// using Server.Accept().
func (srv *Server) Listen() {
	t := time.Now()
	if !srv.started.CompareAndSwap(nil, &t) {
		panic("start server: already started")
	}

	info, _ := debug.ReadBuildInfo()
	if info == nil {
		info = &debug.BuildInfo{GoVersion: "N/A", Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "N/A"}}}
	}
	revision := ""
	for _, set := range info.Settings {
		if set.Key == "vcs.revision" {
			revision = set.Value
		}
	}

	srv.conf.Log.Info("Starting Dragonfly server...", "mc-version", protocol.CurrentVersion, "go-version", info.GoVersion, "commit", revision)
	srv.startListening()
	go srv.wait()
}

// Accept accepts an incoming player into the server. It blocks until a player
// connects to the server. A HandleFunc may be passed which is run immediately
// before a *player.Player is accepted to the Server. This function may be used
// to add a player.Handler to the player and prepare its session. The function
// may be nil if player joining does not need to be handled. Accept returns
// false if the Server is closed using a call to Close.
func (srv *Server) Accept() iter.Seq[*player.Player] {
	return func(yield func(*player.Player) bool) {
		for {
			inc, ok := <-srv.incoming
			if !ok {
				return
			}
			srv.pmu.Lock()
			srv.p[inc.p.handle.UUID()] = inc.p
			srv.pmu.Unlock()

			ret := false
			<-inc.w.Exec(func(tx *world.Tx) {
				p := tx.AddEntity(inc.p.handle).(*player.Player)
				inc.s.Spawn(p, tx)
				ret = !yield(p)
			})
			if ret {
				return
			}
		}
	}
}

// World returns the overworld of the server. Players will be spawned in this
// world and this world will be read from and written to when the world is
// edited.
func (srv *Server) World() *world.World {
	return srv.world
}

// Nether returns the nether world of the server. Players are transported to it
// when entering a nether portal in the world returned by the World method.
func (srv *Server) Nether() *world.World {
	return srv.nether
}

// End returns the end world of the server. Players are transported to it when
// entering an end portal in the world returned by the World method.
func (srv *Server) End() *world.World {
	return srv.end
}

// MaxPlayerCount returns the maximum amount of players that are allowed to
// play on the server at the same time. Players trying to join when the server
// is full will be refused to enter. If the config has a maximum player count
// set to 0, MaxPlayerCount will return Server.PlayerCount + 1.
func (srv *Server) MaxPlayerCount() int {
	if srv.conf.MaxPlayers == 0 {
		srv.pmu.RLock()
		defer srv.pmu.RUnlock()
		return len(srv.p) + 1
	}
	return srv.conf.MaxPlayers
}

// Players returns a list of all players currently connected to the server.
// Note that the slice returned is not updated when new players join or leave,
// so it is only valid for as long as no new players join or players leave.
func (srv *Server) Players() iter.Seq[*player.Player] {
	srv.pmu.RLock()
	handles := make([]*world.EntityHandle, 0, len(srv.p))
	for _, p := range srv.p {
		handles = append(handles, p.handle)
	}
	srv.pmu.RUnlock()

	return func(yield func(*player.Player) bool) {
		for _, handle := range handles {
			ret := false
			handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
				ret = !yield(e.(*player.Player))
			})
			if ret {
				break
			}
		}
	}
}

// Player looks for a player on the server with the UUID passed. If found, the
// player is returned and the bool returns holds a true value. If not, the bool
// returned is false and the player is nil.
func (srv *Server) Player(uuid uuid.UUID) (*world.EntityHandle, bool) {
	srv.pmu.RLock()
	defer srv.pmu.RUnlock()
	p, ok := srv.p[uuid]
	return p.handle, ok
}

// PlayerByName looks for a player on the server with the name passed. If
// found, the player is returned and the bool returned holds a true value. If
// not, the bool is false and the player is nil
func (srv *Server) PlayerByName(name string) (*world.EntityHandle, bool) {
	if p, ok := sliceutil.SearchValue(maps.Values(srv.p), func(p *onlinePlayer) bool {
		return p.name == name
	}); ok {
		return p.handle, true
	}
	return nil, false
}

// PlayerByXUID looks for a player on the server with the XUID passed. If found,
// the player is returned and the bool returned is true. If no player with the
// XUID was found, nil and false are returned.
func (srv *Server) PlayerByXUID(xuid string) (*world.EntityHandle, bool) {
	if p, ok := sliceutil.SearchValue(maps.Values(srv.p), func(p *onlinePlayer) bool {
		return p.xuid == xuid
	}); ok {
		return p.handle, true
	}
	return nil, false
}

// CloseOnProgramEnd closes the server right before the program ends, so that
// all data of the server are saved properly.
func (srv *Server) CloseOnProgramEnd() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		if err := srv.Close(); err != nil {
			srv.conf.Log.Error("close server: " + err.Error())
		}
	}()
}

// Close closes the server, making any call to Run/Accept cancel immediately.
func (srv *Server) Close() error {
	if srv.started.Load() == nil {
		panic("server not yet running")
	}
	srv.once.Do(srv.close)
	return nil
}

// close stops the server, storing player and world data to disk when
// necessary.
func (srv *Server) close() {
	srv.conf.Log.Info("Server closing...")

	srv.conf.Log.Debug("Disconnecting players...")
	for p := range srv.Players() {
		p.Disconnect(text.Colourf("<yellow>%v</yellow>", srv.conf.ShutdownMessage))
	}
	srv.pwg.Wait()

	srv.conf.Log.Debug("Closing player provider...")
	if err := srv.conf.PlayerProvider.Close(); err != nil {
		srv.conf.Log.Error("Close player provider: " + err.Error())
	}

	srv.conf.Log.Debug("Closing worlds...")
	for _, w := range []*world.World{srv.end, srv.nether, srv.world} {
		if err := w.Close(); err != nil {
			srv.conf.Log.Error(fmt.Sprintf("Close dimension %v: ", w.Dimension()) + err.Error())
		}
	}

	srv.conf.Log.Debug("Closing listeners...")
	for _, l := range srv.listeners {
		if err := l.Close(); err != nil {
			srv.conf.Log.Error("Close listener: " + err.Error())
		}
	}
	srv.conf.Log.Info("Server closed.", "uptime", time.Since(*srv.started.Load()).String())
}

// listen makes the Server listen for new connections from the Listener passed.
// This may be used to listen for players on different interfaces. Note that
// the maximum player count of additional Listeners added is not enforced
// automatically. The limit must be enforced by the Listener.
func (srv *Server) listen(l Listener) {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	for {
		c, err := l.Accept()
		if err != nil {
			// Cancel the context so that any call to StartGameContext is
			// cancelled rapidly.
			cancel()
			// First wait until all connections that are being handled are
			// done inserting the player into the channel. Afterwards, when
			// we're sure no more values will be inserted in the players
			// channel, we can return so the player channel can be closed.
			wg.Wait()
			srv.wg.Done()
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if msg, ok := srv.conf.Allower.Allow(c.RemoteAddr(), c.IdentityData(), c.ClientData()); !ok {
				_ = c.WritePacket(&packet.Disconnect{HideDisconnectionScreen: msg == "", Message: msg})
				_ = c.Close()
				return
			}
			srv.finaliseConn(ctx, c, l)
		}()
	}
}

// startListening starts making the EncodeBlock listener listen, accepting new
// connections from players.
func (srv *Server) startListening() {
	srv.makeBlockEntries()
	srv.makeItemComponents()

	srv.wg.Add(len(srv.conf.Listeners))
	for _, lf := range srv.conf.Listeners {
		l, err := lf(srv.conf)
		if err != nil {
			srv.conf.Log.Error("create listener: " + err.Error())
			return
		}
		srv.listeners = append(srv.listeners, l)
		go srv.listen(l)
	}
}

// makeBlockEntries initializes the server's block components map using the registered custom blocks. It allows block
// components to be created only once at startup.
func (srv *Server) makeBlockEntries() {
	custom := maps.Values(world.CustomBlocks())
	srv.customBlocks = make([]protocol.BlockEntry, len(custom))

	for i, b := range custom {
		name, _ := b.EncodeBlock()
		srv.customBlocks[i] = protocol.BlockEntry{
			Name:       name,
			Properties: blockinternal.Components(name, b, 10000+int32(i)),
		}
	}
}

// makeItemComponents initializes the server's item components map using the
// registered custom items. It allows item components to be created only once
// at startup
func (srv *Server) makeItemComponents() {
	custom := world.CustomItems()
	srv.customItems = make([]protocol.ItemComponentEntry, len(custom))

	for _, it := range custom {
		name, _ := it.EncodeItem()
		srv.customItems = append(srv.customItems, protocol.ItemComponentEntry{
			Name: name,
			Data: iteminternal.Components(it),
		})
	}
}

// wait awaits the closing of all Listeners added to the Server through a call
// to listen and closed the players channel once that happens.
func (srv *Server) wait() {
	srv.wg.Wait()
	close(srv.incoming)
}

// finaliseConn finalises the session.Conn passed and subtracts from the
// sync.WaitGroup once done.
func (srv *Server) finaliseConn(ctx context.Context, conn session.Conn, l Listener) {
	id := uuid.MustParse(conn.IdentityData().Identity)
	data := srv.defaultGameData()

	var playerData *player.Data
	if d, err := srv.conf.PlayerProvider.Load(id, srv.dimension); err == nil {
		if d.World == nil {
			d.World = srv.world
		}
		data.PlayerPosition = vec64To32(d.Position).Add(mgl32.Vec3{0, 1.62})
		dim, _ := world.DimensionID(d.World.Dimension())
		data.Dimension = int32(dim)
		data.Yaw, data.Pitch = float32(d.Yaw), float32(d.Pitch)

		playerData = &d
	}

	if err := conn.StartGameContext(ctx, data); err != nil {
		_ = l.Disconnect(conn, "Connection timeout.")

		srv.conf.Log.Debug("spawn failed: "+err.Error(), "raddr", conn.RemoteAddr())
		return
	}
	if _, ok := srv.Player(id); ok {
		_ = l.Disconnect(conn, "Already logged in.")
		srv.conf.Log.Debug("spawn failed: already logged in", "raddr", conn.RemoteAddr())
		return
	}
	_ = conn.WritePacket(&packet.ItemComponent{Items: srv.customItems})
	srv.incoming <- srv.createPlayer(id, conn, playerData)
}

// defaultGameData returns a minecraft.GameData as sent for a new player. It
// may later be modified if the player was saved in the player provider of the
// server.
func (srv *Server) defaultGameData() minecraft.GameData {
	return minecraft.GameData{
		// We set these IDs to 1, because that's how the session will treat them.
		EntityUniqueID:  1,
		EntityRuntimeID: 1,

		WorldName:       srv.conf.Name,
		BaseGameVersion: protocol.CurrentVersion,

		Time:       int64(srv.world.Time()),
		Difficulty: 2,

		PlayerGameMode:    packet.GameTypeCreative,
		PlayerPermissions: packet.PermissionLevelMember,
		PlayerPosition:    vec64To32(srv.world.Spawn().Vec3Centre().Add(mgl64.Vec3{0, 1.62})),

		Items:        srv.itemEntries(),
		CustomBlocks: srv.customBlocks,
		GameRules:    []protocol.GameRule{{Name: "naturalregeneration", Value: false}},

		ServerAuthoritativeInventory: true,
		PlayerMovementSettings: protocol.PlayerMovementSettings{
			MovementType:                     protocol.PlayerMovementModeServer,
			ServerAuthoritativeBlockBreaking: true,
		},
	}
}

// dimension returns a world by a dimension passed.
func (srv *Server) dimension(dimension world.Dimension) *world.World {
	switch dimension {
	default:
		return srv.world
	case world.Nether:
		return srv.nether
	case world.End:
		return srv.end
	}
}

// checkNetIsolation checks if a loopback exempt is in place to allow the
// hosting device to join the server. This is only relevant on Windows. It will
// never log anything for anything but Windows.
func (srv *Server) checkNetIsolation() {
	if runtime.GOOS != "windows" {
		// Only an issue on Windows.
		return
	}
	data, _ := exec.Command("CheckNetIsolation", "LoopbackExempt", "-s", `-n="microsoft.minecraftuwp_8wekyb3d8bbwe"`).CombinedOutput()
	if bytes.Contains(data, []byte("microsoft.minecraftuwp_8wekyb3d8bbwe")) {
		return
	}
	const loopbackExemptCmd = `CheckNetIsolation LoopbackExempt -a -n="Microsoft.MinecraftUWP_8wekyb3d8bbwe"`
	srv.conf.Log.Info("You are currently unable to join the server on this machine. Run " + loopbackExemptCmd + " in an admin PowerShell session to resolve.")
}

// handleSessionClose handles the closing of a session. It removes the player
// of the session from the server.
func (srv *Server) handleSessionClose(tx *world.Tx, c session.Controllable) {
	srv.pmu.Lock()
	_, ok := srv.p[c.UUID()]
	delete(srv.p, c.UUID())
	srv.pmu.Unlock()
	if !ok {
		// When a player disconnects immediately after a session is started, it might not be added to the players map
		// yet. This is expected, but we need to be careful not to crash when this happens.
		return
	}

	if err := srv.conf.PlayerProvider.Save(c.UUID(), c.(*player.Player).Data()); err != nil {
		srv.conf.Log.Error("Save player data: " + err.Error())
	}
	srv.pwg.Done()
}

// createPlayer creates a new player instance using the UUID and connection
// passed.
func (srv *Server) createPlayer(id uuid.UUID, conn session.Conn, data *player.Data) incoming {
	srv.pwg.Add(1)

	w, gm, pos := srv.world, srv.world.DefaultGameMode(), srv.world.Spawn().Vec3Middle()
	if data != nil {
		w, gm, pos = data.World, data.GameMode, data.Position
	}
	s := session.Config{
		Log:            srv.conf.Log,
		MaxChunkRadius: srv.conf.MaxChunkRadius,
		JoinMessage:    srv.conf.JoinMessage,
		QuitMessage:    srv.conf.QuitMessage,
		HandleStop:     srv.handleSessionClose,
	}.New(conn)

	// TODO: Do something with the gamemode and other player data.
	_ = gm

	conf := player.Config{
		Name:    conn.IdentityData().DisplayName,
		XUID:    conn.IdentityData().XUID,
		UUID:    id,
		Locale:  conn.ClientData().LanguageCode,
		Skin:    srv.parseSkin(conn.ClientData()),
		Data:    data,
		Pos:     pos,
		Session: s,
	}
	handle := world.EntitySpawnOpts{Position: pos, ID: id}.New(player.Type, conf)
	return incoming{s: s, w: w, p: &onlinePlayer{name: conf.Name, xuid: conf.XUID, handle: handle}}
}

// createWorld loads a world of the server with a specific dimension, ending
// the program if the world could not be loaded. The layers passed are used to
// create a generator.Flat that is used as generator for the world.
func (srv *Server) createWorld(dim world.Dimension, nether, end **world.World) *world.World {
	logger := srv.conf.Log.With("dimension", strings.ToLower(fmt.Sprint(dim)))
	logger.Debug("Loading dimension...")

	conf := world.Config{
		Log:             logger,
		Dim:             dim,
		Provider:        srv.conf.WorldProvider,
		Generator:       srv.conf.Generator(dim),
		RandomTickSpeed: srv.conf.RandomTickSpeed,
		ReadOnly:        srv.conf.ReadOnlyWorld,
		Entities:        srv.conf.Entities,
		PortalDestination: func(dim world.Dimension) *world.World {
			if dim == world.Nether {
				return *nether
			} else if dim == world.End {
				return *end
			}
			return nil
		},
	}
	w := conf.New()
	logger.Info("Opened world.", "name", w.Name())
	return w
}

// parseSkin parses a skin from the login.ClientData  and returns it.
func (srv *Server) parseSkin(data login.ClientData) skin.Skin {
	// Gophertunnel guarantees the following values are valid data and are of
	// the correct size.
	skinResourcePatch, _ := base64.StdEncoding.DecodeString(data.SkinResourcePatch)

	playerSkin := skin.New(data.SkinImageWidth, data.SkinImageHeight)
	playerSkin.Persona = data.PersonaSkin
	playerSkin.Pix, _ = base64.StdEncoding.DecodeString(data.SkinData)
	playerSkin.Model, _ = base64.StdEncoding.DecodeString(data.SkinGeometry)
	playerSkin.ModelConfig, _ = skin.DecodeModelConfig(skinResourcePatch)
	playerSkin.PlayFabID = data.PlayFabID

	playerSkin.Cape = skin.NewCape(data.CapeImageWidth, data.CapeImageHeight)
	playerSkin.Cape.Pix, _ = base64.StdEncoding.DecodeString(data.CapeData)

	for _, animation := range data.AnimatedImageData {
		var t skin.AnimationType
		switch animation.Type {
		case protocol.SkinAnimationHead:
			t = skin.AnimationHead
		case protocol.SkinAnimationBody32x32:
			t = skin.AnimationBody32x32
		case protocol.SkinAnimationBody128x128:
			t = skin.AnimationBody128x128
		}

		anim := skin.NewAnimation(animation.ImageWidth, animation.ImageHeight, animation.AnimationExpression, t)
		anim.FrameCount = int(animation.Frames)
		anim.Pix, _ = base64.StdEncoding.DecodeString(animation.Image)

		playerSkin.Animations = append(playerSkin.Animations, anim)
	}

	return playerSkin
}

// registerTargetFunc registers a cmd.TargetFunc to be able to get all players
// connected and all entities in the server's world.
func (srv *Server) registerTargetFunc() {
	cmd.AddTargetFunc(func(src cmd.Source) (entities []cmd.Target, players []cmd.NamedTarget) {
		// TODO: Figure out transactions with commands.
		// return sliceutil.Convert[cmd.Target](src.World().Entities()), sliceutil.Convert[cmd.NamedTarget](srv.Players())
		return nil, nil
	})
}

// vec64To32 converts a mgl64.Vec3 to a mgl32.Vec3.
func vec64To32(vec3 mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(vec3[0]), float32(vec3[1]), float32(vec3[2])}
}

// itemEntries loads a list of all custom item entries of the server, ready to
// be sent in the StartGame packet.
func (srv *Server) itemEntries() []protocol.ItemEntry {
	entries := make([]protocol.ItemEntry, 0, len(itemRuntimeIDs))

	for name, rid := range itemRuntimeIDs {
		entries = append(entries, protocol.ItemEntry{
			Name:      name,
			RuntimeID: int16(rid),
		})
	}
	for _, it := range world.CustomItems() {
		name, _ := it.EncodeItem()
		rid, _, _ := world.ItemRuntimeID(it)
		entries = append(entries, protocol.ItemEntry{
			Name:           name,
			ComponentBased: true,
			RuntimeID:      int16(rid),
		})
	}
	return entries
}

// ashyBiome represents a biome that has any form of ash.
type ashyBiome interface {
	// Ash returns the ash and white ash of the biome.
	Ash() (ash float64, whiteAsh float64)
}

// sporingBiome represents a biome that has blue or red spores.
type sporingBiome interface {
	// Spores returns the blue and red spores of the biome.
	Spores() (blueSpores float64, redSpores float64)
}

// biomes builds a mapping of all biome definitions of the server, ready to be set in the biomes field of the server
// listener.
func biomes() map[string]any {
	definitions := make(map[string]any)
	for _, b := range world.Biomes() {
		definition := map[string]any{
			"name_hash":   b.String(), // This isn't actually a hash despite what the field name may suggest.
			"temperature": float32(b.Temperature()),
			"downfall":    float32(b.Rainfall()),
			"rain":        b.Rainfall() > 0,
		}
		if a, ok := b.(ashyBiome); ok {
			ash, whiteAsh := a.Ash()
			definition["ash"], definition["white_ash"] = float32(ash), float32(whiteAsh)
		}
		if s, ok := b.(sporingBiome); ok {
			blueSpores, redSpores := s.Spores()
			definition["blue_spores"], definition["red_spores"] = float32(blueSpores), float32(redSpores)
		}
		definitions[b.String()] = definition
	}
	return definitions
}

var (
	//go:embed world/item_runtime_ids.nbt
	itemRuntimeIDData []byte
	itemRuntimeIDs    = map[string]int32{}
)

// init reads all item entries from the resource JSON, and sets the according
// values in the runtime ID maps. init also seeds the global `rand` with the
// current time.
func init() {
	_ = nbt.Unmarshal(itemRuntimeIDData, &itemRuntimeIDs)
}
