package types

import "time"

// NodeDetails holds information for connection
type NodeDetails struct {
	Addr         string          `json:"addr"`
	FirstContact time.Time       `json:"first_contact"`
	LastContact  time.Time       `json:"last_contact"`
	Actions      []CommandResult `json:"actions"`
}

// Nodes holds nodes details for all known nodes
type Nodes map[string]NodeDetails
