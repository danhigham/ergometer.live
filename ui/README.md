# Ergometer.Live UI

Vue.js frontend for the Ergometer.Live PM5 tracking system.

## Features

- Real-time WebSocket connection to local PM5 server
- Live workout statistics display
- WebSocket connection testing page

## Development Setup

### Prerequisites

- Node.js 18+ and npm
- PM5 WebSocket server running on `localhost:8080`

### Install Dependencies

```sh
npm install
```

### Run Development Server

```sh
npm run dev
```

The app will be available at `http://localhost:5173` by default.

### Testing WebSocket Connection

1. Make sure the PM5 server is running:
   ```bash
   cd ../
   go run main.go
   ```

2. Navigate to the test page at `http://localhost:5173/test`

3. Click "Connect" to establish a WebSocket connection

4. You should see real-time messages from the PM5 server

## Project Structure

```
ui/
├── src/
│   ├── components/       # Reusable Vue components
│   ├── composables/      # Vue composables
│   │   └── useWebSocket.ts  # WebSocket connection management
│   ├── views/           # Page components
│   │   ├── HomeView.vue
│   │   └── WebSocketTest.vue  # WebSocket test page
│   ├── router/          # Vue Router configuration
│   └── App.vue          # Root component
├── public/              # Static assets
└── package.json
```

## WebSocket API

The UI connects to the PM5 server's WebSocket endpoint at `ws://localhost:8080/ws`.

### Message Types Received

- `workout_stats` - Real-time workout statistics (during active workouts)
- `workout_state` - Workout state changes
- `workout_started` - Workout has started
- `workout_ended` - Workout has ended
- `status` - Device status information

### Message Types Sent

- `start_workout` - Start a new workout
- `stop_workout` - Stop the current workout
- `get_status` - Get device status

## Build Commands

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

The built files will be in the `dist/` directory.

### Type Checking

```sh
npm run type-check
```

### Linting

```sh
npm run lint
```

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)
