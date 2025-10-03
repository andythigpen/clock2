#!/bin/bash

# path to the directory containing the lottie.json files
LOTTIE_DIR="$1"
# path to the compiled lottileconverter program (https://github.com/sot-tech/LottieConverter)
LOTTIE_CONVERTER_DIR="$2"
LOTTIE_FILES_480=(clear-day clear-night code-orange code-red fog-day fog-night hail humidity overcast-day overcast-night partly-cloudy-day partly-cloudy-night rain sleet snow sunrise sunset thunderstorms-day thunderstorms-night thunderstorms-day-rain thunderstorms-night-rain wind)
LOTTIE_FILES_400=(pressure-high pressure-low)

mkdir -p pngs cropped sprites

for lottie in "${LOTTIE_FILES_480[@]}"; do
    file="${LOTTIE_DIR}/${lottie}.json"
    echo "$file"
    rm -f pngs/*.png cropped/*.png
    bash ./pngs.sh "$LOTTIE_CONVERTER_DIR" "$file" "480x480"
    bash ./crop.sh
    bash ./sprite.sh "${lottie}"
done

for lottie in "${LOTTIE_FILES_400[@]}"; do
    file="${LOTTIE_DIR}/${lottie}.json"
    echo "$file"
    rm -f pngs/*.png cropped/*.png
    bash ./pngs.sh "$LOTTIE_CONVERTER_DIR" "$file" "400x400"
    bash ./crop.sh
    bash ./sprite.sh "${lottie}"
done

# rename files to more closely match HA states
for img in ./sprites/overcast-*.png; do
    mv "$img" "${img/overcast/cloudy}"
done
for img in ./sprites/thunderstorms-day-rain-*.png; do
    mv "$img" "${img/day-rain/rain-day}"
done
for img in ./sprites/thunderstorms-night-rain-*.png; do
    mv "$img" "${img/night-rain/rain-night}"
done
for img in ./sprites/wind-*.png; do
    mv "$img" "${img/wind/windy}"
done
