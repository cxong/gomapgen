package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/cxong/gomapgen/gmgmap"
)

func main() {
	algo := flag.String("algo", "shop", "generation algorithm: bsp/cell/rogue/shop/walk")
	template := flag.String("template", "dawnlike", "TMX export template: dawnlike/kenney")
	width := flag.Int("width", 32, "map width")
	height := flag.Int("height", 32, "map height")
	iterations := flag.Int("iterations", 3000, "number of iterations for walk algo")
	gridWidth := flag.Int("gridwidth", 3, "grid size, for rogue algo")
	gridHeight := flag.Int("gridheight", 3, "grid size, for rogue algo")
	minRoomPct := flag.Int("minroompct", 50, "percent of rooms per grid, for rogue algo")
	maxRoomPct := flag.Int("maxroompct", 100, "percent of rooms per grid, for rogue algo")
	fillPct := flag.Int("fillpct", 40, "initial fill percent, for cell algo")
	r11 := flag.Int("r11", 5, "R1 cutoff rep 1, for cell algo")
	r12 := flag.Int("r12", 2, "R2 cutoff rep 1, for cell algo")
	reps1 := flag.Int("reps1", 4, "reps for rep 1, for cell algo")
	r21 := flag.Int("r21", 5, "R1 cutoff rep 2, for cell algo")
	r22 := flag.Int("r22", -1, "R2 cutoff rep 2, for cell algo")
	reps2 := flag.Int("reps2", 3, "reps for rep 2, for cell algo")
	splits := flag.Int("splits", 4, "number of splits for bsp algo")
	minRoomSize := flag.Int("minroomsize", 5, "minimum room width/height for bsp algo")
	connectionIterations := flag.Int("connectioniterations", 15, "iterations for connection phase for bsp algo")
	seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "random seed")
	flag.Parse()
	// make map
	fmt.Println("Using seed", *seed)
	rand.Seed(*seed)
	m := gmgmap.NewMap(*width, *height)
	switch *algo {
	case "bsp":
		m = gmgmap.NewBSP(*width, *height, *splits, *minRoomSize, *connectionIterations)
	case "cell":
		m = gmgmap.NewCellularAutomata(*width, *height, *fillPct,
			*r11, *r12, *reps1, *r21, *r22, *reps2)
	case "rogue":
		m = gmgmap.NewRogue(*width, *height, *gridWidth, *gridHeight,
			*minRoomPct, *maxRoomPct)
	case "shop":
		m = gmgmap.NewShop(*width, *height)
	case "walk":
		m = gmgmap.NewRandomWalk(*width, *height, *iterations)
	}
	// print
	m.Print()
	// export TMX
	t := &gmgmap.DawnLikeTemplate
	switch *template {
	case "dawnlike":
		t = &gmgmap.DawnLikeTemplate
	case "kenney":
		t = &gmgmap.KenneyTemplate
	}
	if err := m.ToTMX(t); err != nil {
		panic(err)
	}
}
