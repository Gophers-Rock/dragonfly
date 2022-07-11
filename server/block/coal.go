package block

import (
	"github.com/df-mc/dragonfly/server/item"
)

// Coal is a precious mineral block made from 9 coal.
type Coal struct {
	solid
	bassDrum
}

// FlammabilityInfo ...
func (c Coal) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(5, 5, false)
}

// BreakInfo ...
func (c Coal) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierWood.HarvestLevel
	}, pickaxeEffective, oneOf(c))
}

// EncodeItem ...
func (Coal) EncodeItem() (name string, meta int16) {
	return "minecraft:coal_block", 0
}

// EncodeBlock ...
func (Coal) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:coal_block", nil
}