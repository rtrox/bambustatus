# BambuStatus

An OBS-ready HTML overlay generator for displaying Bambu Lab 3D printer status in real-time via MQTT.

## Features

- Real-time printer status via MQTT
- Direct connection to Bambu Lab printer's built-in MQTT broker
- Progress bar with percentage
- Layer information (current/total)
- Temperature monitoring:
  - Nozzle temperature (current/target)
  - Bed temperature (current/target)
  - Chamber/Ambient temperature
- Time tracking (elapsed/remaining)
- Auto-discovery of printer serial number
- Clean, customizable overlay design
- Transparent background for OBS
- **RTSP stream compositor** - Overlay status onto video streams

## Project Structure

```
bambustatus/
├── cmd/
│   └── bambustatus/
│       └── main.go          # Main application with MQTT client
├── pkg/
│   └── printer/
│       ├── status.go        # Printer status data structures
│       ├── bambu.go         # Bambu Lab message parsing
│       └── mqtt.go          # MQTT client implementation
├── web/
│   ├── templates/
│   │   └── overlay.html     # HTML template for OBS overlay
│   └── static/
│       ├── css/
│       │   └── style.css    # Overlay styling
│       └── js/
│           └── updater.js   # Auto-refresh script
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- A Bambu Lab 3D printer with LAN-only mode enabled
- Printer's IP address and access code

### Installation

1. Clone the repository:
```bash
git clone https://github.com/rtrox/bambustatus.git
cd bambustatus
```

2. Build the application:
```bash
go build -o bambustatus ./cmd/bambustatus
```

### Configuration

Before running, you need:

- **Printer IP address**: Find this in your printer's network settings
- **Access code**: Found in your printer's settings under "LAN Only Mode"

### Running

Start the server with your printer's credentials:
```bash
./bambustatus -host <printer-ip> -password <access-code>
```

Example:
```bash
./bambustatus -host 192.168.1.100 -password 12345678
```

The HTTP server will start on `http://localhost:8080` and automatically connect to your printer via MQTT.

### Command-Line Options

- `-host` (required): Printer IP address or hostname
- `-password` (required): Printer access code
- `-port`: MQTT port (default: 8883)
- `-username`: MQTT username (default: bblp)
- `-serial`: Printer serial number (auto-discovered if not specified)
- `-http-port`: HTTP server port (default: 8080)

## Usage

### OBS Setup

1. In OBS, add a new "Browser" source
2. Set the URL to: `http://localhost:8080/`
3. Set dimensions:
   - Width: 1744
   - Height: 250
4. Check "Shutdown source when not visible" to save resources
5. Adjust position and scale as needed

### API Endpoints

#### GET `/`
Returns the HTML overlay page for OBS

#### GET `/api/status`
Returns the current printer status as JSON (automatically populated from MQTT)

## Customization

### Styling

Edit [web/static/css/style.css](web/static/css/style.css) to customize:
- Colors and fonts
- Layout and spacing
- Background transparency
- Temperature color coding

### Refresh Rate

Edit [web/static/js/updater.js](web/static/js/updater.js) to change the refresh interval:
```javascript
const REFRESH_INTERVAL = 2000; // milliseconds
```

## RTSP Stream Compositor

BambuStatus can overlay printer status onto an RTSP video stream (like a camera feed of your printer).

### Docker Compose Setup

1. Create or edit [docker-compose.streamer.yml](docker-compose.streamer.yml):

```yaml
services:
  bambustatus-streamer:
    build:
      context: .
      dockerfile: Dockerfile.streamer
    ports:
      - "8080:8080"  # Web interface
      - "8554:8554"  # RTSP output
    environment:
      # Printer settings
      PRINTER_HOST: "192.168.1.100"
      PRINTER_PASSWORD: "12345678"

      # Input RTSP stream
      INPUT_RTSP: "rtsp://camera:554/stream"

      # Overlay position (optional)
      OVERLAY_X: "10"
      OVERLAY_Y: "10"
      CAPTURE_INTERVAL: "1"
    restart: unless-stopped
```

2. Start the compositor:

```bash
docker-compose -f docker-compose.streamer.yml up -d
```

3. View the composited stream at `rtsp://localhost:8554/live`

### Configuration Options

- `PRINTER_HOST` - Printer IP address (required)
- `PRINTER_PASSWORD` - Printer access code (required)
- `INPUT_RTSP` - Source RTSP stream URL (required)
- `OUTPUT_RTSP` - Output RTSP stream URL (default: `rtsp://0.0.0.0:8554/live`)
- `OVERLAY_X` - Horizontal position in pixels (default: `10`)
- `OVERLAY_Y` - Vertical position in pixels (default: `10`)
- `CAPTURE_INTERVAL` - How often to update overlay in seconds (default: `1`)

### How It Works

1. **BambuStatus Server**: Runs the web server with real-time printer data
2. **Overlay Capture**: Chromium headless captures the overlay as PNG every second
3. **FFmpeg Compositor**: Composites the overlay onto the input RTSP stream
4. **RTSP Output**: Serves the composited stream on port 8554

## Development

### Running in Development

```bash
go run ./cmd/bambustatus/main.go -host <printer-ip> -password <access-code>
```

### Building

```bash
go build -o bambustatus ./cmd/bambustatus
```

## How It Works

1. **MQTT Connection**: The application connects to your Bambu Lab printer's built-in MQTT broker using TLS (with certificate verification disabled for local connections)
2. **Auto-Discovery**: If you don't provide a serial number, it subscribes to `device/+/report` to automatically discover your printer
3. **Real-Time Updates**: As your printer publishes status updates via MQTT, the application parses the messages and updates the internal status
4. **HTTP Server**: The web interface polls the server every 2 seconds to refresh the overlay display
5. **OBS Integration**: OBS displays the HTML page with a transparent background

## Troubleshooting

### Connection Issues

- Ensure your printer is in LAN-only mode
- Verify the IP address is correct (check your router or printer display)
- Confirm the access code matches what's shown in printer settings
- Make sure your firewall allows connections on port 8883

### No Data Showing

- Wait a few seconds after connecting - the printer sends updates periodically
- Start a print job to see live data
- Check the console output for MQTT messages

## License

MIT
