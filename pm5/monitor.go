package pm5

import (
	"log"
	"time"

	pm5lib "github.com/danhigham/pm5"
	"github.com/danhigham/pm5/csafe"
)

const (
	// PollIntervalActive is the polling interval during active workouts
	PollIntervalActive = 100 * time.Millisecond

	// PollIntervalIdle is the polling interval when idle
	PollIntervalIdle = 1 * time.Second

	// PollIntervalCheck is the default polling interval for state checks
	PollIntervalCheck = 500 * time.Millisecond
)

// WorkoutStats contains real-time workout statistics
type WorkoutStats struct {
	ElapsedTime   float64 `json:"elapsed_time"`    // seconds
	Distance      float64 `json:"distance"`        // meters
	Pace          float64 `json:"pace"`            // seconds per 500m
	AvgPace       float64 `json:"avg_pace"`        // seconds per 500m
	Power         uint32  `json:"power"`           // watts
	AvgPower      uint32  `json:"avg_power"`       // watts
	StrokeRate    byte    `json:"stroke_rate"`     // strokes per minute
	AvgStrokeRate byte    `json:"avg_stroke_rate"` // strokes per minute
	Calories      uint32  `json:"calories"`        // total calories
	HeartRate     byte    `json:"heart_rate"`      // bpm (0 = no data)
	AvgHeartRate  byte    `json:"avg_heart_rate"`  // bpm
	DragFactor    byte    `json:"drag_factor"`     // drag factor

	// State information
	WorkoutType  string `json:"workout_type"`
	WorkoutState string `json:"workout_state"`
	RowingState  string `json:"rowing_state"`
	StrokeState  string `json:"stroke_state"`
	OpState      string `json:"operational_state"`
}

// WorkoutStateInfo contains state change information
type WorkoutStateInfo struct {
	OpState      string `json:"operational_state"`
	WorkoutState string `json:"workout_state"`
	IsActive     bool   `json:"is_active"`
}

// startMonitor starts the monitoring goroutine
func (m *Manager) startMonitor() {
	m.monitorWg.Add(1)
	defer m.monitorWg.Done()

	m.mu.Lock()
	m.isMonitoring = true
	m.mu.Unlock()

	log.Println("Monitor started")
	m.monitorLoop()
}

// monitorLoop is the main monitoring loop with adaptive polling
func (m *Manager) monitorLoop() {
	ticker := time.NewTicker(PollIntervalCheck)
	defer ticker.Stop()

	currentInterval := PollIntervalCheck
	var lastWorkoutState csafe.WorkoutState

	for {
		select {
		case <-ticker.C:
			if m.pm5Device == nil {
				continue
			}

			// Get workout state
			workoutState, err := m.pm5Device.GetWorkoutState()

			if err != nil {
				log.Printf("Failed to get workout state: %v", err)
				continue
			}

			// Detect state transitions
			if lastWorkoutState != workoutState {
				m.broadcastStateChange(lastWorkoutState, workoutState)
				lastWorkoutState = workoutState
			}

			// Check if workout is active (rowing or in intervals)
			isActive := isWorkoutActive(workoutState)

			// Adjust polling rate based on state
			if isActive {
				// Active workout - poll fast
				if currentInterval != PollIntervalActive {
					currentInterval = PollIntervalActive
					ticker.Reset(currentInterval)
					log.Printf("Switched to active polling (%v) - state: %s", currentInterval, workoutState)
				}

				// Get and broadcast full workout snapshot
				m.broadcastWorkoutStats()

			} else {
				// Idle state - poll slowly
				if currentInterval != PollIntervalIdle {
					currentInterval = PollIntervalIdle
					ticker.Reset(currentInterval)
					log.Printf("Switched to idle polling (%v) - state: %s", currentInterval, workoutState)
				}

				// Just send state update
				m.broadcastStateOnly(workoutState)
			}

		case <-m.stopMonitor:
			log.Println("Monitor stopped")
			return
		}
	}
}

// isWorkoutActive returns true if the workout state indicates active rowing
func isWorkoutActive(state csafe.WorkoutState) bool {
	switch state {
	case csafe.WorkoutStateWorkoutRow,
		csafe.WorkoutStateIntervalWorkTime,
		csafe.WorkoutStateIntervalWorkDistance,
		csafe.WorkoutStateIntervalRest,
		csafe.WorkoutStateIntervalRestEndToWorkTime,
		csafe.WorkoutStateIntervalRestEndToWorkDistance,
		csafe.WorkoutStateIntervalWorkTimeToRest,
		csafe.WorkoutStateIntervalWorkDistanceToRest:
		return true
	default:
		return false
	}
}

// broadcastWorkoutStats gets the full workout snapshot and broadcasts it
func (m *Manager) broadcastWorkoutStats() {
	snapshot, err := m.pm5Device.GetWorkoutSnapshot()
	if err != nil {
		log.Printf("Failed to get workout snapshot: %v", err)
		return
	}

	stats := m.convertSnapshot(snapshot)
	m.BroadcastJSON("workout_stats", stats)
}

// broadcastStateOnly broadcasts just the workout state
func (m *Manager) broadcastStateOnly(workoutState csafe.WorkoutState) {
	isActive := isWorkoutActive(workoutState)

	stateInfo := &WorkoutStateInfo{
		OpState:      "", // Will be filled in convertSnapshot
		WorkoutState: workoutState.String(),
		IsActive:     isActive,
	}

	m.BroadcastJSON("workout_state", stateInfo)
}

// broadcastStateChange broadcasts a state transition event
func (m *Manager) broadcastStateChange(oldState, newState csafe.WorkoutState) {
	log.Printf("Workout state transition: %s -> %s", oldState, newState)

	wasActive := isWorkoutActive(oldState)
	isActive := isWorkoutActive(newState)

	stateInfo := &WorkoutStateInfo{
		OpState:      "", // Will be filled in convertSnapshot
		WorkoutState: newState.String(),
		IsActive:     isActive,
	}

	m.BroadcastJSON("workout_state", stateInfo)

	// Send specific events for workout start/end
	if !wasActive && isActive {
		m.BroadcastJSON("workout_started", map[string]string{"message": "Workout started", "state": newState.String()})
	}

	if wasActive && !isActive {
		m.BroadcastJSON("workout_ended", map[string]string{"message": "Workout ended", "state": newState.String()})
	}
}

// convertSnapshot converts a PM5 WorkoutSnapshot to our WorkoutStats format
func (m *Manager) convertSnapshot(snapshot *pm5lib.WorkoutSnapshot) *WorkoutStats {
	// Get operational state separately
	opState := ""
	if state, err := m.pm5Device.GetOperationalState(); err == nil {
		opState = state.String()
	}

	return &WorkoutStats{
		ElapsedTime:   snapshot.ElapsedTime.Seconds(),
		Distance:      snapshot.Distance,
		Pace:          snapshot.Pace.Seconds(),
		AvgPace:       snapshot.AvgPace.Seconds(),
		Power:         uint32(snapshot.Power),
		AvgPower:      uint32(snapshot.AvgPower),
		StrokeRate:    snapshot.StrokeRate,
		AvgStrokeRate: snapshot.AvgStrokeRate,
		Calories:      snapshot.Calories,
		HeartRate:     snapshot.HeartRate,
		AvgHeartRate:  snapshot.AvgHeartRate,
		DragFactor:    snapshot.DragFactor,
		WorkoutType:   snapshot.WorkoutType,
		WorkoutState:  snapshot.WorkoutState,
		RowingState:   snapshot.RowingState,
		StrokeState:   snapshot.StrokeState,
		OpState:       opState,
	}
}
