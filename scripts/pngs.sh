#!/bin/bash

LOTTIE_CONVERTER_DIR="$1"
LOTTIE_FILE="$2"

"$LOTTIE_CONVERTER_DIR/lottieconverter" "$LOTTIE_FILE" pngs/img pngs 480x480 20
