package gmgmap

import (
	"math"
	"math/rand"
)

// NewBSPInterior - Create new BSP interior map
func NewBSPInterior(width, height, minRoomSize int) *Map {
	iterations := 4
	m := NewMap(width, height)

	// Split the map for a number of iterations, choosing alternating axis and random location
	var areas []bspRoom
	hcount := rand.Intn(2)
	areas = append(areas, bspRoomRoot(width, height))
	for i := 0; i < len(areas); i++ {
		if areas[i].level == iterations {
			break
		}
		var r1, r2 bspRoom
		var err error
		horizontal := ((hcount + int(math.Log2(float64(i)))) % 2) == 1
		if horizontal {
			r1, r2, err = bspSplitHorizontal(&areas[i], i, minRoomSize)
		} else {
			r1, r2, err = bspSplitVertical(&areas[i], i, minRoomSize)
		}
		if err == nil {
			// Resize rooms to allow space for street
			if horizontal {
				r1.r.w--
			} else {
				r1.r.h--
			}
			areas[i].child1 = len(areas)
			areas = append(areas, r1)
			areas[i].child2 = len(areas)
			areas = append(areas, r2)
		}
	}

	g := m.Layer("Ground")
	s := m.Layer("Structures")
	// Fill with street
	g.fill(room2)
	// Turn the leaves into rooms
	for i := range areas {
		// Only place rooms in leaf nodes
		if areas[i].child1 >= 0 || areas[i].child2 >= 0 {
			continue
		}
		var r rect
		r.w = areas[i].r.w
		r.x = areas[i].r.x
		r.h = areas[i].r.h
		r.y = areas[i].r.y
		g.rectangleFilled(rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}, room)
		s.rectangleUnfilled(r, wall2)
	}

	return m
}
