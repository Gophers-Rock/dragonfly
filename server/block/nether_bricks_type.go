package block

import "fmt"

// NetherBricksType represents a type of nether bricks.
type NetherBricksType struct {
	netherBricks
}

type netherBricks uint8

// NormalNetherBricks is the normal variant of the nether bricks.
func NormalNetherBricks() NetherBricksType {
	return NetherBricksType{netherBricks(0)}
}

// RedNetherBricks is the red variant of the nether bricks.
func RedNetherBricks() NetherBricksType {
	return NetherBricksType{netherBricks(1)}
}

// CrackedNetherBricks is the cracked variant of the nether bricks.
func CrackedNetherBricks() NetherBricksType {
	return NetherBricksType{netherBricks(2)}
}

// ChiseledNetherBricks is the chiseled variant of the nether bricks.
func ChiseledNetherBricks() NetherBricksType {
	return NetherBricksType{netherBricks(3)}
}

// Uint8 returns the nether bricks as a unit8
func (n netherBricks) Uint8() uint8 {
	return uint8(n)
}

// Name ...
func (n netherBricks) Name() string {
	switch n {
	case 0:
		return "Nether Bricks"
	case 1:
		return "Red Nether Bricks"
	case 2:
		return "Cracked Nether Bricks"
	case 3:
		return "Chiseled Nether Bricks"
	}
	panic("unknown nether brick type")
}

// FromString ...
func (n netherBricks) FromString(str string) (interface{}, error) {
	switch str {
	case "normal", "default":
		return NormalNetherBricks(), nil
	case "red":
		return RedNetherBricks(), nil
	case "cracked":
		return CrackedNetherBricks(), nil
	case "chiseled":
		return ChiseledNetherBricks(), nil
	}
	return nil, fmt.Errorf("unexpected nether brick type '%v'", n)
}

// String ...
func (n netherBricks) String() string {
	switch n {
	case 0:
		return "nether_brick"
	case 1:
		return "red_nether_brick"
	case 2:
		return "cracked_nether_bricks"
	case 3:
		return "chiseled_nether_bricks"
	}
	panic("unknown nether brick type")
}

// NetherBricksTypes returns all nether bricks type.
func NetherBricksTypes() []NetherBricksType {
	return []NetherBricksType{NormalNetherBricks(), RedNetherBricks(), CrackedNetherBricks(), ChiseledNetherBricks()}
}
