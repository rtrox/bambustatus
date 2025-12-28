#!/bin/sh

# Script to continuously capture the overlay as a PNG with transparency

OVERLAY_URL="${OVERLAY_URL:-http://localhost:8080}"
OUTPUT_FILE="${OVERLAY_FILE:-/tmp/overlay/current.png}"
CAPTURE_INTERVAL="${CAPTURE_INTERVAL:-1}"

echo "Starting overlay capture from $OVERLAY_URL"
echo "Saving to $OUTPUT_FILE every ${CAPTURE_INTERVAL}s"

TEMP_FILE="${OUTPUT_FILE}-tmp.png"

while true; do
    chromium-browser \
        --headless=new \
        --disable-gpu \
        --no-sandbox \
        --disable-dev-shm-usage \
        --disable-software-rasterizer \
        --disable-extensions \
        --disable-features=VizDisplayCompositor \
        --run-all-compositor-stages-before-draw \
        --disable-background-networking \
        --disable-sync \
        --metrics-recording-only \
        --disable-default-apps \
        --mute-audio \
        --no-first-run \
        --disable-hang-monitor \
        --disable-prompt-on-repost \
        --disable-breakpad \
        --virtual-time-budget=10000 \
        --window-size=1744,400 \
        --default-background-color=00000000 \
        --screenshot="$TEMP_FILE" \
        "$OVERLAY_URL"

    # Atomically replace the file so FFmpeg never sees a partial write
    mv "$TEMP_FILE" "$OUTPUT_FILE"

    sleep "$CAPTURE_INTERVAL"
done
