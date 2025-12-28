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
