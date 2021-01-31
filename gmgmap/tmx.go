package gmgmap

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// DTO for CSV export
type csvExport struct {
	Name   string
	Width  int
	Height int
	Values string
}

// TMXTemplate - configuration for TMX export
type TMXTemplate struct {
	path       string
	background string
	// Arrays of tile ids (16);
	// first is centre,
	// then 8 tiles from top clockwise,
	// then h/v,
	// then 4 end tiles from top clockwise,
	// then isolated tile
	floorIDs [16]string
	//floor2IDs  [16]string
	roadIDs     [16]string
	road2IDs    [16]string
	wallIDs     [16]string
	wall2IDs    [16]string
	roomIDs     [16]string
	room2IDs    [16]string
	doorH       string
	doorV       string
	doorLockedH string
	doorLockedV string
	stairsUp    string
	stairsDown  string
	// note: trees use different tiling system
	// centre
	// 8 tiles from top clockwise
	// 4 concaves from upper right clockwise
	// upleft/downright diagonal
	// upright/downleft diagonal
	// isolated
	treeIDs  [16]string
	grassIDs [16]string

	// Flavour tiles - randomly chosen
	signIDs       []string
	wallSignIDs   []string // signs that hang off walls TODO: by shop type
	hangingIDs    []string
	windowIDs     []string
	counterHIDs   [3]string // 3 IDs, left middle end
	counterVIDs   [3]string // 3 IDs, top middle bottom
	shopkeeperIDs []string
	shelfID       string    // TODO: per item type
	stockIDs      []string  // TODO: per item type
	tableID       string    // TODO: furniture sets
	chairIDs      [2]string // left/right TODO: furniture sets
	rugIDs        [16]string
	potIDs        []string
	assistantIDs  []string
	playerIDs     []string
	flowerIDs     []string
	keyIDs        []string

	// Parameters used for template export
	Width  int
	Height int
	CSVs   []csvExport
}

// ToTMX - export map as TMX (Tiled XML map)
func (m Map) ToTMX(tmxTemplate *TMXTemplate) error {
	exportDir := "tmx_export"
	err := os.Mkdir(exportDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
		return err
	}
	// Copy data files
	baseDir := path.Join("gmgmap", tmxTemplate.path)
	err = filepath.Walk(baseDir, func(walkPath string, info os.FileInfo, _ error) error {
		destDir := walkPath[len(baseDir):]
		if destDir == "" {
			return nil
		}
		destPath := path.Join(exportDir, destDir)
		if info.IsDir() {
			// Make dir if not exists
			if err := os.Mkdir(destPath, 0755); err != nil && !os.IsExist(err) {
				return err
			}
			return nil
		}
		// Copy file, except for tmx (which we'll be generating)
		if strings.ToLower(filepath.Ext(walkPath)) == ".tmx" {
			return nil
		}
		//fmt.Printf("Copying %s to %s\n", walkPath, destPath)
		src, err := os.Open(walkPath)
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.Create(destPath)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			return err
		}
		return dst.Close()
	})
	if err != nil {
		return err
	}

	populateTemplate(m, tmxTemplate)

	// Generate TMX
	// Use template path as template name
	t, err := template.ParseFiles(path.Join(baseDir, "template.tmx"))
	if err != nil {
		return err
	}
	templateFile, err := os.Create(path.Join(exportDir, "map.tmx"))
	if err != nil {
		return err
	}
	if err := t.Execute(templateFile, tmxTemplate); err != nil {
		return err
	}
	return nil
}

