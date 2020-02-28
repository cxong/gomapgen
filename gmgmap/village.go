package gmgmap

import (
	"fmt"
	"math"
	"math/rand"
)

type Building struct {
	r          rect
	importance int
}

func (b Building) addNPC(c *Layer) {
	// Try to place a random NPC somewhere inside the building
	for i := 0; i < 100; i++ {
		x := rand.Intn(b.r.w-2) + b.r.x + 1
		y := rand.Intn(b.r.h-2) + b.r.y + 1
		if c.getTile(x, y) == nothing {
			c.setTile(x, y, player)
			break
		}
	}
}

// NewVillage - create a village, made up of multiple buildings
func NewVillage(width, height, buildingPadding int) *Map {
	m := NewMap(width, height)
	g := m.Layer("Ground")
	s := m.Layer("Structures")
	f := m.Layer("Furniture")
	c := m.Layer("Characters")

	// Grass
	g.fill(grass)

	// Buildings
	buildings := make([]Building, 0)
	// Keep placing buildings for a while
	for i := 0; i < 500; i++ {
		w := rand.Intn(3) + 5
		h := rand.Intn(3) + 5
		x := rand.Intn(width - w)
		y := rand.Intn(height - h)
		if x < 0 || y < 0 {
			continue
		}
		// Check if it overlaps with any existing buildings
		overlaps := false
		for _, b := range buildings {
			// Add a bit of padding between the buildings
			if rectOverlaps(
				b.r,
				rect{
					x - buildingPadding,
					y - buildingPadding,
					w + buildingPadding*2,
					h + buildingPadding*2}) {
				overlaps = true
				break
			}
		}
		if overlaps {
			continue
		}
		addBuilding(g, s, f, x, y, w, h)
		buildings = append(buildings, Building{rect{x, y, w, h}, 0})
	}

	// Randomly assign importance to buildings
	impSum := 0
	for _, building := range buildings {
		building.importance = int(math.Pow(float64(rand.Intn(3)+1), 2))
		impSum += building.importance
		// Place NPCs based on importance
		for i := 0; i < building.importance; i++ {
			building.addNPC(c)
		}
	}

	// Draw paths between random pairs of entrances via importance
	type Pair struct {
		a, b int
	}
	paths := map[Pair]bool{}
	for i := 0; i < len(buildings)-1; i++ {
		for {
			// Check for path valid and exists
			building1 := rand.Intn(len(buildings))
			// randomly select second building by importance
			impFact := rand.Intn(impSum)
			building2 := 0
			impFactSum := 0
			for i, b2 := range buildings {
				impFactSum += b2.importance
				if impFactSum > impFact {
					building2 = i
					break
				}
			}
			if building1 == building2 {
				continue
			}
			if building1 > building2 {
				building2, building1 = building1, building2
			}
			key := Pair{building1, building2}
			if _, ok := paths[key]; ok {
				continue
			}
			paths[key] = true
			// TODO: find entrance and start/end paths there
			b1 := buildings[building1]
			b2 := buildings[building2]
			startX := b1.r.x + b1.r.w/2
			startY := b1.r.y + b1.r.h
			endX := b2.r.x + b2.r.w/2
			endY := b2.r.y + b2.r.h
			path, _, found := addPath(g, s, startX, startY, endX, endY)
			if !found {
				fmt.Println("Could not find path")
			} else {
				for _, t := range path {
					g.setTile(t.(*Tile).x, t.(*Tile).y, road)
					// Clear walls in the way
					if s != nil {
						s.setTile(t.(*Tile).x, t.(*Tile).y, nothing)
					}
				}
			}
			break
		}
	}

	return m
}

func addBuilding(g, s, f *Layer, x, y, w, h int) {
	// Perimeter
	s.rectangle(rect{x, y, w, h}, wall, false)
	// Floor
	g.rectangle(rect{x + 1, y + 1, w - 2, h - 2}, room, true)
	// Entrance
	entranceX := x + w/2
	entranceY := y + h - 1
	g.setTile(entranceX, entranceY, room)
	s.setTile(entranceX, entranceY, door)
	f.setTile(entranceX-1, entranceY, sign)
}