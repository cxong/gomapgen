package gmgmap

import "math/rand"

// NewShop - create a single shop, surrounded by road tiles.
// A shop contains the following elements:
// - Road around grass, with floor interior
// - Walls
// - Counter area
// - On walls / surroundings:
//   - shop sign
//   - notice board
//   - windows/candles/shelves
// - Counter (at least 2 wide)
//   - shopkeep
// - shelves (up to 3 width aisles)
//   - items on shelves
// - rest area (if at least 3 wide free)
//   - rug
//   - tables/chairs (randomly placed, at least one chair per table)
// - against walls (display items, pots, barrels, leave diagonals free)
// - assistants (1 per 100 tiles, after the first)
// - patrons (1 per 36 tiles)
func NewShop(width, height int) *Map {
	m := NewMap(width, height)

	// Grass with road surroundings
	g := m.Layer("Ground")
	g.fill(grass)
	g.rectangle(rect{0, 0, g.Width, g.Height}, road, false)
	// Shop floor
	g.rectangle(rect{2, 2, g.Width - 4, g.Height - 5}, room, true)

	s := m.Layer("Structures")
	// Shop walls
	// Leave one row along bottom as buffer for lawn/board
	s.rectangle(rect{1, 1, g.Width - 2, g.Height - 3}, wall, false)
	// Entrance - connect with road/floor, replace wall with door, add sign
	entranceX := m.Width / 2
	g.setTile(entranceX, g.Height-2, road)
	g.setTile(entranceX, g.Height-3, room)
	doorY := s.Height - 3
	s.setTile(entranceX, doorY, door)
	s.setTile(entranceX+1, s.Height-2, sign)

	f := m.Layer("Furniture")
	// Items on walls - windows, candles/shelves inside, shop sign
	f.rectangle(rect{2, 1, f.Width - 4, 1}, hanging, false)
	f.setTile(entranceX-1, f.Height-3, sign)
	// Front wall elements must have 1 tile gap
	y := f.Height - 3
	for x := 3; x < f.Width-3; x = x + 2 {
		// Make sure that the area is clear of signs or doors
		if f.isClear(x-1, y, 3, 1) &&
			s.getTile(x-1, y) != door &&
			s.getTile(x, y) != door &&
			s.getTile(x+1, y) != door {
			f.setTile(x, y, window)
		}
	}

	// Counter - place opposite the door, random width (at least 2)
	counterW := rand.Intn(f.Width-4-2) + 2
	counterX := entranceX - counterW/2
	counterY := 3
	for x := counterX; x < counterX+counterW; x++ {
		f.setTile(x, counterY, counter)
	}
	c := m.Layer("Characters")
	// Shopkeep, opposite door
	c.setTile(entranceX, 2, shopkeeper)

	// Shelf area - to the right, at least 3 wide
	// Note: use two rands to get a distribution near the middle
	shelfX := rand.Intn(f.Width-7)/2 + (rand.Intn(f.Width-7)+1)/2 + 2
	shelfW := f.Width - 2 - shelfX
	// If wider than 5, leave at least 3 free to the left for rest area
	if shelfW > 5 && shelfX < 5 {
		shelfX = 5
		shelfW = f.Width - 2 - shelfX
	}
	// Place rows of shelves
	v := m.Layer("Inventory")
	var shelfClear = func(x, y int) bool {
		for x1 := x - 1; x1 <= x+1; x1++ {
			for y1 := y - 1; y1 <= y+1; y1++ {
				ftile := f.getTile(x1, y1)
				if ftile != nothing && ftile != shelf {
					return false
				}
			}
		}
		return true
	}
	for y := 3; y < f.Height-4; y = y + 2 {
		rowCounter := 0
		for x := shelfX; x < f.Width-3; x++ {
			if shelfClear(x, y) && rowCounter < 3 && x != entranceX {
				f.setTile(x, y, shelf)
				// Randomly place items on them
				if rand.Intn(3) < 2 {
					v.setTile(x, y, stock)
				}
				rowCounter++
			} else {
				rowCounter = 0
			}
		}
	}

	// Rest area
	if shelfX >= 5 {
		restRect := rect{2, 2, shelfX - 2, f.Height - 5}
		// Rug, in front of counter
		s.rectangle(rect{restRect.x, restRect.y + 2, restRect.w, restRect.h - 2},
			rug, true)

		// Randomly place tables from the wall to shelfX
		restArea := restRect.w * restRect.h
		for i := 0; i < 2*restArea; i++ {
			x := rand.Intn(restRect.w) + 2
			y := rand.Intn(restRect.h) + 2
			// Don't place in path of entrance
			if x == entranceX {
				continue
			}
			// Check that radius 1 is free of furniture
			if f.isClear(x-1, y-1, 3, 3) {
				f.setTile(x, y, table)
				// Place chairs as well, as long as the tiles behind it are free
				if x-1 >= 2 && f.isClear(x-2, y-1, 1, 3) {
					f.setTile(x-1, y, chair)
				}
				if x+1 < shelfX && f.isClear(x+2, y-1, 1, 3) {
					f.setTile(x+1, y, chair)
				}
			}
		}
	}

	// Items against walls - pots
	for i := 0; i < (f.Width+f.Height)*4; i++ {
		x := rand.Intn(f.Width-4) + 2
		y := rand.Intn(f.Height-5) + 2
		if x != 2 && x != f.Width-3 && y != 2 && y != f.Height-4 {
			continue
		}
		// Check that the 1 radius around is free of furniture that is not pots or
		// counters or hangings, and shopkeepers
		clear := c.isClear(x-1, y-1, 3, 3)
		for x1 := x - 1; x1 <= x+1 && clear; x1++ {
			for y1 := y - 1; y1 <= y+1 && clear; y1++ {
				furniture := f.getTile(x1, y1)
				if furniture != nothing && furniture != pot && furniture != counter &&
					furniture != hanging {
					clear = false
				}
			}
		}
		if clear {
			f.setTile(x, y, pot)
		}
	}

	// Shop assistants
	for i := 0; i < c.Width*c.Height/100-1; i++ {
		// Place assistants in the shop, in front of the counter
		for {
			x := rand.Intn(f.Width-4) + 2
			y := rand.Intn(f.Height-7) + 4
			if f.isClear(x, y, 1, 1) {
				c.setTile(x, y, assistant)
				break
			}
		}
	}

	// patrons - place anywhere except behind the counter
	for i := 0; i < c.Width*c.Height/36; i++ {
		for {
			x := rand.Intn(f.Width)
			y := rand.Intn(f.Height)
			if !(y == 2 && x >= counterX && x < counterX+counterW) &&
				// Allow patrons on rug
				(s.isClear(x, y, 1, 1) || s.getTile(x, y) == rug) &&
				// Allow patrons on chairs
				(f.isClear(x, y, 1, 1) || f.getTile(x, y) == chair) &&
				c.isClear(x, y, 1, 1) {
				c.setTile(x, y, player)
				break
			}
		}
	}

	return m
}
