package item

type OceanMap struct {
	baseMap
}

// EncodeItem ...
func (m OceanMap) EncodeItem() (name string, meta int16) {
	return "minecraft:filled_map", 3
}