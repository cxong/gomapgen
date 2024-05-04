package gmgmap

import (
	"math"
	"math/rand"
)

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

	superpositions := newSuperpositions(width, height)

	rules := []Rule{}

	// Grass with road surroundings
	rules = append(rules, defaultGrassRule)
	rules = append(rules, roadAtEdgeRule)

	// Collapse waves until everything is collapsed
	for {
		// Apply rules on every tile to set up the waves
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				for _, rule := range rules {
					newSP := rule(superpositions, x, y)
					if newSP != nil {
						superpositions.set(x, y, newSP)
					}
				}
			}
		}

		// Choose the tile with lowest entropy
		minEntropy := math.Inf(0)
		minX := 0
		minY := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				entropy := superpositions.get(x, y).entropy()
				if entropy < minEntropy {
					minEntropy = entropy
					minX = x
					minY = y
				}
			}
		}
		if minEntropy == 0.0 {
			// Everything is collapsed
			break
		}
		// Collapse the highest entropy tile
		// Just select the highest weight
		// TODO: investigate other methods of collapse
		var maxKey rune
		var maxWeight float64
		for key, value := range superpositions.get(minX, minY) {
			if value > maxWeight {
				maxWeight = value
				maxKey = key
			}
		}
		superpositions.set(minX, minY, Superposition{maxKey: 1.0})
	}

	g := m.Layer("Ground")
	// s := m.Layer("Structures")
	// f := m.Layer("Furniture")

	// Construct map based on superpositions
	// At this point everything should be collapsed
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			sp := superpositions.get(x, y)
			switch sp.collapsedValue() {
			case grass:
				g.setTile(x, y, grass)
			case road:
				g.setTile(x, y, road)
			}
		}
	}

	return m
}

type Superposition map[rune]float64

func (s Superposition) entropy() float64 {
	sum := float64(0)
	for _, value := range s {
		sum += value * math.Log(value)
	}
	return sum
}

func (s Superposition) collapsedValue() rune {
	for key, value := range s {
		if value == 1.0 {
			return key
		}
	}
	return nothing
}

type Superpositions struct {
	sp     []Superposition
	Width  int
	Height int
}

func (s Superpositions) get(x, y int) Superposition {
	return s.sp[x+y*s.Width]
}

func (s *Superpositions) set(x, y int, sp Superposition) {
	s.sp[x+y*s.Width] = sp
}

func newSuperpositions(width, height int) *Superpositions {
	s := new(Superpositions)
	s.Width, s.Height = width, height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			s.sp = append(s.sp, Superposition{nothing: 1.0})
		}
	}
	return s
}

type Rule func(s *Superpositions, x, y int) Superposition

func defaultGrassRule(s *Superpositions, x, y int) Superposition {
	sp := s.get(x, y)
	if len(sp) == 0 || sp.collapsedValue() == nothing {
		return Superposition{grass: 1.0}
	}
	return nil
}

func roadAtEdgeRule(s *Superpositions, x, y int) Superposition {
	if x == 0 || y == 0 || x == s.Width-1 || y == s.Height-1 {
		return Superposition{road: 1.0}
	}
	return nil
}
