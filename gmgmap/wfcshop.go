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
//   - at least 3 wide from the right wall
//   - does not overlap with rest area
//   - items on shelves
//
// - rest area (if at least 3 wide free)
//   - below counter area
//   - between the left wall and middle of the room
//   - does not overlap with shelf area
//   - rug
//   - tables/chairs (randomly placed, at least one chair per table)
//
// - against walls (display items, pots, barrels, leave diagonals free)
// - assistants (1 per 100 tiles, after the first)
// - patrons (1 per 36 tiles)
func NewWFCShop(rr *rand.Rand, exportFunc func(*Map), width, height int) *Map {
	m := NewMap(width, height)

	exportFunc(m)

	superpositions := []map[rune]float32{}

	// Grass with road surroundings
	g := m.Layer("Ground")
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			superpositions = append(superpositions, map[rune]float32{grass: 1.0})
		}
	}
	g.fill(grass)

	// s := m.Layer("Structures")
	// f := m.Layer("Furniture")

	// Construct map based on superpositions
	// At this point everything should be collapsed
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			sp := superpositions[x+y*g.Width]
			if len(sp) != 1 {
				panic("a superposition is not collapsed")
			}
			for key, value := range sp {
				if value == 1.0 {
					switch key {
					case grass:
						g.setTile(x, y, key)
					}
					break
				}
			}
		}
	}

	return m
}
