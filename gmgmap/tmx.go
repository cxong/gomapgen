package gmgmap

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// TMXTemplate - configuration for TMX export
type TMXTemplate struct {
	path      string
	nothingID int
	floorID   int
	floor2ID  int
	wallID    int
	roomID    int
	Width     int
	Height    int
	CSV       string
}

// DawnLikeTemplate - using DawnLike tile set
var DawnLikeTemplate = TMXTemplate{"dawnlike", 1043, 1176, 1421, 69, 1428, 0, 0, ""}

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
	// Populate template
	tmxTemplate.Width = m.Width
	tmxTemplate.Height = m.Height
	exportTiles := make([]string, len(m.Tiles))
	for i := 0; i < len(m.Tiles); i++ {
		switch m.Tiles[i] {
		case nothing:
			exportTiles[i] = strconv.Itoa(tmxTemplate.nothingID)
		case floor:
			exportTiles[i] = strconv.Itoa(tmxTemplate.floorID)
		case floor2:
			exportTiles[i] = strconv.Itoa(tmxTemplate.floor2ID)
		case wall:
			exportTiles[i] = strconv.Itoa(tmxTemplate.wallID)
		case room:
			exportTiles[i] = strconv.Itoa(tmxTemplate.roomID)
		}
	}
	tmxTemplate.CSV = strings.Join(exportTiles, ",")
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
