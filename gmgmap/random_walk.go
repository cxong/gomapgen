package gmgmap

// NewRandomWalk - create a room-on-floor map using random walk algorithm
func NewRandomWalk(width, height, iterations int) *Map {
	m := NewMap(width, height)
	l := m.Layer("Tiles")
	l.fill(floor)
	// Start walking from the middle, randomly
	x, y := width/2, height/2
	for i := 0; i < iterations; i++ {
		l.setTile(x, y, floor2)
		x, y = randomWalk(x, y, m.Width, m.Height)
	}
	return m
}
