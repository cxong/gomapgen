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
	// then isolated tile
	nothingID string
	floorIDs  []string
	floor2IDs []string
	wallID    string
	roomIDs   []string

	// Parameters for map generation
	floorTerrain  bool
	floor2Terrain bool
	wallTerrain   bool
	roomTerrain   bool

	// Parameters used for template export
	Width  int
	Height int
	CSV    string
}

// DawnLikeTemplate - using DawnLike tile set
var DawnLikeTemplate = TMXTemplate{
	"dawnlike",
	"1043",
	[]string{"1421", "1400", "1401", "1422", "1443", "1442", "1441", "1420", "1399", "1425", "1423", "1402", "1426", "1444", "1424", "1404"},
	[]string{"1176", "1155", "1156", "1177", "1198", "1197", "1196", "1175", "1154", "1180", "1178", "1157", "1181", "1199", "1179", "1159"},
	// TODO: wall tile sets
	"69",
	[]string{"1428", "1407", "1408", "1429", "1450", "1449", "1448", "1427", "1406", "1432", "1430", "1409", "1433", "1451", "1431", "1411"},
	false, true, true, true,
	0, 0, ""}

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
	exportTiles := make([]string, len(m.Tiles))
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tile, err := m.GetTile(x, y)
			if err != nil {
				panic(err)
			}
			var tileIDs []string
			switch tile {
			case nothing:
				exportTiles[x+y*m.Width] = tmxTemplate.nothingID
				continue
			case floor:
				if !tmxTemplate.floorTerrain {
					exportTiles[x+y*m.Width] = tmxTemplate.floorIDs[0]
					continue
				}
				tileIDs = tmxTemplate.floorIDs
			case floor2:
				if !tmxTemplate.floor2Terrain {
					exportTiles[x+y*m.Width] = tmxTemplate.floor2IDs[0]
					continue
				}
				tileIDs = tmxTemplate.floor2IDs
			case wall:
				/*if !tmxTemplate.wallTerrain {
					exportTiles[x+y*m.Width] = tmxTemplate.wallIDs[0]
					continue
				}*/
				exportTiles[x+y*m.Width] = tmxTemplate.wallID
				continue
			case room:
				if !tmxTemplate.roomTerrain {
					exportTiles[x+y*m.Width] = tmxTemplate.roomIDs[0]
					continue
				}
				tileIDs = tmxTemplate.roomIDs
			}
			exportTiles[x+y*m.Width] = get16Tile(m, x, y, tile, tileIDs)
		}
	}
	tmxTemplate.CSV = strings.Join(exportTiles, ",")
}

func get16Tile(m Map, x, y int, tile byte, templateTiles []string) string {
	up := isSameTile(m, x, y-1, tile)
	right := isSameTile(m, x+1, y, tile)
	down := isSameTile(m, x, y+1, tile)
	left := isSameTile(m, x-1, y, tile)
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

func isSameTile(m Map, x, y int, tile byte) bool {
	other, err := m.GetTile(x, y)
	return err != nil || other == tile
}
