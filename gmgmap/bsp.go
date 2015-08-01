package gmgmap

import (
	"errors"
	"math/rand"
)

type bspRoom struct {
	r              rect
	parent         int
	child1, child2 int
	level          int
}

const minRoomSize = 3

// NewBSP - generate a new dungeon, using BSP method
func NewBSP(width, height, iterations int) *Map {
	m := NewMap(width, height)

	// Split the map for a number of iterations, choosing random axis and location
	var rooms []bspRoom
	rooms = append(rooms, bspRoom{rect{0, 0, width, height}, -1, -1, -1, 0})
	for i := 0; i < len(rooms); i++ {
		room := &rooms[i]
		if room.level == iterations-1 {
			break
		}
		if r1, r2, err := split(room, i); err == nil {
			room.child1 = len(rooms)
			rooms = append(rooms, r1)
			room.child2 = len(rooms)
			rooms = append(rooms, r2)
		}
	}

	g := m.Layer("Ground")
	s := m.Layer("Structures")
	// Place rooms randomly into the split areas
	for i := range rooms {
		if rooms[i].child1 < 0 && rooms[i].child2 < 0 {
			var r rect
			if rooms[i].r.w == minRoomSize {
				r.w = minRoomSize
				r.x = rooms[i].r.x
			} else {
				r.w = rand.Intn(rooms[i].r.w-minRoomSize) + minRoomSize
				r.x = rand.Intn(rooms[i].r.w-r.w) + rooms[i].r.x
			}
			if rooms[i].r.h == minRoomSize {
				r.h = minRoomSize
				r.y = rooms[i].r.y
			} else {
				r.h = rand.Intn(rooms[i].r.h-minRoomSize) + minRoomSize
				r.y = rand.Intn(rooms[i].r.h-r.h) + rooms[i].r.y
			}
			g.rectangle(rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}, room, true)
			s.rectangle(r, wall2, false)
		}
	}
	return m
}

func split(room *bspRoom, i int) (bspRoom, bspRoom, error) {
	if rand.Intn(2) == 0 {
		// Split horizontally
		r := room.r.w - minRoomSize*2
		var x int
		if r < 0 {
			return bspRoom{}, bspRoom{}, errors.New("room too small")
		} else if r == 0 {
			x = minRoomSize
		} else {
			x = rand.Intn(r) + minRoomSize
		}
		return bspRoom{rect{room.r.x, room.r.y, x, room.r.h}, i, -1, -1, room.level + 1},
			bspRoom{rect{room.r.x + x, room.r.y, room.r.w - x, room.r.h}, i, -1, -1, room.level + 1},
			nil
	}
	// Split vertically
	r := room.r.h - minRoomSize*2
	var y int
	if r < 0 {
		return bspRoom{}, bspRoom{}, errors.New("room too small")
	} else if r == 0 {
		y = minRoomSize
	} else {
		y = rand.Intn(r) + minRoomSize
	}
	return bspRoom{rect{room.r.x, room.r.y, room.r.w, y}, i, -1, -1, room.level + 1},
		bspRoom{rect{room.r.x, room.r.y + y, room.r.w, room.r.h - y}, i, -1, -1, room.level + 1},
		nil
}
