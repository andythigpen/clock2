#!/bin/bash

LOTTIE_CONVERTER_DIR="$1"
LOTTIE_FILE="$2"
LOTTIE_SIZE="$3"

"$LOTTIE_CONVERTER_DIR/lottieconverter" "$LOTTIE_FILE" pngs/img pngs "$LOTTIE_SIZE" 20

# the last frame always matches the first, so we can remove it
lastframe=$(ls -1 ./pngs/img* | sort | tail -1)
rm "$lastframe"
