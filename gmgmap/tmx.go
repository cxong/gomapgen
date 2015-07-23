package gmgmap

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TMXTemplate - configuration for TMX export
type TMXTemplate struct {
	templatePath string
	nothingID    int
	floorID      int
	wallID       int
	roomID       int
	Width        int
	Height       int
	CSV          []int
}

// DawnLikeTemplate - using DawnLike tile set
var DawnLikeTemplate = TMXTemplate{"dawnlike", 1043, 1176, 69, 1428, 0, 0, []int{}}

// ToTMX - export map as TMX (Tiled XML map)
func (m Map) ToTMX(template *TMXTemplate) error {
	exportDir := "tmx_export"
	err := os.Mkdir(exportDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
		return err
	}
	// Copy data files
	baseDir := path.Join("gmgmap", template.templatePath)
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
	template.Width = m.Width
	template.Height = m.Height
	for i := 0; i < len(m.Tiles); i++ {
		switch m.Tiles[i] {
		case nothing:
			template.CSV = append(template.CSV, template.nothingID)
		case floor:
			template.CSV = append(template.CSV, template.floorID)
		case wall:
			template.CSV = append(template.CSV, template.wallID)
		case room:
			template.CSV = append(template.CSV, template.roomID)
		}
	}
	return nil
}
