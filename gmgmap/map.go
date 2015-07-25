package gmgmap

import "fmt"

// Map - a rectangular tile map
type Map struct {
	Tiles  []rune
	Width  int
	Height int
}

// Tile types
const (
	nothing = ' '
	floor   = 'f'
	floor2  = 'F'
	wall    = 'w'
	wall2   = 'W'
	room    = 'r'
	room2   = 'R'
	door    = 'd'
)

// NewMap - create a new Map for a certain size
func NewMap(width, height int) *Map {
	m := new(Map)
	m.Width = width
	m.Height = height
	m.Tiles = make([]rune, m.Width*m.Height)
	m.Fill(nothing)
	return m
}

// TileOutOfBounds - returned when trying to access a map tile that is out of
// bounds
type TileOutOfBounds struct {
}

func (e *TileOutOfBounds) Error() string { return "Tile out of bounds" }

// GetTile - get tile at x, y
func (m Map) GetTile(x, y int) (rune, error) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return 0, &TileOutOfBounds{}
	}
	return m.Tiles[x+y*m.Width], nil
}

// SetTile - set tile at x, y
func (m Map) SetTile(x, y int, tile rune) error {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return &TileOutOfBounds{}
	}
	m.Tiles[x+y*m.Width] = tile
	return nil
}

// Fill the map with a single tile type
func (m Map) Fill(tile rune) {
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

// Check if rectangular area is clear, i.e. only composed of nothing tiles
func (m Map) isClear(roomX, roomY, roomWidth, roomHeight int) bool {
	for x := roomX; x < roomX+roomWidth; x++ {
		for y := roomY; y < roomY+roomHeight; y++ {
			c, err := m.GetTile(x, y)
			if err != nil {
				panic(err)
			}
			if c != nothing {
				return false
			}
		}
	}
	return true
}

// IsWall - whether a tile is a wall type
func IsWall(tile rune) bool {
	return tile == wall || tile == wall2
}
