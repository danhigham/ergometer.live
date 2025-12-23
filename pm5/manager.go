package pm5

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/danhigham/ergometer.live/broadcast"
	pm5lib "github.com/danhigham/pm5"
	"github.com/danhigham/pm5/csafe"
	"github.com/danhigham/pm5/device"
)

var (
	instance *Manager
	once     sync.Once
)

// WorkoutParams contains parameters for starting a workout
type WorkoutParams struct {
	WorkoutType   string `json:"workout_type"`   // just_row, fixed_distance, fixed_time
	Distance      uint32 `json:"distance"`       // meters (for fixed_distance)
	Time          uint32 `json:"time"`           // seconds (for fixed_time)
	SplitDistance uint32 `json:"split_distance"` // meters (optional)
	SplitTime     uint32 `json:"split_time"`     // seconds (optional)
}

// ControlRequest represents a control command sent to the manager
type ControlRequest struct {
	Type     string         // "start_workout", "stop_workout", "get_status"
	Data     *WorkoutParams // Workout parameters (for start_workout)
	Response chan *ControlResponse
}

// ControlResponse contains the result of a control request
type ControlResponse struct {
	Success bool
	Data    interface{}
	Error   error
}

// Manager manages the PM5 device connection and monitoring
type Manager struct {
	pm5Device *pm5lib.PM5
	hub       *broadcast.Hub

	// Control channel for workout commands
	controlChan chan *ControlRequest

	// Monitor control
	stopMonitor chan struct{}
	monitorWg   sync.WaitGroup

	// State
	mu               sync.RWMutex
	isMonitoring     bool
	lastWorkoutState csafe.WorkoutState
	connected        bool
	deviceInfo       *DeviceInfo
}

// DeviceInfo contains information about the connected PM5 device
type DeviceInfo struct {
	Connected bool   `json:"connected"`
	Serial    string `json:"serial"`
	Model     int    `json:"model"`
	Battery   byte   `json:"battery"`
	ErgType   string `json:"erg_type"`
	OpState   string `json:"operational_state"`
}

// GetManager returns the singleton Manager instance
func GetManager(hub *broadcast.Hub) *Manager {
	once.Do(func() {
		instance = &Manager{
			hub:         hub,
			controlChan: make(chan *ControlRequest),
			stopMonitor: make(chan struct{}),
			connected:   false,
		}

		// Try to connect to PM5
		if err := instance.connect(); err != nil {
			log.Printf("Failed to connect to PM5: %v", err)
		}

		// Start control handler
		go instance.handleControl()

		// Start monitor
		go instance.startMonitor()
	})

	return instance
}

// connect attempts to connect to a PM5 device
func (m *Manager) connect() error {
	devices, err := device.EnumerateDevices()
	if err != nil {
		return fmt.Errorf("failed to enumerate devices: %w", err)
	}

	if len(devices) == 0 {
		return fmt.Errorf("no PM5 devices found")
	}

	usbDev := device.NewUSBDevice(devices[0])
	pm := pm5lib.New(usbDev)
	pm.SetDebug(true)

	if err := pm.Connect(); err != nil {
		return fmt.Errorf("failed to connect to PM5: %w", err)
	}

	m.pm5Device = pm
	m.connected = true

	// Get device info
	if err := m.updateDeviceInfo(); err != nil {
		log.Printf("Failed to get device info: %v", err)
	}

	log.Printf("Connected to PM5: %s", m.deviceInfo.ErgType)
	return nil
}

// updateDeviceInfo retrieves and caches device information
func (m *Manager) updateDeviceInfo() error {
	if m.pm5Device == nil {
		return fmt.Errorf("PM5 not connected")
	}

	info := &DeviceInfo{
		Connected: true,
	}

	// Get version/model
	if version, err := m.pm5Device.GetVersion(); err == nil {
		info.Model = int(version.Model)
	}

	// Get serial number
	if serial, err := m.pm5Device.GetSerial(); err == nil {
		info.Serial = serial
	}

	// Get battery level
	if battery, err := m.pm5Device.GetBatteryLevel(); err == nil {
		info.Battery = battery
	}

	// Get erg type
	if ergType, err := m.pm5Device.GetErgMachineType(); err == nil {
		info.ErgType = ergType.String()
	}

	// Get operational state
	if opState, err := m.pm5Device.GetOperationalState(); err == nil {
		info.OpState = opState.String()
	}

	// Get workout state
	if workoutState, err := m.pm5Device.GetWorkoutState(); err == nil {
		m.lastWorkoutState = workoutState
	}

	m.deviceInfo = info
	return nil
}

