// Code generated by cmd/blockhash; DO NOT EDIT.

package block

const hashAir = 0
const hashAncientDebris = 1
const hashAndesite = 2
const hashBarrier = 3
const hashBasalt = 4
const hashBeacon = 5
const hashBedrock = 6
const hashBeetrootSeeds = 7
const hashBlueIce = 8
const hashBoneBlock = 9
const hashBricks = 10
const hashCake = 11
const hashCarpet = 12
const hashCarrot = 13
const hashChest = 14
const hashChiseledQuartz = 15
const hashClay = 16
const hashCoalBlock = 17
const hashCoalOre = 18
const hashCobblestone = 19
const hashCocoaBean = 20
const hashConcrete = 21
const hashConcretePowder = 22
const hashCoral = 23
const hashCoralBlock = 24
const hashDiamondBlock = 25
const hashDiamondOre = 26
const hashDiorite = 27
const hashDirt = 28
const hashDirtPath = 29
const hashDragonEgg = 30
const hashEmeraldBlock = 31
const hashEmeraldOre = 32
const hashEndBrickStairs = 33
const hashEndBricks = 34
const hashEndStone = 35
const hashFarmland = 36
const hashFire = 37
const hashGildedBlackstone = 38
const hashGlass = 39
const hashGlassPane = 40
const hashGlazedTerracotta = 41
const hashGlowstone = 42
const hashGoldBlock = 43
const hashGoldOre = 44
const hashGranite = 45
const hashGrass = 46
const hashGrassPlant = 47
const hashGravel = 48
const hashInvisibleBedrock = 49
const hashIronBars = 50
const hashIronBlock = 51
const hashIronOre = 52
const hashKelp = 53
const hashLantern = 54
const hashLapisBlock = 55
const hashLapisOre = 56
const hashLava = 57
const hashLeaves = 58
const hashLight = 59
const hashLitPumpkin = 60
const hashLog = 61
const hashMelon = 62
const hashMelonSeeds = 63
const hashNetherBrickFence = 64
const hashNetherGoldOre = 65
const hashNetherQuartzOre = 66
const hashNetherWart = 67
const hashNetheriteBlock = 68
const hashNetherrack = 69
const hashNoteBlock = 70
const hashObsidian = 71
const hashPlanks = 72
const hashPotato = 73
const hashPumpkin = 74
const hashPumpkinSeeds = 75
const hashQuartz = 76
const hashQuartzBricks = 77
const hashQuartzPillar = 78
const hashSand = 79
const hashSandstone = 80
const hashSeaLantern = 81
const hashShroomlight = 82
const hashSoulSand = 83
const hashSoulSoil = 84
const hashSponge = 85
const hashStainedGlass = 86
const hashStainedGlassPane = 87
const hashStainedTerracotta = 88
const hashStone = 89
const hashTerracotta = 90
const hashTorch = 91
const hashWater = 92
const hashWheatSeeds = 93
const hashWoodDoor = 94
const hashWoodFence = 95
const hashWoodFenceGate = 96
const hashWoodSign = 97
const hashWoodSlab = 98
const hashWoodStairs = 99
const hashWoodTrapdoor = 100
const hashWool = 101

func (Air) Hash() uint64 {
	return hashAir
}

func (AncientDebris) Hash() uint64 {
	return hashAncientDebris
}

func (a Andesite) Hash() uint64 {
	return hashAndesite | uint64(boolByte(a.Polished))<<7
}

func (Barrier) Hash() uint64 {
	return hashBarrier
}

func (b Basalt) Hash() uint64 {
	return hashBasalt | uint64(boolByte(b.Polished))<<7 | uint64(b.Axis)<<8
}

func (Beacon) Hash() uint64 {
	return hashBeacon
}

func (b Bedrock) Hash() uint64 {
	return hashBedrock | uint64(boolByte(b.InfiniteBurning))<<7
}

func (b BeetrootSeeds) Hash() uint64 {
	return hashBeetrootSeeds | uint64(b.Growth)<<7
}

func (BlueIce) Hash() uint64 {
	return hashBlueIce
}

func (b BoneBlock) Hash() uint64 {
	return hashBoneBlock | uint64(b.Axis)<<7
}

func (Bricks) Hash() uint64 {
	return hashBricks
}

func (c Cake) Hash() uint64 {
	return hashCake | uint64(c.Bites)<<7
}

func (c Carpet) Hash() uint64 {
	return hashCarpet | uint64(c.Colour.Uint8())<<7
}

func (c Carrot) Hash() uint64 {
	return hashCarrot | uint64(c.Growth)<<7
}

func (c Chest) Hash() uint64 {
	return hashChest | uint64(c.Facing)<<7
}

func (ChiseledQuartz) Hash() uint64 {
	return hashChiseledQuartz
}

func (c Clay) Hash() uint64 {
	return hashClay
}

func (CoalBlock) Hash() uint64 {
	return hashCoalBlock
}

func (CoalOre) Hash() uint64 {
	return hashCoalOre
}

func (c Cobblestone) Hash() uint64 {
	return hashCobblestone | uint64(boolByte(c.Mossy))<<7
}

