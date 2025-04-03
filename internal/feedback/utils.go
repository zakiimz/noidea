package feedback

import (
	"time"
)

// GetTimeOfDay returns a string representation of the time of day
func GetTimeOfDay(t time.Time) string {
	hour := t.Hour()
	
	switch {
	case hour >= 5 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 17:
		return "afternoon"
	case hour >= 17 && hour < 21:
		return "evening"
	default:
		return "night"
	}
} 