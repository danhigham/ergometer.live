#!/bin/bash

# Ergometer.Live Development Environment Launcher
# This script starts all services in a tmux session with three panes

SESSION_NAME="ergometer-dev"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo "Error: tmux is not installed"
    echo "Install with: brew install tmux"
    exit 1
fi

# Check if session already exists
if tmux has-session -t $SESSION_NAME 2>/dev/null; then
    echo "Session '$SESSION_NAME' already exists."
    echo "Attaching to existing session..."
    tmux attach-session -t $SESSION_NAME
    exit 0
fi

# Create new tmux session
echo "Starting Ergometer.Live development environment..."
echo "Project root: $PROJECT_ROOT"
echo ""

# Create session with first pane (WebSocket Server)
tmux new-session -d -s $SESSION_NAME -n "ergometer" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:0.0 "echo '=== WebSocket Server (PM5) ===' && echo 'Port: 8080'" C-m
tmux send-keys -t $SESSION_NAME:0.0 "go run main.go" C-m

# Split horizontally for REST API
tmux split-window -h -t $SESSION_NAME:0 -c "$PROJECT_ROOT/api"
tmux send-keys -t $SESSION_NAME:0.1 "echo '=== REST API Server ===' && echo 'Port: 3000'" C-m
tmux send-keys -t $SESSION_NAME:0.1 "sleep 2 && go run main.go" C-m

# Split the right pane vertically for Vue frontend
tmux split-window -v -t $SESSION_NAME:0.1 -c "$PROJECT_ROOT/ui"
tmux send-keys -t $SESSION_NAME:0.2 "echo '=== Vue Frontend ===' && echo 'Port: 5173'" C-m
tmux send-keys -t $SESSION_NAME:0.2 "sleep 3 && npm run dev" C-m

# Set pane layout (main-vertical gives nice side-by-side view)
tmux select-layout -t $SESSION_NAME:0 even-horizontal

# Select the first pane
tmux select-pane -t $SESSION_NAME:0.0

# Attach to session
echo ""
echo "Development environment started!"
echo ""
echo "Layout:"
echo "  ┌─────────────┬──────────────┐"
echo "  │  WebSocket  │  REST API    │"
echo "  │  :8080      │  :3000       │"
echo "  │             ├──────────────┤"
echo "  │             │  Vue UI      │"
echo "  │             │  :5173       │"
echo "  └─────────────┴──────────────┘"
echo ""
echo "Commands:"
echo "  Switch panes: Ctrl+b then arrow keys"
echo "  Detach: Ctrl+b then d"
echo "  Reattach: ./start-dev.sh (or tmux attach -t $SESSION_NAME)"
echo "  Kill session: tmux kill-session -t $SESSION_NAME"
echo ""
echo "Attaching in 2 seconds..."
sleep 2

tmux attach-session -t $SESSION_NAME
