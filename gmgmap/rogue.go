package gmgmap

import (
	"math/rand"
	"time"
)

type connectInfo struct {
	up    bool
	right bool
	down  bool
	left  bool
}

// NewRogue - generate a new Rogue-like map, with rooms connected with tunnels
func NewRogue(width, height,
	minRoomSize, maxRoomSize int) *Map {
	rand.Seed(time.Now().UTC().UnixNano())
	m := NewMap(width, height)

	// Divide into 3x3 grids, with flags marking grid connetions
	gridWidth := (width + 2) / 3
	gridHeight := (height + 2) / 3
	connected := make([]connectInfo, gridWidth*gridHeight)

	// Pick random grid to start with
	gridIndex := rand.Intn(len(connected))
	gridX := gridIndex % gridWidth
	gridY := gridIndex / gridWidth

	// Connect to a random neighbour
	for {
		// Mark edges as already connected
		if gridX == 0 {
			connected[gridIndex].left = true
		}
		if gridY == 0 {
			connected[gridIndex].up = true
		}
		if gridX == gridWidth-1 {
			connected[gridIndex].right = true
		}
		if gridY == gridHeight-1 {
			connected[gridIndex].down = true
		}
		// If all neighbours connected, end
		if connected[gridIndex].allConnected() {
			break
		}
		// Otherwise, connect to a random unconnected neighbour
		for {
			neighbourX, neighbourY := randomWalk(gridX, gridY, gridWidth, gridHeight)
			neighbourIndex := neighbourX + neighbourY*gridWidth
			if !tryConnect(connected, gridX, gridY, neighbourX, neighbourY,
				gridIndex, neighbourIndex) {
				continue
			}
			// Set neighbour as current grid
			gridX, gridY, gridIndex = neighbourX, neighbourY, neighbourIndex
			break
		}
	}

	// Scan for unconnected grids; if so try to connect them to a neighbouring
	// connected grid
	for {
		hasUnconnected := false
		for gridX = 0; gridX < gridWidth; gridX++ {
			for gridY = 0; gridY < gridHeight; gridY++ {
				gridIndex = gridX + gridY*gridWidth
				if connected[gridIndex].isConnected() {
					// Grid is already connected
					continue
				}
				hasUnconnected = true
				// If no neighbours connected, continue
				if (gridX == 0 || !connected[gridIndex-1].isConnected()) &&
					(gridX == gridWidth-1 || !connected[gridIndex+1].isConnected()) &&
					(gridY == 0 || !connected[gridIndex-gridWidth].isConnected()) &&
					(gridY == gridHeight-1 || !connected[gridIndex+gridWidth].isConnected()) {
					continue
				}
				// Try connecting to a random connected neighbour
				for {
					neighbourX, neighbourY := randomWalk(gridX, gridY, gridWidth, gridHeight)
					neighbourIndex := neighbourX + neighbourY*gridWidth
					if !connected[neighbourIndex].isConnected() {
						continue
					}
					if !tryConnect(connected, gridX, gridY, neighbourX, neighbourY,
						gridIndex, neighbourIndex) {
						panic("unexpected error")
					}
					break
				}
			}
		}
		// Continue until no unconnected grids
		if !hasUnconnected {
			break
		}
	}

	// Make some random connections
	extraConnections := rand.Intn(gridWidth)
	for i := 0; i < extraConnections; i++ {
		for {
			gridIndex = rand.Intn(len(connected))
			gridX = gridIndex % gridWidth
			gridY = gridIndex / gridWidth
			if connected[gridIndex].allConnected() {
				break
			}
			// Try connecting to a random neighbour
			for {
				neighbourX, neighbourY := randomWalk(gridX, gridY, gridWidth, gridHeight)
				neighbourIndex := neighbourX + neighbourY*gridWidth
				if tryConnect(connected, gridX, gridY, neighbourX, neighbourY,
					gridIndex, neighbourIndex) {
					break
				}
			}
			break
		}
	}

	// Try to place rooms
	for i := 0; i < 100; i++ {
		// Generate random room
		roomWidth := rand.Intn(maxRoomSize-minRoomSize) + minRoomSize
		roomHeight := rand.Intn(maxRoomSize-minRoomSize) + minRoomSize
		roomX, roomY := rand.Intn(width-roomWidth), rand.Intn(height-roomHeight)
		// Check if the room overlaps with anything
		if !m.isClear(roomX, roomY, roomWidth, roomHeight) {
			continue
		}
		// Place the room
		for x := roomX; x < roomX+roomWidth; x++ {
			for y := roomY; y < roomY+roomHeight; y++ {
				tile := room
				if x == roomX || x == roomX+roomWidth-1 ||
					y == roomY || y == roomY+roomHeight-1 {
					tile = wall
				}
				if err := m.SetTile(x, y, tile); err != nil {
					panic(err)
				}
			}
		}
	}
	return m
}

func (c connectInfo) isConnected() bool {
	return c.up || c.right || c.down || c.left
}
func (c connectInfo) allConnected() bool {
	return c.up && c.right && c.down && c.left
}
func tryConnect(connected []connectInfo, x, y, x1, y1, index, index2 int) bool {
	switch {
	case y > y1:
		// up
		if connected[index].up {
			return false
		}
		connected[index].up = true
		connected[index2].down = true
	case x < x1:
		// right
		if connected[index].right {
			return false
		}
		connected[index].right = true
		connected[index2].left = true
	case y < y1:
		// down
		if connected[index].down {
			return false
		}
		connected[index].down = true
		connected[index2].up = true
	case x > x1:
		// left
		if connected[index].left {
			return false
		}
		connected[index].left = true
		connected[index2].right = true
	}
	return true
}
