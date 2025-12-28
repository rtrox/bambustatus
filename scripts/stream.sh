#!/bin/sh

# FFmpeg script to composite overlay onto RTMP stream

INPUT_RTSP="${INPUT_RTSP:-rtsp://camera:554/stream}"
OUTPUT_RTMP="${OUTPUT_RTMP:-rtmp://localhost/live}"
OVERLAY_FILE="${OVERLAY_FILE:-/tmp/overlay/current.png}"
OVERLAY_X="${OVERLAY_X:-10}"
OVERLAY_Y="${OVERLAY_Y:-10}"

echo "Starting FFmpeg stream compositor"
echo "Input: $INPUT_RTSP"
echo "Output: $OUTPUT_RTMP"
echo "Overlay position: ${OVERLAY_X},${OVERLAY_Y}"

# Wait for overlay file to be created
while [ ! -f "$OVERLAY_FILE" ]; do
    echo "Waiting for overlay file to be created..."
    sleep 1
done

# Give the overlay capture a moment to stabilize
sleep 2

echo "Starting FFmpeg stream..."

# FFmpeg command to composite overlay
# Note: Using exec to replace shell process and properly handle signals
# Generates silent audio track since RTSP source has no audio
exec ffmpeg \
    -loglevel info \
    -rtsp_transport tcp \
    -i "$INPUT_RTSP" \
    -f lavfi \
    -i anullsrc=channel_layout=stereo:sample_rate=44100 \
    -loop 1 \
    -framerate 1 \
    -i "$OVERLAY_FILE" \
    -filter_complex "[0:v][2:v]overlay=${OVERLAY_X}:${OVERLAY_Y}:format=auto:repeatlast=1" \
    -c:v libx264 \
    -preset ultrafast \
    -tune zerolatency \
    -g 60 \
    -keyint_min 60 \
    -sc_threshold 0 \
    -b:v 2500k \
    -maxrate 2500k \
    -bufsize 5000k \
    -pix_fmt yuv420p \
    -c:a aac \
    -b:a 128k \
    -ar 44100 \
    -shortest \
    -f flv \
    -flvflags no_duration_filesize \
    "$OUTPUT_RTMP"
