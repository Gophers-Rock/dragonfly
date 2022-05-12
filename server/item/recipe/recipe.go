package recipe

import "github.com/df-mc/dragonfly/server/item"

// Recipe is implemented by all recipe types.
type Recipe interface {
	// Input returns the items required to craft the recipe.
	Input() []item.Stack
	// Output returns the items that are produced when the recipe is crafted.
	Output() []item.Stack
	// Block returns the block that is used to craft the recipe.
	Block() string
	// Priority returns the priority of the recipe. Recipes with lower priority are preferred compared to recipes with
	// higher priority.
	Priority() int
}

// ShapelessRecipe is a recipe that has no particular shape.
type ShapelessRecipe struct {
	recipe
}

// NewShapelessRecipe creates a new shapeless recipe and returns it.
func NewShapelessRecipe(input []item.Stack, output []item.Stack, block string, priority int) ShapelessRecipe {
	return ShapelessRecipe{recipe: recipe{
		input:    input,
		output:   output,
		block:    block,
		priority: priority,
	}}
}

// ShapedRecipe is a recipe that has a specific shape that must be used to craft the output of the recipe.
type ShapedRecipe struct {
	recipe
	// Shape contains the width and height of the shaped recipe.
	Shape Shape
}

// NewShapedRecipe creates a new shaped recipe and returns it.
func NewShapedRecipe(input []item.Stack, output []item.Stack, block string, priority int, shape Shape) ShapedRecipe {
	return ShapedRecipe{
		Shape: shape,
		recipe: recipe{
			input:    input,
			output:   output,
			block:    block,
			priority: priority,
		},
	}
}

// recipe implements the Recipe interface. Structs in this package may embed it to gets its functionality
// out of the box.
type recipe struct {
	// input is a list of items that serve as the input of the shaped recipe. These items are the items
	// required to craft the output. The amount of input items must be exactly equal to Width * Height.
	input []item.Stack
	// output contains items that are created as a result of crafting the recipe.
	output []item.Stack
	// block is the block that is used to craft the recipe.
	block string
	// priority is the priority of the recipe versus others.
	priority int
}

// Input ...
func (r recipe) Input() []item.Stack {
	return r.input
}

// Output ...
func (r recipe) Output() []item.Stack {
	return r.output
}

// Block ...
func (r recipe) Block() string {
	return r.block
}

// Priority ...
func (r recipe) Priority() int {
	return r.priority
}
