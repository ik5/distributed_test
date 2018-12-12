package types

// Command holds the node actions to do
type Command struct {
	Path  string `json:"path"`
	Type  string `json:"type"`
	Check string `json:"check,omitempty"`
}

// Action holds a map of command
type Action = map[string]Command

// CommandResult holds the answer for a Command
type CommandResult struct {
	Path   string `json:"path"`
	Type   string `json:"type"`
	Check  string `json:"check,omitempty"`
	Result bool   `json:"result"`
}

// ActionResult holds a map of command result based on given action
type ActionResult = map[string]CommandResult

// ActionType holds the type of action to use
type ActionType int

// Action types
const (
	ActionUnknown ActionType = iota
	ActionHeartbeat
	ActionRequest
	ActionAnswer
)

func (a ActionType) String() string {
	str := ""
	switch a {
	case ActionUnknown:
		str = "Unknown"
	case ActionHeartbeat:
		str = "heartbeat"
	case ActionRequest:
		str = "action request"
	case ActionAnswer:
		str = "action answer"
	}
	return str
}
