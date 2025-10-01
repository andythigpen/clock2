#!/bin/bash

input_dir="./cropped"
output_file="sprites/$1.png"

mkdir -p "$(dirname "$output_file")"

montage "$input_dir"/*.png -verbose -background transparent -tile 30x4 -geometry +0+0 "$output_file"
