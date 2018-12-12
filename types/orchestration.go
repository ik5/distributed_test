package types

import "time"

// OrchDetails holds information about an orchestration server
type OrchDetails struct {
	Addr         string    `json:"addr"`
	FirstContact time.Time `json:"first_contact"`
	LastContact  time.Time `json:"last_contact"`
	Nodes        Nodes     `json:"nodes"`
}

// Orchs holds a list of all known orchestration servers
type Orchs map[string]OrchDetails