func populateTemplate(m Map, tmp *TMXTemplate) {
	tmp.Width = m.Width
	tmp.Height = m.Height
	var arrayToCSV = func(xt []string, w, h int) string {
		var xtline []string
		for y := 0; y < h; y++ {
			xtline = append(xtline, strings.Join(xt[y*w:(y+1)*w], ","))
		}
		return strings.Join(xtline, ",\n")
	}
	var makeCSV = func(l *Layer, wallLayer *Layer) csvExport {
		xt := make([]string, l.Width*l.Height)
		for y := 0; y < l.Height; y++ {
			for x := 0; x < l.Width; x++ {
				tile := l.getTile(x, y)
				var tileIDs *[16]string
				switch tile {
				case nothing:
					xt[x+y*l.Width] = "0"
				case floor:
					tileIDs = &tmp.floorIDs
				//case floor2:
				//tileIDs = &tmp.floor2IDs
				case road:
					tileIDs = &tmp.roadIDs
				case road2:
					tileIDs = &tmp.road2IDs
				case wall:
					tileIDs = &tmp.wallIDs
				case wall2:
					tileIDs = &tmp.wall2IDs
				case room:
					tileIDs = &tmp.roomIDs
				case room2:
					tileIDs = &tmp.room2IDs
				case door:
					left := wall
					if x > 0 {
						left = wallLayer.getTile(x-1, y)
					}
					if IsWall(left) || IsDoor(left) {
						xt[x+y*l.Width] = tmp.doorH
					} else {
						xt[x+y*l.Width] = tmp.doorV
					}
				case doorLocked:
					left := wall
					if x > 0 {
						left = wallLayer.getTile(x-1, y)
					}
					if IsWall(left) || IsDoor(left) {
						xt[x+y*l.Width] = tmp.doorLockedH
					} else {
						xt[x+y*l.Width] = tmp.doorLockedV
					}
				case stairsUp:
					xt[x+y*l.Width] = tmp.stairsUp
				case stairsDown:
					xt[x+y*l.Width] = tmp.stairsDown
				case tree:
					xt[x+y*l.Width] = get16Tile2(m, x, y, tile, &tmp.treeIDs)
				case grass:
					tileIDs = &tmp.grassIDs
				case sign:
					// choose from on-wall sign or stand-alone sign
					if IsWall(wallLayer.getTile(x, y)) {
						xt[x+y*l.Width] = tmp.wallSignIDs[rand.Intn(len(tmp.wallSignIDs))]
					} else {
						xt[x+y*l.Width] = tmp.signIDs[rand.Intn(len(tmp.signIDs))]
					}
				case hanging:
					xt[x+y*l.Width] = tmp.hangingIDs[rand.Intn(len(tmp.hangingIDs))]
				case window:
					xt[x+y*l.Width] = tmp.windowIDs[rand.Intn(len(tmp.windowIDs))]
				case counter:
					top := y > 0 && l.getTile(x, y-1) == counter
					bottom := y < l.Height-1 && l.getTile(x, y+1) == counter
					left := x > 0 && l.getTile(x-1, y) == counter
					right := x < l.Width-1 && l.getTile(x+1, y) == counter
					switch {
					case left && right:
						xt[x+y*l.Width] = tmp.counterHIDs[1]
					case top && bottom:
						xt[x+y*l.Width] = tmp.counterVIDs[1]
					case left:
						xt[x+y*l.Width] = tmp.counterHIDs[2]
					case right:
						xt[x+y*l.Width] = tmp.counterHIDs[0]
					case top:
						xt[x+y*l.Width] = tmp.counterVIDs[2]
					case bottom:
						xt[x+y*l.Width] = tmp.counterVIDs[0]
					}
				case shopkeeper:
					xt[x+y*l.Width] = tmp.shopkeeperIDs[rand.Intn(len(tmp.shopkeeperIDs))]
				case shelf:
					xt[x+y*l.Width] = tmp.shelfID
				case stock:
					xt[x+y*l.Width] = tmp.stockIDs[rand.Intn(len(tmp.stockIDs))]
				case table:
					xt[x+y*l.Width] = tmp.tableID
				case chair:
					// Find the table to face
					if x == 0 || l.getTile(x-1, y) == table {
						xt[x+y*l.Width] = tmp.chairIDs[1]
					} else {
						xt[x+y*l.Width] = tmp.chairIDs[0]
					}
				case rug:
					tileIDs = &tmp.rugIDs
				case pot:
					xt[x+y*l.Width] = tmp.potIDs[rand.Intn(len(tmp.potIDs))]
				case assistant:
					xt[x+y*l.Width] = tmp.assistantIDs[rand.Intn(len(tmp.assistantIDs))]
				case player:
					xt[x+y*l.Width] = tmp.playerIDs[rand.Intn(len(tmp.playerIDs))]
				case flower:
					xt[x+y*l.Width] = tmp.flowerIDs[rand.Intn(len(tmp.flowerIDs))]
				case key:
					xt[x+y*l.Width] = tmp.keyIDs[rand.Intn(len(tmp.keyIDs))]
				default:
					fmt.Println("Unhandled tile", tile)
					panic(tile)
				}
				if tileIDs != nil {
					xt[x+y*l.Width] = get16Tile(m, x, y, tile, tileIDs)
				}
			}
		}
		return csvExport{l.Name, l.Width, l.Height,
			arrayToCSV(xt, l.Width, l.Height)}
	}
	// Add a background layer export for appearance
	backArr := make([]string, m.Width*m.Height)
	for i := 0; i < len(backArr); i++ {
		backArr[i] = tmp.background
	}
	tmp.CSVs = append(tmp.CSVs,
		csvExport{"Background", m.Width, m.Height,
			arrayToCSV(backArr, m.Width, m.Height)})
	for _, l := range m.Layers {
		tmp.CSVs = append(tmp.CSVs, makeCSV(l, m.Layer("Structures")))
	}
}

