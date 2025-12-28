package printer

import "time"

type Status struct {
	PrintName         string    `json:"print_name"`
	Progress          float64   `json:"progress"`
	CurrentLayer      int       `json:"current_layer"`
	TotalLayers       int       `json:"total_layers"`
	NozzleTemp        float64   `json:"nozzle_temp"`
	NozzleTempTarget  float64   `json:"nozzle_temp_target"`
	BedTemp           float64   `json:"bed_temp"`
	BedTempTarget     float64   `json:"bed_temp_target"`
	AmbientTemp       float64   `json:"ambient_temp"`
	TimeRemaining     int       `json:"time_remaining"` // in seconds
	TimeElapsed       int       `json:"time_elapsed"`   // in seconds
	LastUpdated       time.Time `json:"last_updated"`
}

func NewStatus() *Status {
	return &Status{
		PrintName:         "Idle",
		Progress:          0.0,
		CurrentLayer:      0,
		TotalLayers:       0,
		NozzleTemp:        0.0,
		NozzleTempTarget:  0.0,
		BedTemp:           0.0,
		BedTempTarget:     0.0,
		AmbientTemp:       0.0,
		TimeRemaining:     0,
		TimeElapsed:       0,
		LastUpdated:       time.Now(),
	}
}

func (s *Status) FormatTimeRemaining() string {
	if s.TimeRemaining <= 0 {
		return "--:--"
	}
	hours := s.TimeRemaining / 3600
	minutes := (s.TimeRemaining % 3600) / 60
	if hours > 0 {
		return formatTime(hours, minutes)
	}
	return formatTime(0, minutes)
}

func (s *Status) FormatTimeElapsed() string {
	if s.TimeElapsed <= 0 {
		return "--:--"
	}
	hours := s.TimeElapsed / 3600
	minutes := (s.TimeElapsed % 3600) / 60
	return formatTime(hours, minutes)
}

func formatTime(hours, minutes int) string {
	if hours > 0 {
		return formatDuration(hours, "h", minutes, "m")
	}
	return formatDuration(minutes, "m", 0, "")
}

func formatDuration(value1 int, unit1 string, value2 int, unit2 string) string {
	if value2 > 0 {
		return formatWithUnit(value1, unit1) + " " + formatWithUnit(value2, unit2)
	}
	return formatWithUnit(value1, unit1)
}

func formatWithUnit(value int, unit string) string {
	if unit == "" {
		return ""
	}
	return string(rune(value/10+'0')) + string(rune(value%10+'0')) + unit
}
