#!/bin/bash

# Stop the Ergometer.Live development environment

SESSION_NAME="ergometer-dev"

if tmux has-session -t $SESSION_NAME 2>/dev/null; then
    echo "Stopping development environment..."
    tmux kill-session -t $SESSION_NAME
    echo "Session '$SESSION_NAME' terminated."
else
    echo "No active session found."
fi
