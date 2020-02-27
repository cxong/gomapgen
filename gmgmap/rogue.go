package gmgmap

import "math/rand"

type connectInfo struct {
	up    bool
	right bool
	down  bool
	left  bool
}

// NewRogue - generate a new Rogue-like map, with rooms connected with tunnels
func NewRogue(width, height,
	gridWidth, gridHeight, minRoomPct, maxRoomPct int) *Map {
	m := NewMap(width, height)

	// Divide into grid, with flags marking grid connections
	totalGrids := gridWidth * gridHeight
	connected := make([]connectInfo, totalGrids)

	// Pick random grid to start with
	gridIndex := rand.Intn(len(connected))
	firstRoomIndex := gridIndex
	var lastRoomIndex int
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
			lastRoomIndex = neighbourIndex
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
					lastRoomIndex = gridIndex
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

	g := m.Layer("Ground")
	s := m.Layer("Structures")

	// Try to place rooms - one for each grid
	numRooms := (rand.Intn(maxRoomPct-minRoomPct) + minRoomPct) * totalGrids / 100
	roomIndices := rand.Perm(totalGrids)
	rooms := make([]rect, totalGrids)
	gridWidthTiles := width / gridWidth
	gridHeightTiles := height / gridHeight
	for i := 0; i < totalGrids; i++ {
		// Coordinates of grid top-left corner
		grid.x, grid.y = roomIndices[i]%gridWidth, roomIndices[i]/gridWidth
		grid.w, grid.h = gridWidth, gridHeight
		gridStartX, gridStartY := grid.x*gridWidthTiles, grid.y*gridHeightTiles
		var roomRect rect
		// force dead ends to be rooms
		numConnections := connected[roomIndices[i]].numConnections(grid)
		// also force first/last room to be rooms
		if i < numRooms || numConnections <= 1 ||
			roomIndices[i] == firstRoomIndex || roomIndices[i] == lastRoomIndex {
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
				if roomRect.w > 1 &&
					(x == roomRect.x || x == roomRect.x+roomRect.w-1 ||
						y == roomRect.y || y == roomRect.y+roomRect.h-1) {
					s.setTile(x, y, wall2)
				} else {
					g.setTile(x, y, room)
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
			addCorridor(g, s, roomRect.x+roomRect.w-1, roomRect.y+roomRect.h/2,
				neighbour.x, neighbour.y+neighbour.h/2, room2)
		}
		if connections.down && y < gridHeight-1 {
			// Connect with neighbour below
			neighbour := rooms[i+gridWidth]
			addCorridor(g, s, roomRect.x+roomRect.w/2, roomRect.y+roomRect.h-1,
				neighbour.x+neighbour.w/2, neighbour.y, room2)
		}
	}

	// Find door tiles: those with 2 neighbour walls and 1 each of corridor/room
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if IsWall(s.getTile(x, y)) {
				continue
			}
			walls, corridors, rooms := 0, 0, 0
			var countTile = func(x, y int) {
				if IsWall(s.getTile(x, y)) {
					walls++
				} else {
					switch g.getTile(x, y) {
					case room:
						rooms++
					case room2:
						corridors++
					}
				}
			}
			if y > 0 {
				countTile(x, y-1)
			}
			if x < m.Width-1 {
				countTile(x+1, y)
			}
			if y < m.Height-1 {
				countTile(x, y+1)
			}
			if x > 0 {
				countTile(x-1, y)
			}
			if walls == 2 && corridors == 1 && rooms == 1 {
				s.setTile(x, y, door)
			}
		}
	}

	// Put stairs in the first and last room
	firstRoom := rooms[firstRoomIndex]
	lastRoom := rooms[lastRoomIndex]
	s.setTile(firstRoom.x+firstRoom.w/2, firstRoom.y+firstRoom.h/2, stairsUp)
	s.setTile(lastRoom.x+lastRoom.w/2, lastRoom.y+lastRoom.h/2, stairsDown)

	return m
}

// Don't count edges as connections
func (c connectInfo) numConnections(grid rect) int {
	n := 0
	if c.up && grid.y > 0 {
		n++
	}
	if c.right && grid.x < grid.w-1 {
		n++
	}
	if c.down && grid.y < grid.h-1 {
		n++
	}
	if c.left && grid.x > 0 {
		n++
	}
	return n
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
