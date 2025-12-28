#!/bin/sh

# Script to continuously capture the overlay as a PNG with transparency

OVERLAY_URL="${OVERLAY_URL:-http://localhost:8080}"
OUTPUT_FILE="${OVERLAY_FILE:-/tmp/overlay/current.png}"
CAPTURE_INTERVAL="${CAPTURE_INTERVAL:-1}"

echo "Starting overlay capture from $OVERLAY_URL"
echo "Saving to $OUTPUT_FILE every ${CAPTURE_INTERVAL}s"

while true; do
    chromium-browser \
        --headless \
        --disable-gpu \
        --no-sandbox \
        --disable-dev-shm-usage \
        --window-size=1744,250 \
        --screenshot="$OUTPUT_FILE" \
        --default-background-color=0 \
        "$OVERLAY_URL" \
        2>/dev/null

    sleep "$CAPTURE_INTERVAL"
done
