package gmgmap

import "math/rand"

// NewVillage - create a village, made up of multiple buildings
func NewVillage(width, height int) *Map {
	m := NewMap(width, height)

	// Grass
	g := m.Layer("Ground")
	g.fill(grass)

	// Buildings
	s := m.Layer("Structures")
	for i := 0; i < 5; i++ {
		w := rand.Intn(3) + 5
		h := rand.Intn(3) + 5
		x := rand.Intn(width - w)
		y := rand.Intn(height - h)
		if x < 0 || y < 0 {
			continue
		}
		addBuilding(s, x, y, w, h)
	}

	return m
}

func addBuilding(s *Layer, x, y, w, h int) {
	// Perimeter
	s.rectangle(rect{x, y, w, h}, wall, false)
}
