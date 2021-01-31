package gmgmap

import (
	"fmt"
	"math/rand"

	"github.com/beefsack/go-astar"
)

// Layer - a rectangular collection of tiles
type Layer struct {
	Name   string
	Tiles  []rune
	Width  int
	Height int
}

// Map - a rectangular tile map
type Map struct {
	Layers []*Layer
	Width  int
	Height int
}

// Tile types
const (
	nothing    = ' '
	floor      = 'f'
	floor2     = 'F'
	road       = 'r'
	road2      = 'R'
	wall       = 'w'
	wall2      = 'W'
	room       = '.'
	room2      = '#'
	door       = '+'
	doorLocked = 'x'
	stairsUp   = '<'
	stairsDown = '>'
	tree       = 'T'
	grass      = 'g'

	// flavour
	sign       = 's'
	hanging    = 'h' // stuff that goes on indoor walls
	window     = 'o'
	counter    = '_'
	shopkeeper = 'A'
	shelf      = 'S'
	stock      = ')'
	table      = 't'
	chair      = 'c'
	rug        = '~'
	pot        = '{'
	assistant  = 'a'
	player     = '@'
	flower     = 'v'
	key        = '('
)

// NewMap - create a new Map for a certain size
func NewMap(width, height int) *Map {
	m := new(Map)
	m.Width = width
	m.Height = height
	return m
}

func newLayer(name string, width, height int) *Layer {
	l := new(Layer)
	l.Name = name
	l.Width, l.Height = width, height
	l.Tiles = make([]rune, width*height)
	l.fill(nothing)
	return l
}

// Layer - get a map layer by name
// If it doesn't exist, add the layer
func (m *Map) Layer(name string) *Layer {
	for _, l := range m.Layers {
		if l.Name == name {
			return l
		}
	}
	m.Layers = append(m.Layers, newLayer(name, m.Width, m.Height))
	return m.Layers[len(m.Layers)-1]
}

func (m *Map) removeLayer(name string) {
	for i, l := range m.Layers {
		if l.Name == name {
			m.Layers = append(m.Layers[:i], m.Layers[i+1:]...)
			return
		}
	}
}

func (l Layer) getTile(x, y int) rune {
	if x < 0 || x >= l.Width || y < 0 || y >= l.Height {
		return rune(0)
	}
	return l.Tiles[x+y*l.Width]
}

func (l *Layer) setTile(x, y int, tile rune) {
	l.Tiles[x+y*l.Width] = tile
}

func (l *Layer) setTileInAreaIfEmpty(r rect, tile rune) {
	for i := 0; i < 100; i++ {
		x := rand.Intn(r.w) + r.x
		y := rand.Intn(r.h) + r.y
		if l.getTile(x, y) == nothing {
			l.setTile(x, y, tile)
			break
		}
	}
}

func (l Layer) isIn(x, y int) bool {
	return x >= 0 && x < l.Width && y >= 0 && y < l.Height
}

// Fill the map with a single tile type
func (l *Layer) fill(tile rune) {
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			l.setTile(x, y, tile)
		}
	}
}

// Draw a rectangle - optional filled
func (l *Layer) rectangle(r rect, tile rune, filled bool) {
	for y := r.y; y < r.y+r.h; y++ {
		for x := r.x; x < r.x+r.w; x++ {
			if filled || x == r.x || y == r.y || x == r.x+r.w-1 || y == r.y+r.h-1 {
				l.setTile(x, y, tile)
			}
		}
	}
}

func (l *Layer) rectangleFilled(r rect, tile rune) {
	l.rectangle(r, tile, true)
}

func (l *Layer) rectangleUnfilled(r rect, tile rune) {
	l.rectangle(r, tile, false)
}

// Perform a flood fill starting from a location
// Floods up, down, left and right
func (l *Layer) floodFill(x, y int, tile rune) {
	indices := []int{x + y*l.Width}
	floodTile := l.Tiles[indices[0]]
	l.Tiles[indices[0]] = tile
	for i := 0; i < len(indices); i++ {
		x = indices[i] % l.Width
		y = indices[i] / l.Width
		var index int
		// top
		index = (y-1)*l.Width + x
		if y > 0 && l.Tiles[index] == floodTile {
			indices = append(indices, index)
			l.Tiles[index] = tile
		}
		// bottom
		index = (y+1)*l.Width + x
		if y < l.Height-1 && l.Tiles[index] == floodTile {
			indices = append(indices, index)
			l.Tiles[index] = tile
		}
		// left
		index = y*l.Width + x - 1
		if x > 0 && l.Tiles[index] == floodTile {
			indices = append(indices, index)
			l.Tiles[index] = tile
		}
		// right
		index = y*l.Width + x + 1
		if x < l.Width-1 && l.Tiles[index] == floodTile {
			indices = append(indices, index)
			l.Tiles[index] = tile
		}
	}
}

// Print - print map in ascii, with a border
func (m Map) Print() {
	for y := 0; y < m.Height; y++ {
		// Upper frame
		if y == 0 {
			fmt.Print("+")
			for x := 0; x < m.Width; x++ {
				fmt.Print("-")
			}
			fmt.Print("+")
			fmt.Println()
		}

		// Left of frame
		fmt.Print("|")

		// Interior cells
		for x := 0; x < m.Width; x++ {
			// Print the top-most cell in the Layers
			for i := len(m.Layers) - 1; i >= 0; i-- {
				l := m.Layers[i]
				tile := l.getTile(x, y)
				if i == 0 || tile != nothing {
					fmt.Printf("%c", tile)
					break
				}
			}
		}

		// Right of frame
		fmt.Print("|")

		// Bottom frame
		if y == m.Height-1 {
			fmt.Println()
			fmt.Print("+")
			for x := 0; x < m.Width; x++ {
				fmt.Print("-")
			}
			fmt.Print("+")
		}

		fmt.Println()
	}
}

