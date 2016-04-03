package gmgmap

import "math/rand"

const (
  LOBBY_EDGE = iota
  LOBBY_INTERIOR
  LOBBY_ANY
)

// NewInterior - create a building interior layout map.
// The layout has a "lobby", multiple rooms and corridors.
// The lobby must be at most 2 steps from any room.
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
    if r1, r2, err := bspSplit(&rooms[i], i, minRoomSize, maxRoomSize); err == nil {
      rooms[i].child1 = len(rooms)
      rooms = append(rooms, r1)
      rooms[i].child2 = len(rooms)
      rooms = append(rooms, r2)
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

  // Form the room walls using the leaf rooms
  for i := 0; i < len(rooms); i++ {
    if rooms[i].child1 >= 0 || rooms[i].child2 >= 0 {
      continue
    }
    r := rooms[i].r
    s.rectangleUnfilled(r, wall2)
    groundRect := rect{r.x+1, r.y+1, r.w-2, r.h-2}
    g.rectangleFilled(groundRect, room)
  }

  // Choose one of the rooms to be the lobby
  var lobby rect
  var roomIndices = rand.Perm(len(rooms))
  for i := 0; i < len(roomIndices); i++ {
    room := rooms[roomIndices[i]]
    // Only choose leaves
    if room.child1 >= 0 || room.child2 >= 0 {
      continue
    }
    lobby = room.r
    // Check if the lobby placement is ok;
    // If it's on the edge or in the interior
    if lobbyEdge == LOBBY_EDGE {
      if lobby.x == 0 || lobby.y == 0 ||
        lobby.x + lobby.w == width || lobby.y + lobby.h == height {
        break
      }
    } else if lobbyEdge == LOBBY_INTERIOR {
      if lobby.x > 0 && lobby.y > 0 &&
        lobby.x + lobby.w < width && lobby.y + lobby.h < height {
        break
      }
    } else {
      break
    }
  }
  // Place the lobby
  lobby.x++
  lobby.y++
  lobby.w -= 2
  lobby.h -= 2
  g.rectangleFilled(lobby, room2)

  return m
}