func (c CocoaBean) Hash() uint64 {
	return hashCocoaBean | uint64(c.Facing)<<7 | uint64(c.Age)<<9
}

func (c Concrete) Hash() uint64 {
	return hashConcrete | uint64(c.Colour.Uint8())<<7
}

func (c ConcretePowder) Hash() uint64 {
	return hashConcretePowder | uint64(c.Colour.Uint8())<<7
}

func (c Coral) Hash() uint64 {
	return hashCoral | uint64(c.Type.Uint8())<<7 | uint64(boolByte(c.Dead))<<11
}

func (c CoralBlock) Hash() uint64 {
	return hashCoralBlock | uint64(c.Type.Uint8())<<7 | uint64(boolByte(c.Dead))<<11
}

func (DiamondBlock) Hash() uint64 {
	return hashDiamondBlock
}

func (DiamondOre) Hash() uint64 {
	return hashDiamondOre
}

func (d Diorite) Hash() uint64 {
	return hashDiorite | uint64(boolByte(d.Polished))<<7
}

func (d Dirt) Hash() uint64 {
	return hashDirt | uint64(boolByte(d.Coarse))<<7
}

func (DirtPath) Hash() uint64 {
	return hashDirtPath
}

func (DragonEgg) Hash() uint64 {
	return hashDragonEgg
}

func (EmeraldBlock) Hash() uint64 {
	return hashEmeraldBlock
}

func (EmeraldOre) Hash() uint64 {
	return hashEmeraldOre
}

func (s EndBrickStairs) Hash() uint64 {
	return hashEndBrickStairs | uint64(boolByte(s.UpsideDown))<<7 | uint64(s.Facing)<<8
}

func (EndBricks) Hash() uint64 {
	return hashEndBricks
}

func (EndStone) Hash() uint64 {
	return hashEndStone
}

func (f Farmland) Hash() uint64 {
	return hashFarmland | uint64(f.Hydration)<<7
}

func (f Fire) Hash() uint64 {
	return hashFire | uint64(f.Type.Uint8())<<7 | uint64(f.Age)<<11
}

func (GildedBlackstone) Hash() uint64 {
	return hashGildedBlackstone
}

func (Glass) Hash() uint64 {
	return hashGlass
}

func (GlassPane) Hash() uint64 {
	return hashGlassPane
}

func (t GlazedTerracotta) Hash() uint64 {
	return hashGlazedTerracotta | uint64(t.Colour.Uint8())<<7 | uint64(t.Facing)<<11
}

func (Glowstone) Hash() uint64 {
	return hashGlowstone
}

func (GoldBlock) Hash() uint64 {
	return hashGoldBlock
}

func (GoldOre) Hash() uint64 {
	return hashGoldOre
}

func (g Granite) Hash() uint64 {
	return hashGranite | uint64(boolByte(g.Polished))<<7
}

func (Grass) Hash() uint64 {
	return hashGrass
}

func (g GrassPlant) Hash() uint64 {
	return hashGrassPlant | uint64(boolByte(g.UpperPart))<<7 | uint64(g.Type.Uint8())<<8
}

func (Gravel) Hash() uint64 {
	return hashGravel
}

func (InvisibleBedrock) Hash() uint64 {
	return hashInvisibleBedrock
}

func (IronBars) Hash() uint64 {
	return hashIronBars
}

func (IronBlock) Hash() uint64 {
	return hashIronBlock
}

func (IronOre) Hash() uint64 {
	return hashIronOre
}

func (k Kelp) Hash() uint64 {
	return hashKelp | uint64(k.Age)<<7
}

func (l Lantern) Hash() uint64 {
	return hashLantern | uint64(boolByte(l.Hanging))<<7 | uint64(l.Type.Uint8())<<8
}

func (LapisBlock) Hash() uint64 {
	return hashLapisBlock
}

func (LapisOre) Hash() uint64 {
	return hashLapisOre
}

func (l Lava) Hash() uint64 {
	return hashLava | uint64(boolByte(l.Still))<<7 | uint64(l.Depth)<<8 | uint64(boolByte(l.Falling))<<16
}

func (l Leaves) Hash() uint64 {
	return hashLeaves | uint64(l.Wood.Uint8())<<7 | uint64(boolByte(l.Persistent))<<11 | uint64(boolByte(l.ShouldUpdate))<<12
}

func (l Light) Hash() uint64 {
	return hashLight | uint64(l.Level)<<7
}

func (l LitPumpkin) Hash() uint64 {
	return hashLitPumpkin | uint64(l.Facing)<<7
}

func (l Log) Hash() uint64 {
	return hashLog | uint64(l.Wood.Uint8())<<7 | uint64(boolByte(l.Stripped))<<11 | uint64(l.Axis)<<12
}

func (Melon) Hash() uint64 {
	return hashMelon
}

func (m MelonSeeds) Hash() uint64 {
	return hashMelonSeeds | uint64(m.Growth)<<7 | uint64(m.Direction)<<15
}

func (NetherBrickFence) Hash() uint64 {
	return hashNetherBrickFence
}

