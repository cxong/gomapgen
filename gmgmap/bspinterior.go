package gmgmap

import (
	"math"
	"math/rand"
)

type street struct {
	r          rect
	horizontal bool
	level      int
}

// NewBSPInterior - Create new BSP interior map
// Implementation of https://gamedev.stackexchange.com/questions/47917/procedural-house-with-rooms-generator/48216#48216
func NewBSPInterior(width, height, minRoomSize int) *Map {
	iterations := 4
	corridorWidth := 2
	// corridorLevelDiffBlock := 2
	m := NewMap(width, height)

	// Split the map for a number of iterations, choosing alternating axis and random location
	var areas []bspRoom
	var streets []street
	hcount := rand.Intn(2)
	areas = append(areas, bspRoomRoot(width, height))
	for i := 0; i < len(areas); i++ {
		if areas[i].level == iterations {
			break
		}
		var r1, r2 bspRoom
		var err error = nil
		horizontal := ((hcount + int(math.Log2(float64(i)))) % 2) == 1
		if horizontal {
			r1, r2, err = bspSplitHorizontal(&areas[i], i, minRoomSize+corridorWidth/2)
		} else {
			r1, r2, err = bspSplitVertical(&areas[i], i, minRoomSize+corridorWidth/2)
		}
		if err == nil {
			// Resize rooms to allow space for street
			for j := 0; j < corridorWidth; j++ {
				if horizontal {
					if j%2 == 0 {
						r1.r.w--
					} else {
						r2.r.x++
						r2.r.w--
					}
				} else {
					if j%2 == 0 {
						r1.r.h--
					} else {
						r2.r.y++
						r2.r.h--
					}
				}
			}
			areas[i].child1 = len(areas)
			areas = append(areas, r1)
			areas[i].child2 = len(areas)
			areas = append(areas, r2)
			var s street
			if horizontal {
				s.r = rect{r1.r.x + r1.r.w, r1.r.y, corridorWidth, r1.r.h}
			} else {
				s.r = rect{r1.r.x, r1.r.y + r1.r.h, r1.r.w, corridorWidth}
			}
			s.level = r1.level
			s.horizontal = !horizontal
			streets = append(streets, s)
		}
	}

	g := m.Layer("Ground")
	s := m.Layer("Structures")
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
	// Fill streets
	for i := range streets {
		g.rectangleFilled(streets[i].r, room2)
		// Check ends of street - if next to much older street, block off with wall
		end1 := vec2{streets[i].r.x, streets[i].r.y}
		end2 := vec2{streets[i].r.x + streets[i].r.w - 1, streets[i].r.y + streets[i].r.h - 1}
		var dAlong, dAcross vec2
		if streets[i].horizontal {
			dAlong = vec2{1, 0}
			dAcross = vec2{0, 1}
		} else {
			dAlong = vec2{0, 1}
			dAcross = vec2{1, 0}
		}
		capStreet(g, s, end1, dAcross, dAlong, corridorWidth)
		capStreet(g, s, end2, vec2{-dAcross.x, -dAcross.y}, vec2{-dAlong.x, -dAlong.y}, corridorWidth)
	}

	return m
}

func capStreet(g, s *Layer, end, dAcross, dAlong vec2, corridorWidth int) {
	// Check ends of street - if outside map, or next to much older street, block off with wall
	if !g.isIn(end.x-dAlong.x, end.y-dAlong.y) {
		for i := 0; i < corridorWidth; i++ {
			g.setTile(end.x+dAcross.x*i, end.y+dAcross.y*i, nothing)
			s.setTile(end.x+dAcross.x*i, end.y+dAcross.y*i, wall2)
		}
	}
}
