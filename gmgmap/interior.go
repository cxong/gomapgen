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

  // Fill the map with room tiles surrounded by wall
  // We'll place the "road" tiles later
  g := m.Layer("Ground")
  g.fill(room)
  s := m.Layer("Structures")
  s.rectangleUnfilled(rect{0, 0, width, height}, wall)

  // Place the lobby
  var lobby rect
  for {
    x := rand.Intn(width - minRoomSize)
    y := rand.Intn(height - minRoomSize)
    lobby = rect{
      x, y,
      rand.Intn(maxRoomSize - minRoomSize) + minRoomSize,
      rand.Intn(maxRoomSize - minRoomSize) + minRoomSize}
    if lobby.x + lobby.w >= width {
      lobby.w = width - lobby.x
    }
    if lobby.y + lobby.h >= height {
      lobby.h = height - lobby.y
    }
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
  s.rectangleUnfilled(lobby, wall)

  return m
}
