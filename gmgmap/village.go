package gmgmap

// NewVillage - create a village, made up of multiple buildings
func NewVillage(width, height int) *Map {
	m := NewMap(width, height)

	// Grass
	g := m.Layer("Ground")
	g.fill(grass)

	// Buildings
	s := m.Layer("Structures")
	s.rectangle(rect{1, 1, g.Width - 2, g.Height - 3}, wall, false)

	return m
}
