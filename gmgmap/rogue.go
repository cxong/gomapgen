package gmgmap

import (
	"fmt"
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
	gridWidth, gridHeight, minRoomPct, maxRoomPct int) *Map {
	rand.Seed(time.Now().UTC().UnixNano())
	m := NewMap(width, height)

	// Divide into grid, with flags marking grid connections
	totalGrids := gridWidth * gridHeight
	connected := make([]connectInfo, totalGrids)

	// Pick random grid to start with
	gridIndex := rand.Intn(len(connected))
	grid := rect{gridIndex % gridWidth, gridIndex / gridWidth,
		gridWidth, gridHeight}

	// Connect to a random neighbour
	for {
		// Mark edges as already connected
		if grid.x == 0 {
			connected[gridIndex].left = true
		}
		if grid.y == 0 {
			connected[gridIndex].up = true
		}
		if grid.x == gridWidth-1 {
			connected[gridIndex].right = true
		}
		if grid.y == gridHeight-1 {
			connected[gridIndex].down = true
		}
		// If all neighbours connected, end
		if connected[gridIndex].allConnected() {
			break
		}
		// Otherwise, connect to a random unconnected neighbour
		for {
			neighbourX, neighbourY := randomWalk(grid.x, grid.y, gridWidth, gridHeight)
			neighbourIndex := neighbourX + neighbourY*gridWidth
			if !tryConnect(connected, grid.x, grid.y, neighbourX, neighbourY,
				gridIndex, neighbourIndex) {
				continue
			}
			// Set neighbour as current grid
			grid.x, grid.y, gridIndex = neighbourX, neighbourY, neighbourIndex
			break
		}
	}

	// Scan for unconnected grids; if so try to connect them to a neighbouring
	// connected grid
	for {
		hasUnconnected := false
		for grid.x = 0; grid.x < gridWidth; grid.x++ {
			for grid.y = 0; grid.y < gridHeight; grid.y++ {
				gridIndex = grid.x + grid.y*gridWidth
				if connected[gridIndex].isConnected() {
					// Grid is already connected
					continue
				}
				hasUnconnected = true
				// If no neighbours connected, continue
				if (grid.x == 0 || !connected[gridIndex-1].isConnected()) &&
					(grid.x == gridWidth-1 || !connected[gridIndex+1].isConnected()) &&
					(grid.y == 0 || !connected[gridIndex-gridWidth].isConnected()) &&
					(grid.y == gridHeight-1 || !connected[gridIndex+gridWidth].isConnected()) {
					continue
				}
				// Try connecting to a random connected neighbour
				for {
					neighbourX, neighbourY := randomWalk(grid.x, grid.y, gridWidth, gridHeight)
					neighbourIndex := neighbourX + neighbourY*gridWidth
					if !connected[neighbourIndex].isConnected() {
						continue
					}
					if !tryConnect(connected, grid.x, grid.y, neighbourX, neighbourY,
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
			grid.x = gridIndex % gridWidth
			grid.y = gridIndex / gridWidth
			if connected[gridIndex].allConnected() {
				break
			}
			// Try connecting to a random neighbour
			for {
				neighbourX, neighbourY := randomWalk(grid.x, grid.y, gridWidth, gridHeight)
				neighbourIndex := neighbourX + neighbourY*gridWidth
				if tryConnect(connected, grid.x, grid.y, neighbourX, neighbourY,
					gridIndex, neighbourIndex) {
					break
				}
			}
			break
		}
	}

	// Try to place rooms - one for each grid
	numRooms := (rand.Intn(maxRoomPct-minRoomPct) + minRoomPct) * totalGrids / 100
	roomIndices := rand.Perm(totalGrids)
	rooms := make([]rect, totalGrids)
	gridWidthTiles := width / gridWidth
	gridHeightTiles := height / gridHeight
	for i := 0; i < totalGrids; i++ {
		// Coordinates of grid top-left corner
		gridStartX := (roomIndices[i] % gridWidth) * gridWidthTiles
		gridStartY := (roomIndices[i] / gridWidth) * gridHeightTiles
		var roomRect rect
		if i < numRooms {
			// Generate random room
			roomRect.w = rand.Intn(gridWidthTiles-4) + 4
			roomRect.h = rand.Intn(gridHeightTiles-4) + 4
		} else {
			// Generate "gone rooms"
			roomRect.w = 1
			roomRect.h = 1
		}
		// Place the room
		roomRect.x = rand.Intn(width/gridWidth-roomRect.w) + gridStartX
		roomRect.y = rand.Intn(height/gridHeight-roomRect.h) + gridStartY
		for x := roomRect.x; x < roomRect.x+roomRect.w; x++ {
			for y := roomRect.y; y < roomRect.y+roomRect.h; y++ {
				tile := room
				if roomRect.w > 1 &&
					(x == roomRect.x || x == roomRect.x+roomRect.w-1 ||
						y == roomRect.y || y == roomRect.y+roomRect.h-1) {
					tile = wall2
				}
				if err := m.SetTile(x, y, tile); err != nil {
					panic(err)
				}
			}
		}
		rooms[roomIndices[i]] = roomRect
	}

	// Connect each room to connected neighbours
	for i := 0; i < totalGrids; i++ {
		connections := connected[i]
		x, y := i%gridWidth, i/gridHeight
		roomRect := rooms[i]
		// Only connect to the right and below
		if connections.right && x < gridWidth-1 {
			// Connect with neighbour on right
			neighbour := rooms[i+1]
			addCorridor(m, roomRect.x+roomRect.w-1, roomRect.y+roomRect.h/2,
				neighbour.x, neighbour.y+neighbour.h/2, 1, 0, room)
		} else if connections.down && y < gridHeight-1 {
			// Connect with neighbour below
			neighbour := rooms[i+gridWidth]
			addCorridor(m, roomRect.x+roomRect.w/2, roomRect.y+roomRect.h-1,
				neighbour.x+neighbour.w/2, neighbour.y, 0, 1, room)
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

func addCorridor(m *Map, startX, startY, endX, endY, dx, dy int, tile rune) {
	var dxAlt, dyAlt int
	var halfX, halfY int
	if dx > 0 {
		// horizontal
		dxAlt, dyAlt = 0, 1
		halfX, halfY = (endX-startX)/2+startX, endY+1
		if endY < startY {
			dyAlt = -1
			halfY = endY - 1
		}
	} else {
		// vertical
		dxAlt, dyAlt = 1, 0
		halfX, halfY = endX+1, (endY-startY)/2+startY
		if endX < startX {
			dxAlt = -1
			halfX = endX - 1
		}
	}
	// Initial direction
	x, y := startX, startY
	fmt.Printf("%d,%d to %d,%d half %d,%d at %d,%d\n",
		x, y, endX, endY, halfX, halfY, dx, dy)
	for ; x != halfX && y != halfY; x, y = x+dx, y+dy {
		if err := m.SetTile(x, y, tile); err != nil {
			panic(err)
		}
	}
	// Turn
	for ; x != endX && y != endY; x, y = x+dxAlt, y+dyAlt {
		if err := m.SetTile(x, y, tile); err != nil {
			panic(err)
		}
	}
	// Finish
	for ; x != endX || y != endY; x, y = x+dx, y+dy {
		if err := m.SetTile(x, y, tile); err != nil {
			panic(err)
		}
	}
	if err := m.SetTile(endX, endY, tile); err != nil {
		panic(err)
	}
}
