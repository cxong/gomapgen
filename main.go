package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/cxong/gomapgen/gmgmap"
)

func main() {
	algo := flag.String("algo", "bspinterior", "generation algorithm: bsp/bspinterior/cell/rogue/shop/walk")
	template := flag.String("template", "dawnlike", "TMX export template: dawnlike/kenney")
	width := flag.Int("width", 32, "map width")
	height := flag.Int("height", 32, "map height")
	export := flag.Bool("export", true, "enable TMX export")
	iterations := flag.Int("iterations", 3000, "number of iterations for walk algo")
	gridWidth := flag.Int("gridwidth", 3, "grid size, for rogue algo")
	gridHeight := flag.Int("gridheight", 3, "grid size, for rogue algo")
	minRoomPct := flag.Int("minroompct", 50, "percent of rooms per grid, for rogue algo")
	maxRoomPct := flag.Int("maxroompct", 100, "percent of rooms per grid, for rogue algo")
	fillPct := flag.Int("fillpct", 40, "initial fill percent, for cell algo")
	r1 := flag.Int("r1", 5, "R1 cutoff, for cell algo")
	r2 := flag.Int("r2", 2, "R2 cutoff, for cell algo")
	reps := flag.Int("reps", 4, "reps, for cell algo")
	splits := flag.Int("splits", 4, "number of splits for bsp/bspinterior algo")
	minRoomSize := flag.Int("minroomsize", 5, "minimum room width/height")
	maxRoomSize := flag.Int("maxroomsize", 10, "maximum room width/height")
	connectionIterations := flag.Int(
		"connectioniterations", 15, "iterations for connection phase for bsp algo")
	lobbyEdgeType := flag.Int(
		"lobbyedge", gmgmap.LobbyEdge,
		"lobby placement for interior algo; 0=edge, 1=interior, 2=any")
	buildingPadding := flag.Int(
		"buildingPadding", 1, "padding between village buildings")
	corridorWidth := flag.Int(
		"corridorWidth", 1, "width of corridors (bspinterior only)")
	seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "random seed")
	flag.Parse()
	// make map
	fmt.Println("Using seed", *seed)
	rand.Seed(*seed)
	m := gmgmap.NewMap(*width, *height)
	switch *algo {
	case "bsp":
		m = gmgmap.NewBSP(*width, *height, *splits, *minRoomSize, *connectionIterations)
	case "bspinterior":
		m = gmgmap.NewBSPInterior(*width, *height, *splits, *minRoomSize, *corridorWidth)
	case "cell":
		m = gmgmap.NewCellularAutomata(*width, *height, *fillPct, *reps, *r1, *r2)
	case "interior":
		m = gmgmap.NewInterior(
			*width, *height, *minRoomSize, *maxRoomSize, *lobbyEdgeType)
	case "rogue":
		m = gmgmap.NewRogue(*width, *height, *gridWidth, *gridHeight,
			*minRoomPct, *maxRoomPct)
	case "shop":
		m = gmgmap.NewShop(*width, *height)
	case "walk":
		m = gmgmap.NewRandomWalk(*width, *height, *iterations)
	case "village":
		m = gmgmap.NewVillage(*width, *height, *buildingPadding)
	}
	// print
	m.Print()
	//m.PrintCSV()
	// export TMX
	if *export {
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
}
