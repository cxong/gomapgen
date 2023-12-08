package gmgmap

import "math/rand"

// WFCShop - create a single shop, surrounded by road tiles.
// The interior is filled using Wave Function Collapse, with these rules:
// - Road on the edge
// - Walls inside road, except for the bottom edge
//   - On top wall: windows/candles/shelves
//   - On bottom wall:
//   - door (at least 1 tile away from corner)
//   - shop sign next to door
//
// - Grass between the road and bottom edge wall
//   - On grass: notice sign besides door-road, not below shop sign
//
// - Floor inside the walls
// - Counter area 1 tile below
//
// - Counter (at least 2 wide)
//   - shopkeep
//
// - shelves (up to 3 width aisles)
//   - items on shelves
//
// - rest area (if at least 3 wide free)
//   - rug
//   - tables/chairs (randomly placed, at least one chair per table)
//
// - against walls (display items, pots, barrels, leave diagonals free)
// - assistants (1 per 100 tiles, after the first)
// - patrons (1 per 36 tiles)
func NewWFCShop(rr *rand.Rand, width, height int) *Map {
	m := NewMap(width, height)

	// Grass with road surroundings
	g := m.Layer("Ground")
	g.fill(grass)

	// s := m.Layer("Structures")
	// f := m.Layer("Furniture")

	return m
}
