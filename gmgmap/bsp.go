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

const minRoomSize = 4

// NewBSP - generate a new dungeon, using BSP method
func NewBSP(width, height, iterations int) *Map {
	m := NewMap(width, height)

	// Split the map for a number of iterations, choosing random axis and location
	var areas []bspRoom
	areas = append(areas, bspRoom{rect{0, 0, width, height}, -1, -1, -1, 0})
	for i := 0; i < len(areas); i++ {
		if areas[i].level == iterations {
			break
		}
		if r1, r2, err := split(&areas[i], i); err == nil {
			areas[i].child1 = len(areas)
			areas = append(areas, r1)
			areas[i].child2 = len(areas)
			areas = append(areas, r2)
		}
	}

	g := m.Layer("Ground")
	s := m.Layer("Structures")
	rooms := make([]bspRoom, len(areas))
	// Place rooms randomly into the split areas
	for i := range areas {
		rooms[i] = areas[i]
		if areas[i].child1 < 0 && areas[i].child2 < 0 {
			var r rect
			// make sure center of area is inside room
			if areas[i].r.w == minRoomSize {
				r.w = minRoomSize
				r.x = areas[i].r.x
			} else {
				r.w = rand.Intn(areas[i].r.w-minRoomSize) + minRoomSize
				r.x = rand.Intn(areas[i].r.w-r.w) + areas[i].r.x
				xmid := areas[i].r.x + areas[i].r.w/2
				r.x = iclamp(r.x, xmid-(r.w-2), xmid-1)
			}
			if areas[i].r.h == minRoomSize {
				r.h = minRoomSize
				r.y = areas[i].r.y
			} else {
				r.h = rand.Intn(areas[i].r.h-minRoomSize) + minRoomSize
				r.y = rand.Intn(areas[i].r.h-r.h) + areas[i].r.y
				ymid := areas[i].r.y + areas[i].r.h/2
				r.y = iclamp(r.y, ymid-(r.h-2), ymid-1)
			}
			g.rectangle(rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}, room, true)
			s.rectangle(r, wall2, false)
			rm := areas[i]
			rm.r = r
			rooms[i] = rm
		}
	}

	// Connect leaves to siblings
	for i := range rooms {
		r := rooms[i]
		if r.parent < 0 || r.child1 >= 0 || r.child2 >= 0 {
			continue
		}
		// Only connect to second sibling, so we don't double up
		// Also needs to have siblings
		siblingIndex := rooms[rooms[i].parent].child2
		if siblingIndex == i || siblingIndex < 0 {
			continue
		}
		// Sibling also needs to be leaf
		sibling := rooms[siblingIndex]
		if sibling.child1 >= 0 || sibling.child2 >= 0 {
			continue
		}
		xmin := imax(r.r.x, sibling.r.x)
		xmax := imin(r.r.x+r.r.w, sibling.r.x+sibling.r.w)
		ymin := imax(r.r.y, sibling.r.y)
		ymax := imin(r.r.y+r.r.h, sibling.r.y+sibling.r.h)
		if xmin+1 < xmax-1 {
			// connect from top to bottom
			x := irand(xmin+1, xmax-1)
			addCorridor(g, s, x, ymax-1, x, ymin, 0, 1, room2)
		} else if ymin+1 < ymax-1 {
			// Connect from left to right
			y := irand(ymin+1, ymax-1)
			addCorridor(g, s, xmax-1, y, xmin, y, 1, 0, room2)
		}
	}

	return m
}

func split(room *bspRoom, i int) (bspRoom, bspRoom, error) {
	// If more than 3:2, split the long dimension, otherwise randomise
	if room.r.w*3 > room.r.h*2 || (room.r.h*3 < room.r.w*2 && rand.Intn(2) == 0) {
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
