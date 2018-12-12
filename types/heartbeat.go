package types

import "time"

// Heartbeat holds the heartbeat message structure
type Heartbeat struct {
	Action string    `json:"action"`
	From   string    `json:"from"`
	Time   time.Time `json:"time"`
}
