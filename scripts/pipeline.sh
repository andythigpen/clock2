#!/bin/bash

# path to the directory containing the lottie.json files
LOTTIE_DIR="$1"
# path to the compiled lottileconverter program (https://github.com/sot-tech/LottieConverter)
LOTTIE_CONVERTER_DIR="$2"
LOTTIE_FILES=(clear-day clear-night code-orange code-red fog-day fog-night hail humidity overcast-day overcast-night partly-cloudy-day partly-cloudy-night pressure-high pressure-low rain sleet snow sunrise sunset thunderstorms-day thunderstorms-night thunderstorms-day-rain thunderstorms-night-rain wind)

mkdir -p pngs cropped sprites

for lottie in "${LOTTIE_FILES[@]}"; do
    file="${LOTTIE_DIR}/${lottie}.json"
    echo "$file"

    rm -f pngs/*.png cropped/*.png
    bash ./pngs.sh "$LOTTIE_CONVERTER_DIR" "$file"
    bash ./crop.sh
    bash ./sprite.sh "${lottie}"
done

rm ./sprites/*-1.png
for file in ./sprites/*-0.png; do
    if [[ -e "$file" ]]; then
        mv "$file" "${file/-0/}"
    fi
done

mv ./sprites/{overcast,cloudy}-day.png
mv ./sprites/{overcast,cloudy}-night.png
mv ./sprites/thunderstorms-{day-rain,rain-day}.png
mv ./sprites/thunderstorms-{night-rain,rain-night}.png
mv ./sprites/wind{,y}.png
