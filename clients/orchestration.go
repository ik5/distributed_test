package clients

import (
	"sync"
	"time"

	"github.com/ik5/distributed_test/types"
)

// while Go v1.9 added sync.Map, making sure that this code can work also on
// older versions as well
var orchs types.Orchs
var orchsMutex sync.RWMutex

// UpdateOrchHeartbeat update the orchestration server heartbeat
func UpdateOrchHeartbeat(addr string, t time.Time) {
	orchsMutex.Lock()
	defer orchsMutex.Unlock()
	if orchs == nil {
		orchs = make(types.Orchs)
	}

	orch, found := orchs[addr]
	if !found {
		orch = types.OrchDetails{
			Addr:         addr,
			FirstContact: t,
			LastContact:  t,
		}
	} else {
		orch.LastContact = t
	}

	orchs[addr] = orch
}

// UpdateOrchNodes updates the nodes for a given orchestration server
func UpdateOrchNodes(addr string, nodes types.Nodes) {
	orchsMutex.Lock()
	defer orchsMutex.Unlock()
	if orchs == nil {
		orchs = make(types.Orchs)
	}

	orch, found := orchs[addr]
	if !found {
		orch.Addr = addr
		orch.FirstContact = time.Now()
		orch.LastContact = time.Now()
	}

	if orch.Nodes == nil {
		orch.Nodes = nodes
		orchs[addr] = orch
		return
	}

	for key, orchVal := range orch.Nodes {
		node, found := nodes[key]
		if !found {
			continue
		}
		// orchestration server is the source of truth in this case, even if it was
		// down because it is the direct connection for the node
		orchVal.FirstContact = node.FirstContact
		orchVal.LastContact = node.LastContact
		orchVal.Actions = append(orchVal.Actions, node.Actions...)

		orch.Nodes[key] = orchVal
	}

	orchs[addr] = orch
}

//GetOrchs returns all current known orchestration servers
func GetOrchs() types.Orchs {
	orchsMutex.Lock()
	defer orchsMutex.Unlock()

	return orchs
}

// GetOrch returns a given orch by it's address
func GetOrch(addr string) types.OrchDetails {
	orchsMutex.Lock()
	defer orchsMutex.Unlock()
	if orchs == nil {
		return types.OrchDetails{}
	}
	return orchs[addr]
}
