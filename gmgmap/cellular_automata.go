package gmgmap

import "math/rand"

const (
	floorTile = floor
	roadTile  = road2
)

// NewCellularAutomata - create a stone-on-floor map using cellular automata
// For a number of repetitions:
// If the number of stones within one step (including itself) is at least r1 OR
// the number of stones within 2 steps at most r2, turn into a stone,
// else turn into a floor
func NewCellularAutomata(width, height, fillPct, repeat, r1, r2 int) *Map {

	m := NewMap(width, height)
	g := m.Layer("Ground")
	g.fill(floorTile)
	l := m.Layer("Structures")
	// Randomly set a percentage of the tiles as stones
	for i := 0; i < fillPct*width*height/100; i++ {
		l.Tiles[i] = roadTile
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

	// Use flood fill to identify disconnected areas
	fl := m.Layer("Flood")
	fl.fill(0)
	// First copy across the stone tiles
	for i, tile := range l.Tiles {
		if tile == roadTile {
			fl.Tiles[i] = roadTile
		}
	}
	// Then perform flood fill conditionally on the flood layer
	index := rune(-1)
	for i := range fl.Tiles {
		if fl.Tiles[i] == 0 {
			fl.floodFill(i%fl.Width, i/fl.Width, index)
			index--
		}
	}

	numAreas := -index - 1
	// Connect the disconnected areas, first to second, second to third etc.
	// Select random tile from each area
	areaTiles := make([]rune, len(fl.Tiles))
	for i := range areaTiles {
		areaTiles[i] = rune(i)
	}
	for i := range areaTiles {
		j := rand.Intn(i + 1)
		areaTiles[i], areaTiles[j] = areaTiles[j], areaTiles[i]
	}
	areaStarts := make([]rune, numAreas)
	for _, index := range areaTiles {
		tile := -fl.Tiles[index] - 1
		if tile >= 0 && tile < numAreas && areaStarts[tile] == 0 {
			areaStarts[tile] = index
		}
	}
	m.removeLayer(fl.Name)
	// Connect consecutive pairs of tiles
	for i := 0; i < len(areaStarts)-1; i++ {
		x1 := int(areaStarts[i]) % l.Width
		y1 := int(areaStarts[i]) / l.Width
		x2 := int(areaStarts[i+1]) % l.Width
		y2 := int(areaStarts[i+1]) / l.Width
		addCorridor(g, l, x1, y1, x2, y2, floorTile)
	}

	return m
}

func rep(l *Layer, r1, r2 int) {
	buf := make([]rune, len(l.Tiles))
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			i := x + y*l.Width
			if l.countTiles(x, y, 1, roadTile) >= r1 ||
				l.countTiles(x, y, 2, roadTile) <= r2 {
				buf[i] = roadTile
			} else {
				buf[i] = floorTile
			}
		}
	}
	l.Tiles = buf
}
