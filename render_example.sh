#!/bin/bash

CMD="go run main.go --algo=cell --width=24 --height=20 --template=kenney"
MAP=tmx_export/map.tmx
TMX_RASTERIZER=/Applications/Tiled.app/Contents/MacOS/tmxrasterizer
OUT=/tmp
ITERATIONS=5

for i in $(seq 1 $ITERATIONS)
do
    $CMD
    $TMX_RASTERIZER $MAP $OUT/map$i.png
done

convert -delay 100 -dispose previous $OUT/map*.png $OUT/map.gif
echo "Rendered to ${OUT}/map.gif"
