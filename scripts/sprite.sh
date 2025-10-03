#!/bin/bash

input_dir="./cropped"
output_file="sprites/$1.png"

mkdir -p "$(dirname "$output_file")"

rows=30
cols=4

for img in "$input_dir"/*.png; do
    dimensions=$(convert "$img" -format "%wx%h" info:)
    width=${dimensions%x*}
    height=${dimensions#*x}
    cols=$((2048 / width))
    rows=$((2048 / height))
    echo "width=$width height=$height cols=$cols rows=$rows"
    break
done

montage "$input_dir"/*.png -verbose -background transparent -tile "${cols}x${rows}" -geometry +0+0 "$output_file"
