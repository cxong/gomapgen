package gmgmap

import "math/rand"

// NewCellularAutomata - create a stone-on-floor map using cellular automata
// For a number of repetitions:
// If the number of stones within one step (including itself) is at least r1 OR
// the number of stones within 2 steps at most r2, turn into a stone,
// else turn into a floor
func NewCellularAutomata(width, height, fillPct, repeat, r1, r2 int) *Map {
	m := NewMap(width, height)
	m.Layer("Ground").fill(floor)
	l := m.Layer("Structures")
	// Randomly set a percentage of the tiles as stones
	for i := 0; i < fillPct*width*height/100; i++ {
		l.Tiles[i] = road
	}
	// Shuffle
	for i := range l.Tiles {
		j := rand.Intn(i + 1)
		l.Tiles[i], l.Tiles[j] = l.Tiles[j], l.Tiles[i]
	}
	// Repetitions
	for i := 0; i < repeat; i++ {
		rep(l, r1, r2)
	}
	return m
}

func rep(l *Layer, r1, r2 int) {
	buf := make([]rune, len(l.Tiles))
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			i := x + y*l.Width
			if l.countTiles(x, y, 1, road) >= r1 ||
				l.countTiles(x, y, 2, road) <= r2 {
				buf[i] = road
			} else {
				buf[i] = floor
			}
		}
	}
	l.Tiles = buf
}
