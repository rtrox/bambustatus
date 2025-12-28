#!/bin/sh

# Test script to debug RTMP connectivity

RTMP_URL="${1:-rtmp://owncast-stream.rtrox.io/live/Bambulabs^1}"

echo "=== RTMP Connection Test ==="
echo "Target: $RTMP_URL"
echo ""

# Extract hostname from RTMP URL
HOSTNAME=$(echo "$RTMP_URL" | sed 's|rtmp://||' | cut -d'/' -f1)
echo "Hostname: $HOSTNAME"

# Test DNS resolution
echo ""
echo "=== DNS Resolution ==="
nslookup "$HOSTNAME" || echo "DNS lookup failed"

# Test network connectivity
echo ""
echo "=== Network Connectivity ==="
ping -c 3 "$HOSTNAME" || echo "Ping failed (may be blocked)"

# Test RTMP port (default 1935)
echo ""
echo "=== Port 1935 Test ==="
nc -zv "$HOSTNAME" 1935 2>&1 || echo "Port 1935 not reachable"

# Try a minimal FFmpeg test
echo ""
echo "=== FFmpeg RTMP Test (5 second test pattern) ==="
timeout 10 ffmpeg \
    -re \
    -f lavfi \
    -i testsrc=duration=5:size=320x240:rate=30 \
    -f lavfi \
    -i sine=frequency=1000:duration=5 \
    -c:v libx264 \
    -preset ultrafast \
    -b:v 500k \
    -c:a aac \
    -b:a 128k \
    -f flv \
    "$RTMP_URL" 2>&1 | tail -20

echo ""
echo "=== Test Complete ==="
