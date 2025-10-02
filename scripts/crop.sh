#!/bin/bash

input_dir="./pngs"
output_dir="./cropped"

mkdir -p "$output_dir"

min_x=99999
min_y=99999
max_x=0
max_y=0

for img in "$input_dir"/*.png; do
    dimensions=$(convert "$img" -trim -format "%wx%h" info:)
    width=${dimensions%x*}
    height=${dimensions#*x}
    offset=$(convert "$img" -trim -format "%[fx:page.x],%[fx:page.y]" info:)
    if [[ "$offset" == "-1,-1" ]]; then
        continue
    fi
    x_offset=${offset%,*}
    y_offset=${offset#*,}

    # Update the overall bounding box
    min_x=$(($min_x < $x_offset ? $min_x : $x_offset))
    min_y=$(($min_y < $y_offset ? $min_y : $y_offset))
    max_x=$(($max_x > $x_offset + width ? $max_x : $x_offset + width))
    max_y=$(($max_y > $y_offset + height ? $max_y : $y_offset + height))
done

final_width=$(($max_x - $min_x))
final_height=$(($max_y - $min_y))

echo "${final_width}x${final_height}"

for img in "$input_dir"/*.png; do
    filename=$(basename "$img")
    convert "$img" -crop "${final_width}x${final_height}+$min_x+$min_y" +repage "$output_dir/$filename"
done
