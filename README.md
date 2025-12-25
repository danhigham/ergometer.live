# Ergometer.Live

A full-stack real-time workout tracking application for Concept2 PM5 rowing machines with dual-mode support (online/local) and customizable widget dashboards.

## Features

### Current (Phase 1)
- **Dual-Mode Authentication**: Google OAuth (online) or local-only mode
- **Real-time WebSocket**: Live workout data streaming from PM5 device
- **REST API**: Firebase-authenticated backend for data persistence
- **Vue 3 Frontend**: Modern reactive UI with TypeScript
- **Multiple Workout Types**:
  - Just Row (free rowing with optional splits)
  - Fixed Distance (e.g., 2000m with 500m splits)
  - Fixed Time (e.g., 20min with 60s splits)

### Coming Soon
- Customizable widget dashboards with drag-and-drop
- Cloud workout history and analytics (online mode)
- Browser-based storage (local mode)
- Multiple named dashboard views

## Architecture

```
┌─────────────┐     WebSocket      ┌──────────────┐
│   Browser   │ ←─────────────────→ │  WebSocket   │
│  (Vue 3)    │                     │  Server      │
│             │     REST API        │  :8080       │
│             │ ←─────────────────→ │  (Go)        │
└─────────────┘                     └──────────────┘
                                           ↓
                                     PM5 Device
                                      (USB/HID)

┌─────────────┐     REST API       ┌──────────────┐
│   Browser   │ ←─────────────────→ │  REST API    │
│             │                     │  Server      │
│             │                     │  :3000       │
│             │                     │  (Go)        │
└─────────────┘                     └──────┬───────┘
                                           │
                    ┌──────────────────────┼───────────────┐
                    ↓                      ↓               ↓
              Firebase Auth          InfluxDB       (Future: PostgreSQL)
              (Google OAuth)         (Workouts)     (User metadata)
```

### Components
- **WebSocket Server** (port 8080): Real-time PM5 data streaming
- **REST API Server** (port 3000): Authentication, data persistence
- **Vue 3 Frontend** (port 5173): User interface with auth and widgets
- **PM5 Manager**: USB device control with adaptive polling
- **Broadcast Hub**: Fan-out WebSocket broadcast system

## Requirements

- Go 1.23 or later
- Node.js 20+ and npm
- USB connection to Concept2 PM5 rowing machine
- tmux (recommended for development)
- macOS, Linux, or Windows

## Quick Start (Development)

### 1. Install tmux (if not installed)
```bash
brew install tmux  # macOS
# or
apt-get install tmux  # Linux
```

### 2. Set up environment files

**Frontend** (`ui/.env`):
```bash
cd ui
cp .env.example .env
# Edit .env with your Firebase credentials
```

**Backend** (`api/.env`):
```bash
cd api
cp .env.example .env
# Edit .env with your Firebase and InfluxDB credentials
```

### 3. Install frontend dependencies
```bash
cd ui
npm install
```

### 4. Start all services with tmux
```bash
./start-dev.sh
```

This launches all three services in a single tmux session:
- **WebSocket Server** (left pane): `localhost:8080`
- **REST API Server** (top-right pane): `localhost:3000`
- **Vue Frontend** (bottom-right pane): `localhost:5173`

### tmux Quick Reference
- **Switch panes**: `Ctrl+b` then arrow keys
- **Detach session**: `Ctrl+b` then `d`
- **Reattach**: `./start-dev.sh` (or `tmux attach -t ergometer-dev`)
- **Stop all services**: `./stop-dev.sh`

### Manual Start (without tmux)

**Terminal 1 - WebSocket Server:**
```bash
go run main.go
```

**Terminal 2 - REST API Server:**
```bash
cd api
go run main.go
```

**Terminal 3 - Vue Frontend:**
```bash
cd ui
npm run dev
```

Then open your browser to `http://localhost:5173`

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
├── main.go                  # WebSocket server entry point
├── start-dev.sh             # tmux launcher for all services
├── stop-dev.sh              # Stop all services
├── socketserver/            # HTTP & WebSocket server
│   ├── server.go
│   ├── websocket.go
│   └── handler.go
├── pm5/                     # PM5 device manager
│   ├── manager.go
│   └── monitor.go
├── broadcast/               # WebSocket broadcast hub
│   ├── hub.go
│   └── client.go
├── web/                     # Static test pages
│   ├── index.html
│   └── test.html
├── api/                     # REST API server (NEW)
│   ├── main.go
│   ├── config/              # Configuration
│   ├── middleware/          # Auth, CORS, logging
│   ├── services/            # Firebase, InfluxDB
│   ├── handlers/            # HTTP handlers
│   └── models/              # Data models
└── ui/                      # Vue 3 frontend (NEW)
    ├── src/
    │   ├── views/           # Page components
    │   ├── components/      # Reusable components
    │   ├── stores/          # Pinia state stores
    │   ├── services/        # Firebase, API client
    │   └── router/          # Vue Router config
    └── package.json
```

## Development Notes

### Firebase Setup

1. Create a Firebase project at https://console.firebase.google.com
2. Enable Google Authentication
3. Get your web app config (apiKey, authDomain, projectId)
4. Download service account JSON for backend
5. Update environment files with credentials

### InfluxDB Setup (Optional - Phase 3+)

1. Create account at https://cloud2.influxdata.com
2. Create bucket: `ergometer-workouts`
3. Generate API token with read/write access
4. Update `api/.env` with credentials

### Local USB Interface

The WebSocket server uses the local `usb-interface` library:

```go
replace github.com/danhigham/pm5 => ../usb-interface
```

## Contributing

See the implementation plan at `.claude/plans/` for upcoming features and architecture decisions.

## License

See parent project for license information.
