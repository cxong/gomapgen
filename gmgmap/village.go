package gmgmap

import "math/rand"

// NewVillage - create a village, made up of multiple buildings
func NewVillage(width, height, buildingPadding int) *Map {
	m := NewMap(width, height)

	// Grass
	g := m.Layer("Ground")
	g.fill(grass)

	// Buildings
	s := m.Layer("Structures")
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
		addBuilding(s, x, y, w, h)
		buildings = append(buildings, rect{x, y, w, h})
	}

	return m
}

func addBuilding(s *Layer, x, y, w, h int) {
	// Perimeter
	s.rectangle(rect{x, y, w, h}, wall, false)
}
