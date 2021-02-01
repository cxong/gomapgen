#!/usr/bin/env bash

CMD="go run main.go --corridorWidth 2"
MAP=tmx_export/map.tmx
MACOS_TMX_RASTERIZER=/Applications/Tiled.app/Contents/MacOS/tmxrasterizer
WIN_TMX_RASTERIZER="/c/Program Files/Tiled/tmxrasterizer.exe"
if [ -f "${MACOS_TMX_RASTERIZER}" ]; then
    TMX_RASTERIZER=$MACOS_TMX_RASTERIZER
elif [ -f "${WIN_TMX_RASTERIZER}" ]; then
    TMX_RASTERIZER=$WIN_TMX_RASTERIZER
else
    echo "Tiled not found"
    read -p "Press enter to continue"
    exit
fi
OUT=/tmp
ITERATIONS=5
DELAY=200

for i in $(seq 1 $ITERATIONS)
do
    $CMD
    "${TMX_RASTERIZER}" $MAP $OUT/map$i.png
done

convert -delay $DELAY -dispose previous $OUT/map*.png $OUT/map.gif
echo "Rendered to ${OUT}/map.gif"
read -p "Press enter to continue"
