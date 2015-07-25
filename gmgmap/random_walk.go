package gmgmap

// NewRandomWalk - create a room-on-floor map using random walk algorithm
func NewRandomWalk(width, height, iterations int) *Map {
	m := NewMap(width, height)
	m.Layer("Tiles").fill(floor)
	// Start walking from the middle, randomly
	x, y := width/2, height/2
	l := m.Layer("Furniture")
	for i := 0; i < iterations; i++ {
		l.setTile(x, y, tree)
		x, y = randomWalk(x, y, m.Width, m.Height)
	}
	return m
}
