package gmgmap

import "fmt"

// Map - a rectangular tile map
type Map struct {
	Tiles  []byte
	Width  int
	Height int
}

// Tile types
const (
	nothing = ' '
	floor   = 'f'
	wall    = 'w'
	room    = 'r'
)

// NewMap - create a new Map for a certain size
func NewMap(width, height int) *Map {
	m := new(Map)
	m.Width = width
	m.Height = height
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			m.Tiles = append(m.Tiles, nothing)
		}
	}
	return m
}

// GetTile - get tile at x, y
func (m Map) GetTile(x, y int) byte {
	return m.Tiles[x+y*m.Width]
}

// SetTile - set tile at x, y
func (m Map) SetTile(x, y int, tile byte) {
	m.Tiles[x+y*m.Width] = tile
}

// Fill the map with a single tile type
func (m Map) Fill(tile byte) {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			m.SetTile(x, y, tile)
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
			var i = y*m.Width + x
			fmt.Printf("%c", m.Tiles[i])
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

// ToTMX - export map as TMX (Tiled XML map)
func (m Map) ToTMX() {

}
