package gmgmap

import (
	"math/rand"
	"time"
)

// NewRandomWalk - create a room-on-floor map using random walk algorithm
func NewRandomWalk(width, height, iterations int) *Map {
	m := NewMap(width, height)
	m.Fill(floor)
	// Start walking from the middle, randomly
	rand.Seed(time.Now().UTC().UnixNano())
	x, y := width/2, height/2
	for i := 0; i < iterations; i++ {
		m.SetTile(x, y, floor2)
		x, y = randomWalk(x, y, m.Width, m.Height)
	}
	return m
}