// handleControl processes control requests from the control channel
func (m *Manager) handleControl() {
	for req := range m.controlChan {
		switch req.Type {
		case "start_workout":
			err := m.startWorkout(req.Data)
			req.Response <- &ControlResponse{
				Success: err == nil,
				Error:   err,
			}

		case "stop_workout":
			err := m.stopWorkout()
			req.Response <- &ControlResponse{
				Success: err == nil,
				Error:   err,
			}

		case "get_status":
			m.mu.RLock()
			info := m.deviceInfo
			m.mu.RUnlock()

			req.Response <- &ControlResponse{
				Success: true,
				Data:    info,
			}

		default:
			req.Response <- &ControlResponse{
				Success: false,
				Error:   fmt.Errorf("unknown request type: %s", req.Type),
			}
		}
	}
}

// startWorkout starts a workout with the given parameters
func (m *Manager) startWorkout(params *WorkoutParams) error {
	if m.pm5Device == nil {
		return fmt.Errorf("PM5 not connected")
	}

	var err error

	switch params.WorkoutType {
	case "just_row":
		withSplits := params.SplitDistance > 0 || params.SplitTime > 0
		err = m.pm5Device.StartJustRowWorkout(withSplits)

	case "fixed_distance":
		if params.Distance == 0 {
			return fmt.Errorf("distance is required for fixed_distance workout")
		}
		err = m.pm5Device.StartFixedDistanceWorkout(params.Distance, params.SplitDistance)

	case "fixed_time":
		if params.Time == 0 {
			return fmt.Errorf("time is required for fixed_time workout")
		}
		// Convert seconds to hundredths of seconds
		duration := params.Time * 100
		splitDuration := params.SplitTime * 100
		err = m.pm5Device.StartFixedTimeWorkout(duration, splitDuration)

	default:
		return fmt.Errorf("unknown workout type: %s", params.WorkoutType)
	}

	if err != nil {
		return fmt.Errorf("failed to start workout: %w", err)
	}

	log.Printf("Started %s workout", params.WorkoutType)
	return nil
}

// stopWorkout terminates the current workout
func (m *Manager) stopWorkout() error {
	if m.pm5Device == nil {
		return fmt.Errorf("PM5 not connected")
	}

	if err := m.pm5Device.TerminateWorkout(); err != nil {
		return fmt.Errorf("failed to stop workout: %w", err)
	}

	log.Println("Workout stopped")
	return nil
}

// SendControl sends a control request and waits for the response
func (m *Manager) SendControl(reqType string, data *WorkoutParams) (*ControlResponse, error) {
	req := &ControlRequest{
		Type:     reqType,
		Data:     data,
		Response: make(chan *ControlResponse),
	}

	m.controlChan <- req

	// Wait for response with timeout
	select {
	case resp := <-req.Response:
		return resp, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("control request timeout")
	}
}

// BroadcastJSON marshals data to JSON and broadcasts it
func (m *Manager) BroadcastJSON(messageType string, data interface{}) {
	msg := map[string]interface{}{
		"type":      messageType,
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	m.hub.Broadcast(jsonData)
}

// Shutdown gracefully shuts down the manager
func (m *Manager) Shutdown() {
	log.Println("Shutting down PM5 manager...")

	// Stop monitor
	close(m.stopMonitor)
	m.monitorWg.Wait()

	// Close control channel
	close(m.controlChan)

	// Disconnect PM5
	if m.pm5Device != nil {
		m.pm5Device.Disconnect()
	}

	log.Println("PM5 manager shutdown complete")
}

// IsConnected returns whether the PM5 device is connected
func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}
