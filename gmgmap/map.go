package gmgmap

import "fmt"

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
	floor3     = 's'
	wall       = 'w'
	wall2      = 'W'
	room       = '.'
	room2      = '#'
	door       = '+'
	stairsUp   = '<'
	stairsDown = '>'
	tree       = 't'
)

// NewMap - create a new Map for a certain size
func NewMap(width, height int) *Map {
	m := new(Map)
	m.Width = width
	m.Height = height
	m.Layers = append(m.Layers, newLayer("Furniture", width, height))
	m.Layers = append(m.Layers, newLayer("Tiles", width, height))
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
func (m Map) Layer(name string) *Layer {
	for _, l := range m.Layers {
		if l.Name == name {
			return l
		}
	}
	return nil
}

func (l Layer) getTile(x, y int) rune {
	return l.Tiles[x+y*l.Width]
}

func (l Layer) setTile(x, y int, tile rune) {
	l.Tiles[x+y*l.Width] = tile
}

// Fill the map with a single tile type
func (l Layer) fill(tile rune) {
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			l.setTile(x, y, tile)
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
			for index, layer := range m.Layers {
				tile := layer.getTile(x, y)
				if index == len(m.Layers)-1 || tile != nothing {
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
