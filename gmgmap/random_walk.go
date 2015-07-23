package gmgmap

import (
	"math/rand"
	"time"
)

// NewRandomWalk - create a room-on-floor map using random walk algorithm
func NewRandomWalk(width, height, iterations int) *Map {
	rand.Seed(time.Now().UTC().UnixNano())
	m := NewMap(width, height)
	m.Fill(floor)
	// Start walking from the middle, randomly
	x, y := width/2, height/2
	for i := 0; i < iterations; i++ {
		m.SetTile(x, y, floor2)
		x, y = nextDirection(m, x, y)
	}
	return m
}

func nextDirection(m *Map, x, y int) (int, int) {
	for {
		// Choose random direction, up/right/down/left
		switch rand.Intn(4) {
		case 0:
			// up
			if y > 0 {
				return x, y - 1
			}
		case 1:
			// right
			if x < m.Width-1 {
				return x + 1, y
			}
		case 2:
			// down
			if y < m.Height-1 {
				return x, y + 1
			}
		case 3:
			// left
			if x > 0 {
				return x - 1, y
			}
		}
	}
}
