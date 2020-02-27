package gmgmap

import "math/rand"

// NewVillage - create a village, made up of multiple buildings
func NewVillage(width, height, buildingPadding int) *Map {
	m := NewMap(width, height)
	g := m.Layer("Ground")
	s := m.Layer("Structures")
	f := m.Layer("Furniture")

	// Grass
	g.fill(grass)

	// Buildings
	buildings := make([]rect, 0)
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
		for _, r := range buildings {
			// Add a bit of padding between the buildings
			if rectOverlaps(
				r,
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
		buildings = append(buildings, rect{x, y, w, h})
	}

	// Draw paths between random pairs of entrances
	type Pair struct {
		a, b int
	}
	paths := map[Pair]bool{}
	for i := 0; i < len(buildings)-1; i++ {
		for {
			// Check for path valid and exists
			building1 := rand.Intn(len(buildings))
			building2 := rand.Intn(len(buildings))
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
			startX := b1.x + b1.w/2
			startY := b1.y + b1.h
			endX := b2.x + b2.w/2
			endY := b2.y + b2.h
			addPath(g, s, startX, startY, endX, endY, road)
			//addCorridor(g, s, startX, startY, endX, endY, road)
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
