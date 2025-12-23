# Ergometer.Live - PM5 WebSocket Server

A real-time WebSocket server for Concept2 PM5 rowing machines that provides live workout data streaming to web clients.

## Features

- Real-time workout statistics streaming via WebSocket
- Support for multiple concurrent clients
- Adaptive polling (100ms during workouts, 1s when idle)
- Web-based test interface for controlling and monitoring workouts
- Support for multiple workout types:
  - Just Row (free rowing with optional splits)
  - Fixed Distance (e.g., 2000m with 500m splits)
  - Fixed Time (e.g., 20min with 60s splits)

## Architecture

- **Broadcast Hub**: Fan-out WebSocket broadcast with buffered channels
- **PM5 Manager**: Singleton managing USB device connection and control
- **Monitor Goroutine**: Adaptive polling for real-time stats
- **Web Interface**: Simple HTML/CSS/JS test client

## Requirements

- Go 1.21 or later
- USB connection to Concept2 PM5 rowing machine
- macOS, Linux, or Windows

## Building

```bash
go mod tidy
go build -o ergometer-live
```

## Running

```bash
./ergometer-live
```

The server will start on port 8080. Open your web browser to:

```
http://localhost:8080
```

## WebSocket API

### Client → Server Messages

**Start Workout:**
```json
{
  "type": "start_workout",
  "data": {
    "workout_type": "fixed_distance",
    "distance": 2000,
    "split_distance": 500
  }
}
```

**Stop Workout:**
```json
{
  "type": "stop_workout"
}
```

**Get Status:**
```json
{
  "type": "get_status"
}
```

### Server → Client Messages

**Workout Stats (real-time):**
```json
{
  "type": "workout_stats",
  "data": {
    "elapsed_time": 125.5,
    "distance": 512.5,
    "pace": 125.5,
    "power": 185,
    "stroke_rate": 24,
    "calories": 42,
    "heart_rate": 145,
    "workout_state": "Workout Row",
    "rowing_state": "Active"
  }
}
```

**Device Status:**
```json
{
  "type": "status",
  "data": {
    "connected": true,
    "serial": "PM5-123456",
    "model": 5,
    "battery": 85,
    "erg_type": "Rower Model D"
  }
}
```

**Error:**
```json
{
  "type": "error",
  "data": {
    "message": "Failed to start workout",
    "code": "ERROR"
  }
}
```

## Project Structure

```
ergometer.live/
├── main.go                  # Entry point
├── server/                  # HTTP & WebSocket server
│   ├── server.go
│   ├── websocket.go
│   └── handler.go
├── pm5/                     # PM5 device manager
│   ├── manager.go
│   └── monitor.go
├── broadcast/               # WebSocket broadcast hub
│   ├── hub.go
│   └── client.go
└── web/                     # Test web interface
    ├── index.html
    ├── app.js
    └── styles.css
```

## Development

The server uses the local `usb-interface` library via Go module replace directive:

```go
replace github.com/danhigham/pm5 => ../usb-interface
```

## License

See parent project for license information.
