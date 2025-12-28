#!/bin/sh

# FFmpeg script to composite overlay onto RTSP stream

INPUT_RTSP="${INPUT_RTSP:-rtsp://camera:554/stream}"
OUTPUT_RTSP="${OUTPUT_RTSP:-rtsp://localhost:8554/live}"
OVERLAY_FILE="${OVERLAY_FILE:-/tmp/overlay/current.png}"
OVERLAY_X="${OVERLAY_X:-10}"
OVERLAY_Y="${OVERLAY_Y:-10}"

echo "Starting FFmpeg stream compositor"
echo "Input: $INPUT_RTSP"
echo "Output: $OUTPUT_RTSP"
echo "Overlay position: ${OVERLAY_X},${OVERLAY_Y}"

# Wait for overlay file to be created
while [ ! -f "$OVERLAY_FILE" ]; do
    echo "Waiting for overlay file to be created..."
    sleep 1
done

# FFmpeg command to composite overlay
ffmpeg \
    -rtsp_transport tcp \
    -i "$INPUT_RTSP" \
    -loop 1 \
    -framerate 1 \
    -i "$OVERLAY_FILE" \
    -filter_complex "[0:v][1:v]overlay=${OVERLAY_X}:${OVERLAY_Y}:format=auto:repeatlast=1" \
    -c:v libx264 \
    -preset ultrafast \
    -tune zerolatency \
    -c:a copy \
    -f rtsp \
    "$OUTPUT_RTSP"
