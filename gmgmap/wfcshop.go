package gmgmap

import (
	"fmt"
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

	g := m.Layer("Ground")
	s := m.Layer("Structures")
	f := m.Layer("Furniture")

	superpositions := newSuperpositions(m)

	rules := []Rule{}

	// Grass with road surroundings
	defaultTile := grass
	rules = append(rules, roadAtEdgeRule)

	// Walls inside road
	rules = append(rules, wallsInsideRoadRule)

	// Collapse waves until everything is collapsed
	collapseCounter := 0
	cd := 16
	for {
		autoCollapsed := false
		// Apply rules on every uncollapsed tile to set up the waves
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if superpositions.get(x, y).isCollapsed() {
					continue
				}
				for _, rule := range rules {
					newSP := rule(superpositions, x, y)
					if newSP != nil {
						superpositions.set(x, y, newSP)
					}
				}
				newCV := superpositions.get(x, y).collapsedValue()
				if newCV != nothing {
					autoCollapsed = true
					applyCollapsedValue(x, y, newCV, g, s, f)
					if collapseCounter == 0 {
						exportFunc(m)
						collapseCounter = cd
						cd *= 2
					}
					collapseCounter -= 1
				}
			}
		}

		// Choose the non-collapsed tile with lowest entropy
		minEntropy := math.Inf(0)
		minX := 0
		minY := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				sp := superpositions.get(x, y)
				if sp.isCollapsed() {
					continue
				}
				entropy := sp.entropy()
				if entropy < minEntropy {
					minEntropy = entropy
					minX = x
					minY = y
				}
			}
		}
		if minEntropy == math.Inf(0) && !autoCollapsed {
			// We can't collapse anymore
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
		applyCollapsedValue(minX, minY, maxKey, g, s, f)
		if collapseCounter == 0 {
			exportFunc(m)
			collapseCounter = cd
			cd *= 2
		}
		collapseCounter -= 1
	}
	exportFunc(m)
	// Collapse uncollapsed tiles with default rule
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sp := superpositions.get(x, y)
			if sp.isCollapsed() {
				continue
			}
			applyCollapsedValue(x, y, defaultTile, g, s, f)
		}
	}
	exportFunc(m)

	return m
}

func applyCollapsedValue(x, y int, t rune, g, s, f *Layer) {
	switch t {
	case wall:
		s.setTile(x, y, t)
	case grass:
		g.setTile(x, y, t)
	case road:
		g.setTile(x, y, t)
	default:
		panic(fmt.Sprintf("unknown tile %s", string(t)))
	}
}

type Superposition map[rune]float64

func (s Superposition) entropy() float64 {
	if len(s) == 0 {
		return math.Inf(0)
	}
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

func (s Superposition) isCollapsed() bool {
	return s.collapsedValue() != nothing
}

type Superpositions struct {
	sp     []Superposition
	m      *Map
	Width  int
	Height int
}

func (s Superpositions) get(x, y int) Superposition {
	return s.sp[x+y*s.Width]
}

func (s *Superpositions) set(x, y int, sp Superposition) {
	s.sp[x+y*s.Width] = sp
}

// Count the number of collapsed values around a tile that match a certain tile
// Boundary tiles don't count
func (s *Superpositions) countCollapsed(x, y, r int, tile rune) int {
	c := 0
	for xi := x - r; xi <= x+r; xi++ {
		for yi := y - r; yi <= y+r; yi++ {
			if xi < 0 || xi >= s.m.Width || yi < 0 || yi >= s.m.Height {
				continue
			}
			if s.get(xi, yi).collapsedValue() == tile {
				c++
			}
		}
	}
	return c
}

func newSuperpositions(m *Map) *Superpositions {
	s := new(Superpositions)
	s.m = m
	s.Width, s.Height = m.Width, m.Height
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			s.sp = append(s.sp, Superposition{})
		}
	}
	return s
}

type Rule func(s *Superpositions, x, y int) Superposition

func roadAtEdgeRule(s *Superpositions, x, y int) Superposition {
	if x == 0 || y == 0 || x == s.m.Width-1 || y == s.m.Height-1 {
		return Superposition{road: 1.0}
	}
	return nil
}

func wallsInsideRoadRule(s *Superpositions, x, y int) Superposition {
	if s.get(x, y).collapsedValue() != road && s.countCollapsed(x, y, 1, road) >= 1 {
		return Superposition{wall: 1.0}
	}
	return nil
}
