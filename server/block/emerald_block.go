package block

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world/sound"
)

// EmeraldBlock is a precious mineral block crafted using 9 emeralds.
type EmeraldBlock struct {
	solid
}

// Instrument ...
func (e EmeraldBlock) Instrument() sound.Instrument {
	return sound.Bit()
}

// BreakInfo ...
func (e EmeraldBlock) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, oneOf(e)).withExplosionInfo(6, false)
}

// PowersBeacon ...
func (EmeraldBlock) PowersBeacon() bool {
	return true
}

// EncodeItem ...
func (EmeraldBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:emerald_block", 0
}

// EncodeBlock ...
func (EmeraldBlock) EncodeBlock() (string, map[string]any) {
	return "minecraft:emerald_block", nil
}
