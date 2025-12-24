package socketserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danhigham/ergometer.live/broadcast"
	"github.com/danhigham/ergometer.live/pm5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict to your domain in production
		return true
	},
}

// ClientMessage represents a message received from a WebSocket client
type ClientMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// serveWs handles websocket requests from clients
func serveWs(hub *broadcast.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := broadcast.NewClient(hub, conn)
	hub.Register(client)

	// Start the client's read and write pumps
	client.Run()
}

// handleClientMessage processes messages received from WebSocket clients
func handleClientMessage(manager *pm5.Manager, client *broadcast.Client, message []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Failed to unmarshal client message: %v", err)
		sendError(client, "Invalid message format")
		return
	}

	log.Printf("Received message type: %s", msg.Type)

	switch msg.Type {
	case "start_workout":
		handleStartWorkout(manager, client, msg.Data)

	case "stop_workout":
		handleStopWorkout(manager, client)

	case "get_status":
		handleGetStatus(manager, client)

	default:
		sendError(client, "Unknown message type: "+msg.Type)
	}
}

// handleStartWorkout processes a start_workout request
func handleStartWorkout(manager *pm5.Manager, client *broadcast.Client, data map[string]interface{}) {
	params := &pm5.WorkoutParams{}

	// Parse workout type
	if workoutType, ok := data["workout_type"].(string); ok {
		params.WorkoutType = workoutType
	} else {
		sendError(client, "workout_type is required")
		return
	}

	// Parse distance (for fixed_distance)
	if distance, ok := data["distance"].(float64); ok {
		params.Distance = uint32(distance)
	}

	// Parse time (for fixed_time)
	if timeVal, ok := data["time"].(float64); ok {
		params.Time = uint32(timeVal)
	}

	// Parse split distance
	if splitDistance, ok := data["split_distance"].(float64); ok {
		params.SplitDistance = uint32(splitDistance)
	}

	// Parse split time
	if splitTime, ok := data["split_time"].(float64); ok {
		params.SplitTime = uint32(splitTime)
	}

	// Send control request to manager
	resp, err := manager.SendControl("start_workout", params)
	if err != nil {
		sendError(client, "Failed to send start command: "+err.Error())
		return
	}

	if !resp.Success {
		sendError(client, "Failed to start workout: "+resp.Error.Error())
		return
	}

	sendSuccess(client, "start_workout", "Workout started successfully")
}

// handleStopWorkout processes a stop_workout request
func handleStopWorkout(manager *pm5.Manager, client *broadcast.Client) {
	resp, err := manager.SendControl("stop_workout", nil)
	if err != nil {
		sendError(client, "Failed to send stop command: "+err.Error())
		return
	}

	if !resp.Success {
		sendError(client, "Failed to stop workout: "+resp.Error.Error())
		return
	}

	sendSuccess(client, "stop_workout", "Workout stopped successfully")
}

// handleGetStatus processes a get_status request
func handleGetStatus(manager *pm5.Manager, client *broadcast.Client) {
	resp, err := manager.SendControl("get_status", nil)
	if err != nil {
		sendError(client, "Failed to get status: "+err.Error())
		return
	}

	if !resp.Success {
		sendError(client, "Failed to get status")
		return
	}

	sendStatus(client, resp.Data)
}

// sendError sends an error message to a client
func sendError(client *broadcast.Client, message string) {
	msg := map[string]interface{}{
		"type": "error",
		"data": map[string]string{
			"message": message,
			"code":    "ERROR",
		},
	}

	data, _ := json.Marshal(msg)
	client.Send(data)
}

// sendSuccess sends a success message to a client
func sendSuccess(client *broadcast.Client, action, message string) {
	msg := map[string]interface{}{
		"type": "success",
		"data": map[string]string{
			"action":  action,
			"message": message,
		},
	}

	data, _ := json.Marshal(msg)
	client.Send(data)
}

// sendStatus sends a status response to a client
func sendStatus(client *broadcast.Client, statusData interface{}) {
	msg := map[string]interface{}{
		"type": "status",
		"data": statusData,
	}

	data, _ := json.Marshal(msg)
	client.Send(data)
}
