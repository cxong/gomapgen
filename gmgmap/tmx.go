package gmgmap

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// TMXTemplate - configuration for TMX export
type TMXTemplate struct {
	path string
	// Arrays of tile ids (16);
	// first is centre,
	// then 8 tiles from top clockwise,
	// then h/v,
	// then 4 end tiles from top clockwise,
	// then isolated tile\
	floorIDs   []string
	floor2IDs  []string
	wallIDs    []string
	wall2IDs   []string
	roomIDs    []string
	room2IDs   []string
	doorH      string
	doorV      string
	stairsUp   string
	stairsDown string
	treeIDs    []string

	// Parameters for map generation
	floorTerrain  bool
	floor2Terrain bool
	wallTerrain   bool
	wall2Terrain  bool
	roomTerrain   bool
	room2Terrain  bool
	treeTerrain   bool

	// Parameters used for template export
	Width        int
	Height       int
	CSV          string
	CSVFurniture string
}

// DawnLikeTemplate - using DawnLike tile set
var DawnLikeTemplate = TMXTemplate{
	"dawnlike",
	[]string{"1421", "1400", "1401", "1422", "1443", "1442", "1441", "1420", "1399", "1425", "1423", "1402", "1426", "1444", "1424", "1404"},
	[]string{"1176", "1155", "1156", "1177", "1198", "1197", "1196", "1175", "1154", "1180", "1178", "1157", "1181", "1199", "1179", "1159"},
	[]string{"92", "72", "70", "93", "110", "112", "108", "91", "68", "69", "88", "88", "110", "89", "108", "71"},
	[]string{"85", "65", "63", "86", "103", "105", "101", "84", "61", "62", "81", "81", "103", "82", "101", "64"},
	[]string{"1428", "1407", "1408", "1429", "1450", "1449", "1448", "1427", "1406", "1432", "1430", "1409", "1433", "1451", "1431", "1411"},
	[]string{"1232", "1211", "1212", "1233", "1254", "1253", "1252", "1231", "1210", "1236", "1234", "1213", "1237", "1255", "1235", "1215"},
	"2096", "2097",
	"3304", "3305",
	[]string{"2537", "2525", "2526", "2538", "2550", "2549", "2548", "2536", "2524", "2540", "2528", "2529", "2541", "2553", "2552", "2527"},
	false, true, true, true, true, true, true,
	0, 0, "", ""}

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
		fmt.Printf("Copying %s to %s\n", walkPath, destPath)
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
	t.Execute(templateFile, tmxTemplate)
	return nil
}

func populateTemplate(m Map, tmxTemplate *TMXTemplate) {
	tmxTemplate.Width = m.Width
	tmxTemplate.Height = m.Height
	var makeCSV = func(l *Layer, tileLayer *Layer) string {
		exportTiles := make([]string, l.Width*l.Height)
		for y := 0; y < l.Height; y++ {
			for x := 0; x < l.Width; x++ {
				tile := l.getTile(x, y)
				var tileIDs []string
				switch tile {
				case nothing:
					exportTiles[x+y*l.Width] = "0"
					continue
				case floor:
					if !tmxTemplate.floorTerrain {
						exportTiles[x+y*l.Width] = tmxTemplate.floorIDs[0]
						continue
					}
					tileIDs = tmxTemplate.floorIDs
				case floor2:
					if !tmxTemplate.floor2Terrain {
						exportTiles[x+y*l.Width] = tmxTemplate.floor2IDs[0]
						continue
					}
					tileIDs = tmxTemplate.floor2IDs
				case wall:
					if !tmxTemplate.wallTerrain {
						exportTiles[x+y*l.Width] = tmxTemplate.wallIDs[0]
						continue
					}
					tileIDs = tmxTemplate.wallIDs
				case wall2:
					if !tmxTemplate.wall2Terrain {
						exportTiles[x+y*l.Width] = tmxTemplate.wall2IDs[0]
						continue
					}
					tileIDs = tmxTemplate.wall2IDs
				case room:
					if !tmxTemplate.roomTerrain {
						exportTiles[x+y*l.Width] = tmxTemplate.roomIDs[0]
						continue
					}
					tileIDs = tmxTemplate.roomIDs
				case room2:
					if !tmxTemplate.room2Terrain {
						exportTiles[x+y*l.Width] = tmxTemplate.room2IDs[0]
						continue
					}
					tileIDs = tmxTemplate.room2IDs
				case door:
					left := wall
					if x > 0 {
						left = tileLayer.getTile(x-1, y)
					}
					if IsWall(left) {
						exportTiles[x+y*l.Width] = tmxTemplate.doorH
					} else {
						exportTiles[x+y*l.Width] = tmxTemplate.doorV
					}
					continue
				case stairsUp:
					exportTiles[x+y*l.Width] = tmxTemplate.stairsUp
					continue
				case stairsDown:
					exportTiles[x+y*l.Width] = tmxTemplate.stairsDown
					continue
				case tree:
					if !tmxTemplate.treeTerrain {
						exportTiles[x+y*l.Width] = tmxTemplate.treeIDs[9]
						continue
					}
					exportTiles[x+y*l.Width] = get16Tile2(m, x, y, tile, tmxTemplate.treeIDs)
					continue
				}
				exportTiles[x+y*l.Width] = get16Tile(m, x, y, tile, tileIDs)
			}
		}
		return strings.Join(exportTiles, ",")
	}
	tmxTemplate.CSV = makeCSV(m.Layer("Tiles"), m.Layer("Tiles"))
	tmxTemplate.CSVFurniture = makeCSV(m.Layer("Furniture"), m.Layer("Tiles"))
}

func get16Tile(m Map, x, y int, tile rune, templateTiles []string) string {
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

func get16Tile2(m Map, x, y int, tile rune, templateTiles []string) string {
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
