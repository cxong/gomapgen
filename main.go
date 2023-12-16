package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cxong/gomapgen/gmgmap"
)

func main() {
	algo := flag.String("algo", "bspinterior", "generation algorithm: bsp/bspinterior/cell/rogue/shop/wfcshop/walk/village")
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
	rr := rand.New(rand.NewSource(*seed))
	// Use different RNG for export
	rr2 := rand.New(rand.NewSource(*seed))
	m := gmgmap.NewMap(*width, *height)
	t := &gmgmap.DawnLikeTemplate
	switch *template {
	case "dawnlike":
		t = &gmgmap.DawnLikeTemplate
	case "kenney":
		t = &gmgmap.KenneyTemplate
	}

	// Iteratively generate and export map
	imgId := 0
	exportFunc := func(m_ *gmgmap.Map) {
	}
	if *export {
		// Remove existing images
		files, err := filepath.Glob("tmx_export/map*.png")
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				panic(err)
			}
		}
		exportFunc = func(m_ *gmgmap.Map) {
			if err := m_.ToTMX(rr2, t, imgId); err != nil {
				panic(err)
			}
			imgId++
		}
	}

	switch *algo {
	case "bsp":
		m = gmgmap.NewBSP(rr, *width, *height, *splits, *minRoomSize, *connectionIterations)
	case "bspinterior":
		m = gmgmap.NewBSPInterior(rr, exportFunc, *width, *height, *splits, *minRoomSize, *corridorWidth)
	case "cell":
		m = gmgmap.NewCellularAutomata(rr, *width, *height, *fillPct, *reps, *r1, *r2)
	case "interior":
		m = gmgmap.NewInterior(
			rr, *width, *height, *minRoomSize, *maxRoomSize, *lobbyEdgeType)
	case "rogue":
		m = gmgmap.NewRogue(rr, *width, *height, *gridWidth, *gridHeight,
			*minRoomPct, *maxRoomPct)
	case "shop":
		m = gmgmap.NewShop(rr, exportFunc, *width, *height)
	case "walk":
		m = gmgmap.NewRandomWalk(rr, *width, *height, *iterations)
	case "wfcshop":
		m = gmgmap.NewWFCShop(rr, exportFunc, *width, *height)
	case "village":
		m = gmgmap.NewVillage(rr, exportFunc, *width, *height, *buildingPadding)
	}

	// print
	m.Print()
	//m.PrintCSV()
	// export TMX
	exportFunc(m)
	// export gif
	if *export {
		cmd := exec.Command("convert", "-delay", strconv.Itoa(2000/(imgId+1)), "-dispose", "previous", "tmx_export/map*.png",
			// make the last frame last longer
			"-delay", "200", fmt.Sprintf("tmx_export/map%04d.png", imgId-1), "tmx_export/map.gif")
		_, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