func (NetherGoldOre) Hash() uint64 {
	return hashNetherGoldOre
}

func (NetherQuartzOre) Hash() uint64 {
	return hashNetherQuartzOre
}

func (n NetherWart) Hash() uint64 {
	return hashNetherWart | uint64(n.Age)<<7
}

func (NetheriteBlock) Hash() uint64 {
	return hashNetheriteBlock
}

func (Netherrack) Hash() uint64 {
	return hashNetherrack
}

func (n NoteBlock) Hash() uint64 {
	return hashNoteBlock
}

func (o Obsidian) Hash() uint64 {
	return hashObsidian | uint64(boolByte(o.Crying))<<7
}

func (p Planks) Hash() uint64 {
	return hashPlanks | uint64(p.Wood.Uint8())<<7
}

func (p Potato) Hash() uint64 {
	return hashPotato | uint64(p.Growth)<<7
}

func (p Pumpkin) Hash() uint64 {
	return hashPumpkin | uint64(boolByte(p.Carved))<<7 | uint64(p.Facing)<<8
}

func (p PumpkinSeeds) Hash() uint64 {
	return hashPumpkinSeeds | uint64(p.Growth)<<7 | uint64(p.Direction)<<15
}

func (q Quartz) Hash() uint64 {
	return hashQuartz | uint64(boolByte(q.Smooth))<<7
}

func (QuartzBricks) Hash() uint64 {
	return hashQuartzBricks
}

func (q QuartzPillar) Hash() uint64 {
	return hashQuartzPillar | uint64(q.Axis)<<7
}

func (s Sand) Hash() uint64 {
	return hashSand | uint64(boolByte(s.Red))<<7
}

func (s Sandstone) Hash() uint64 {
	return hashSandstone | uint64(s.Type.Uint8())<<7 | uint64(boolByte(s.Red))<<11
}

func (SeaLantern) Hash() uint64 {
	return hashSeaLantern
}

func (Shroomlight) Hash() uint64 {
	return hashShroomlight
}

func (SoulSand) Hash() uint64 {
	return hashSoulSand
}

func (SoulSoil) Hash() uint64 {
	return hashSoulSoil
}

func (s Sponge) Hash() uint64 {
	return hashSponge | uint64(boolByte(s.Wet))<<7
}

func (g StainedGlass) Hash() uint64 {
	return hashStainedGlass | uint64(g.Colour.Uint8())<<7
}

func (p StainedGlassPane) Hash() uint64 {
	return hashStainedGlassPane | uint64(p.Colour.Uint8())<<7
}

func (t StainedTerracotta) Hash() uint64 {
	return hashStainedTerracotta | uint64(t.Colour.Uint8())<<7
}

func (s Stone) Hash() uint64 {
	return hashStone | uint64(boolByte(s.Smooth))<<7
}

func (Terracotta) Hash() uint64 {
	return hashTerracotta
}

func (t Torch) Hash() uint64 {
	return hashTorch | uint64(t.Facing)<<7 | uint64(t.Type.Uint8())<<10
}

func (w Water) Hash() uint64 {
	return hashWater | uint64(boolByte(w.Still))<<7 | uint64(w.Depth)<<8 | uint64(boolByte(w.Falling))<<16
}

func (s WheatSeeds) Hash() uint64 {
	return hashWheatSeeds | uint64(s.Growth)<<7
}

func (d WoodDoor) Hash() uint64 {
	return hashWoodDoor | uint64(d.Wood.Uint8())<<7 | uint64(d.Facing)<<11 | uint64(boolByte(d.Open))<<13 | uint64(boolByte(d.Top))<<14 | uint64(boolByte(d.Right))<<15
}

func (w WoodFence) Hash() uint64 {
	return hashWoodFence | uint64(w.Wood.Uint8())<<7
}

func (f WoodFenceGate) Hash() uint64 {
	return hashWoodFenceGate | uint64(f.Wood.Uint8())<<7 | uint64(f.Facing)<<11 | uint64(boolByte(f.Open))<<13 | uint64(boolByte(f.Lowered))<<14
}

func (s WoodSign) Hash() uint64 {
	return hashWoodSign | uint64(s.Wood.Uint8())<<7 | uint64(s.Facing)<<11 | uint64(boolByte(s.Standing))<<19
}

func (s WoodSlab) Hash() uint64 {
	return hashWoodSlab | uint64(s.Wood.Uint8())<<7 | uint64(boolByte(s.Top))<<11 | uint64(boolByte(s.Double))<<12
}

func (s WoodStairs) Hash() uint64 {
	return hashWoodStairs | uint64(s.Wood.Uint8())<<7 | uint64(boolByte(s.UpsideDown))<<11 | uint64(s.Facing)<<12
}

func (t WoodTrapdoor) Hash() uint64 {
	return hashWoodTrapdoor | uint64(t.Wood.Uint8())<<7 | uint64(t.Facing)<<11 | uint64(boolByte(t.Open))<<13 | uint64(boolByte(t.Top))<<14
}

func (w Wool) Hash() uint64 {
	return hashWool | uint64(w.Colour.Uint8())<<7
}