// PrintCSV - print raw rune values as CSV
func (m Map) PrintCSV() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			// Print the top-most cell in the Layers
			printed := false
			for i := len(m.Layers) - 1; i >= 0; i-- {
				l := m.Layers[i]
				tile := l.getTile(x, y)
				if i == 0 || tile != nothing {
					fmt.Printf("%d", tile)
					printed = true
					break
				}
			}
			if !printed {
				fmt.Print(" ")
			}
			fmt.Print(",")
		}
		fmt.Println()
	}
}

// Check if rectangular area is clear, i.e. only composed of nothing tiles
func (l Layer) isClear(roomX, roomY, roomWidth, roomHeight int) bool {
	for x := roomX; x < roomX+roomWidth; x++ {
		for y := roomY; y < roomY+roomHeight; y++ {
			if l.getTile(x, y) != nothing {
				return false
			}
		}
	}
	return true
}

// Count the number of tiles around a tile that match a certain tile
// Boundary tiles count
func (l Layer) countTiles(x, y, r int, tile rune) int {
	c := 0
	for xi := x - r; xi <= x+r; xi++ {
		for yi := y - r; yi <= y+r; yi++ {
			if xi < 0 || xi >= l.Width || yi < 0 || yi >= l.Height {
				c++
			} else if l.getTile(xi, yi) == tile {
				c++
			}
		}
	}
	return c
}

// IsWall - whether a tile is a wall type
func IsWall(tile rune) bool {
	return tile == wall || tile == wall2
}

// IsDoor - whether a tile is a door type
func IsDoor(tile rune) bool {
	return tile == door || tile == doorLocked
}

// Add a corridor with two turns
// This can connect any two points; the S-shaped turn occurs at the middle
func addCorridor(g, s *Layer, startX, startY, endX, endY int, tile rune) {
	deltax := startX - endX
	if deltax < 0 {
		deltax = -deltax
	}
	deltay := startY - endY
	if deltay < 0 {
		deltay = -deltay
	}
	dx := 0
	if deltax > deltay {
		dx = 1
	}
	dy := 1 - dx
	var dxAlt, dyAlt int
	var halfX, halfY int
	if dx > 0 {
		// horizontal
		dx = 1
		dy = 0
		if startX > endX {
			tmp := startX
			startX = endX
			endX = tmp
			tmp = startY
			startY = endY
			endY = tmp
		}
		dxAlt, dyAlt = 0, 1
		halfX, halfY = (endX-startX)/2+startX, endY+1
		if endY < startY {
			dyAlt = -1
			halfY = endY - 1
		}
	} else {
		// vertical
		dx = 0
		dy = 1
		if startY > endY {
			tmp := startX
			startX = endX
			endX = tmp
			tmp = startY
			startY = endY
			endY = tmp
		}
		dxAlt, dyAlt = 1, 0
		halfX, halfY = endX+1, (endY-startY)/2+startY
		if endX < startX {
			dxAlt = -1
			halfX = endX - 1
		}
	}
	set := func(x, y int) {
		g.setTile(x, y, tile)
		// Clear walls in the way
		if s != nil {
			s.setTile(x, y, nothing)
		}
	}
	// Initial direction
	x, y := startX, startY
	for ; x != halfX && y != halfY; x, y = x+dx, y+dy {
		set(x, y)
	}
	// Turn
	for ; x != endX && y != endY; x, y = x+dxAlt, y+dyAlt {
		set(x, y)
	}
	// Finish
	for ; x != endX || y != endY; x, y = x+dx, y+dy {
		set(x, y)
	}
	set(endX, endY)
}

// Tile - Single tile on the map for astar
type Tile struct {
	x, y int
	s    *Layer
	w    World
}

// PathNeighbors - Get neighbours for astar pathfinding
func (t *Tile) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.s.getTile(t.x+offset[0], t.y+offset[1]); n == nothing {
			neighbors = append(neighbors, t.w.tile(t.x+offset[0], t.y+offset[1]))
		}
	}
	return neighbors
}

// PathNeighborCost - cost of traveling to neighbour for astar
func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

// PathEstimatedCost - heuristic cost of path for astar, using manhattan distance
func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	return float64(manhattanDistance(t.x, t.y, toT.x, toT.y))
}

// World - 2D array of tiles
type World map[int]map[int]*Tile

func (w World) tile(x, y int) *Tile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

func (w World) setTile(t *Tile, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*Tile{}
	}
	w[x][y] = t
	t.x = x
	t.y = y
	t.w = w
}

// Use A* to find and return a path between two points
// A* will avoid any tiles where there's something in the structure (s) layer
func addPath(g, s *Layer, x1, y1, x2, y2 int) (path []astar.Pather, distance float64, found bool) {
	w := World{}
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			w.setTile(&Tile{x, y, s, w}, x, y)
		}
	}
	return astar.Path(w.tile(x1, y1), w.tile(x2, y2))
}