func get16Tile(m Map, x, y int, tile rune, templateTiles *[16]string) string {
	up := hasSameTile(m, x, y-1, tile)
	right := hasSameTile(m, x+1, y, tile)
	down := hasSameTile(m, x, y+1, tile)
	left := hasSameTile(m, x-1, y, tile)
	switch {
	case up && right && down && left:
		return templateTiles[0]
	case !up && right && down && left:
		// upper edge
		return templateTiles[1]
	case !up && !right && down && left:
		// upper right corner
		return templateTiles[2]
	case up && !right && down && left:
		// right edge
		return templateTiles[3]
	case up && !right && !down && left:
		// bottom right corner
		return templateTiles[4]
	case up && right && !down && left:
		// bottom edge
		return templateTiles[5]
	case up && right && !down && !left:
		// bottom left corner
		return templateTiles[6]
	case up && right && down && !left:
		// left edge
		return templateTiles[7]
	case !up && right && down && !left:
		// upper left corner
		return templateTiles[8]
	case !up && right && !down && left:
		// horizontal
		return templateTiles[9]
	case up && !right && down && !left:
		// vertical
		return templateTiles[10]
	case !up && !right && down && !left:
		// upper end
		return templateTiles[11]
	case !up && !right && !down && left:
		// right end
		return templateTiles[12]
	case up && !right && !down && !left:
		// bottom end
		return templateTiles[13]
	case !up && right && !down && !left:
		// left end
		return templateTiles[14]
	case !up && !right && !down && !left:
		// isolated
		return templateTiles[15]
	}
	panic("unknown error")
}

func get16Tile2(m Map, x, y int, tile rune, templateTiles *[16]string) string {
	up := hasSameTile(m, x, y-1, tile)
	upright := hasSameTile(m, x+1, y-1, tile)
	right := hasSameTile(m, x+1, y, tile)
	downright := hasSameTile(m, x+1, y+1, tile)
	down := hasSameTile(m, x, y+1, tile)
	downleft := hasSameTile(m, x-1, y+1, tile)
	left := hasSameTile(m, x-1, y, tile)
	upleft := hasSameTile(m, x-1, y-1, tile)
	switch {
	case up && upright && right && downright && down && downleft && left && upleft:
		return templateTiles[0]
	case up && right && downright && down && downleft && left && upleft:
		// upper right concave
		return templateTiles[9]
	case up && upright && right && down && downleft && left && upleft:
		// lower right concave
		return templateTiles[10]
	case up && upright && right && downright && down && left && upleft:
		// lower left concave
		return templateTiles[11]
	case up && upright && right && downright && down && downleft && left:
		// upper left concave
		return templateTiles[12]
	case right && downright && down && left && upleft && up:
		// upleft-downright diagonal
		return templateTiles[13]
	case up && upright && right && down && downleft && left:
		// upright-downleft diagonal
		return templateTiles[14]
	case right && downright && down && downleft && left:
		// upper edge
		return templateTiles[1]
	case up && down && downleft && left && upleft:
		// right edge
		return templateTiles[3]
	case up && upright && right && left && upleft:
		// bottom edge
		return templateTiles[5]
	case up && upright && right && downright && down:
		// left edge
		return templateTiles[7]
	case down && downleft && left:
		// upper right corner
		return templateTiles[2]
	case up && left && upleft:
		// bottom right corner
		return templateTiles[4]
	case up && upright && right:
		// bottom left corner
		return templateTiles[6]
	case right && downright && down:
		// upper left corner
		return templateTiles[8]
	}
	return templateTiles[15]
}

func hasSameTile(m Map, x, y int, tile rune) bool {
	// Walls don't extend to edge
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return !IsWall(tile)
	}
	for _, l := range m.Layers {
		t := l.getTile(x, y)
		if t == tile {
			return true
		} else if IsWall(tile) && (IsWall(t) || t == door) {
			return true
		}
	}
	return false
}
