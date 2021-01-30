package gmgmap

import (
	"math/rand"
)

// Lobby placement - on edge, off the edge, or any position
const (
	LobbyEdge = iota
	LobbyInterior
	LobbyAny
)

// NewInterior - create a building interior layout map.
// The layout has a "lobby", multiple rooms that all connect to the lobby.
// Idea taken from http://www.redactedgame.com/?p=106
// TODO: improve connectedness of leaf nodes, to make it less tree-like
func NewInterior(width, height, minRoomSize, maxRoomSize,
	lobbyEdge int) *Map {
	m := NewMap(width, height)

	// We'll place the "road" tiles later
	g := m.Layer("Ground")
	s := m.Layer("Structures")

	// Randomly partition the space using bsp
	// Keep splitting as long as we can
	var rooms []bspRoom
	rooms = append(rooms, bspRoomRoot(width, height))
	for i := 0; i < len(rooms); i++ {
		if r1, r2, err := rooms[i].Split(i, minRoomSize, maxRoomSize); err == nil {
			rooms[i].child1 = len(rooms)
			rooms = append(rooms, r1)
			rooms[i].child2 = len(rooms)
			rooms = append(rooms, r2)
		}
	}
	// Discard non-leaf rooms
	for i := 0; i < len(rooms); i++ {
		if !rooms[i].IsLeaf() {
			rooms[i] = rooms[len(rooms)-1]
			rooms = rooms[0 : len(rooms)-1]
			i--
		}
	}
	// Adjust the rooms;
	// BSP algo produces rooms with non-overlapping walls
	for i := 0; i < len(rooms); i++ {
		r := rooms[i].r
		if r.x > 0 {
			r.x--
			r.w++
		}
		if r.y > 0 {
			r.y--
			r.h++
		}
		rooms[i].r = r
	}

	// Form the room walls
	for i := 0; i < len(rooms); i++ {
		r := rooms[i].r
		s.rectangleUnfilled(r, wall2)
		groundRect := rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}
		g.rectangleFilled(groundRect, room)
	}

	// Choose one of the rooms to be the lobby
	var lobby bspRoom
	var roomIndices = rand.Perm(len(rooms))
	for i := 0; i < len(roomIndices); i++ {
		lobby = rooms[roomIndices[i]]
		// Check if the lobby placement is ok;
		// If it's on the edge or in the interior
		if lobbyEdge == LobbyEdge {
			if lobby.r.x == 0 || lobby.r.y == 0 ||
				lobby.r.x+lobby.r.w == width || lobby.r.y+lobby.r.h == height {
				break
			}
		} else if lobbyEdge == LobbyInterior {
			if lobby.r.x > 0 && lobby.r.y > 0 &&
				lobby.r.x+lobby.r.w < width && lobby.r.y+lobby.r.h < height {
				break
			}
		} else {
			break
		}
	}
	// Place the lobby
	lobbyRect := lobby.r
	lobbyRect.x++
	lobbyRect.y++
	lobbyRect.w -= 2
	lobbyRect.h -= 2
	g.rectangleFilled(lobbyRect, room2)

	// Mark all the rooms according to their distance from the lobby (depth)
	// Re-use the level parameter
	overlapSize := 1
	for i := 0; i < len(rooms); i++ {
		room := rooms[i]
		room.level = -1
		if room.r.x == lobby.r.x && room.r.y == lobby.r.y {
			room.level = 0
		}
		rooms[i] = room
	}
	hasMoreRooms := true
	for level := 1; hasMoreRooms; level++ {
		hasMoreRooms = false
		for i := 0; i < len(rooms); i++ {
			room := rooms[i]
			r := rect{room.r.x, room.r.y, room.r.w - 1, room.r.h - 1}
			for j := 0; j < len(rooms) && room.level < 0; j++ {
				roomOther := rooms[j]
				if roomOther.level != level-1 {
					continue
				}
				rOther := rect{roomOther.r.x, roomOther.r.y, roomOther.r.w - 1, roomOther.r.h - 1}
				if r.IsAdjacent(rOther, overlapSize) {
					room.level = level
					rooms[i] = room
					hasMoreRooms = true
				}
			}
		}
	}

	// For every room, connect it to a random room with lower depth
	for i := 0; i < len(rooms); i++ {
		room := rooms[i]
		r := rect{room.r.x, room.r.y, room.r.w - 1, room.r.h - 1}
		for j := 0; j < len(rooms); j++ {
			roomOther := rooms[j]
			if roomOther.level != room.level-1 {
				continue
			}
			rOther := rect{roomOther.r.x, roomOther.r.y, roomOther.r.w - 1, roomOther.r.h - 1}
			if !r.IsAdjacent(rOther, overlapSize) {
				continue
			}
			// Rooms are adjacent; pick the cell that's in the middle of the
			// adjacent area and turn into a door
			minOverlapX := imin(
				room.r.x+room.r.w, roomOther.r.x+roomOther.r.w)
			maxOverlapX := imax(room.r.x, roomOther.r.x)
			minOverlapY := imin(
				room.r.y+room.r.h, roomOther.r.y+roomOther.r.h)
			maxOverlapY := imax(room.r.y, roomOther.r.y)
			overlapX := (minOverlapX + maxOverlapX) / 2
			overlapY := (minOverlapY + maxOverlapY) / 2
			g.setTile(overlapX, overlapY, room2)
			s.setTile(overlapX, overlapY, door)
			break
		}
	}

	return m
}
