package printer

import (
	"encoding/json"
	"time"
)

// BambuMessage represents the top-level MQTT message structure
type BambuMessage struct {
	Print PrintData `json:"print"`
}

// PrintData contains all printer status information
type PrintData struct {
	BedTargetTemper    float64 `json:"bed_target_temper"`
	BedTemper          float64 `json:"bed_temper"`
	ChamberTemper      float64 `json:"chamber_temper"`
	NozzleTargetTemper float64 `json:"nozzle_target_temper"`
	NozzleTemper       float64 `json:"nozzle_temper"`

	GcodeState       string `json:"gcode_state"`
	SubtaskName      string `json:"subtask_name"`
	GcodeFile        string `json:"gcode_file"`

	McPercent        int `json:"mc_percent"`
	McRemainingTime  int `json:"mc_remaining_time"`
	LayerNum         int `json:"layer_num"`
	TotalLayerNum    int `json:"total_layer_num"`

	GcodeStartTime   string `json:"gcode_start_time"`
}

// ToStatus converts a BambuMessage to our internal Status structure
func (bm *BambuMessage) ToStatus() *Status {
	status := NewStatus()

	// Print name - prefer subtask_name, fallback to gcode_file
	if bm.Print.SubtaskName != "" {
		status.PrintName = bm.Print.SubtaskName
	} else if bm.Print.GcodeFile != "" {
		status.PrintName = bm.Print.GcodeFile
	} else {
		status.PrintName = "Idle"
	}

	// Progress
	status.Progress = float64(bm.Print.McPercent)

	// Layers
	status.CurrentLayer = bm.Print.LayerNum
	status.TotalLayers = bm.Print.TotalLayerNum

	// Temperatures
	status.NozzleTemp = bm.Print.NozzleTemper
	status.NozzleTempTarget = bm.Print.NozzleTargetTemper
	status.BedTemp = bm.Print.BedTemper
	status.BedTempTarget = bm.Print.BedTargetTemper
	status.AmbientTemp = bm.Print.ChamberTemper

	// Time
	status.TimeRemaining = bm.Print.McRemainingTime * 60 // Convert minutes to seconds

	// Calculate elapsed time from start time
	if bm.Print.GcodeStartTime != "" {
		// Parse Unix timestamp
		var startTime int64
		if err := json.Unmarshal([]byte(bm.Print.GcodeStartTime), &startTime); err == nil {
			elapsed := time.Since(time.Unix(startTime, 0))
			status.TimeElapsed = int(elapsed.Seconds())
		}
	}

	status.LastUpdated = time.Now()

	return status
}
