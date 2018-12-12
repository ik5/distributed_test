package clients

import (
	"sync"
	"time"

	"github.com/ik5/distributed_test/types"
)

// while Go v1.9 added sync.Map, making sure that this code can work also on
// older versions as well
var nodes types.Nodes
var nodesMutex sync.RWMutex

// UpdateNodeHeartbeat updates a node details
func UpdateNodeHeartbeat(addr string, t time.Time) {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()
	if nodes == nil {
		nodes = make(types.Nodes)
	}

	_, found := nodes[addr]
	if !found {
		nodes[addr] = types.NodeDetails{
			Addr:         addr,
			FirstContact: t,
			LastContact:  t,
		}
	} else {
		details := nodes[addr]
		details.LastContact = t
		nodes[addr] = details
	}
}

// NodeUptime calculate how long is the up
func NodeUptime(addr string) time.Duration {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()
	if nodes == nil {
		nodes = make(types.Nodes)
		return 0
	}
	node, found := nodes[addr]
	if !found {
		return 0
	}
	return node.LastContact.Sub(node.FirstContact)
}

// UpdateNodeCommandResult update command result to a node
func UpdateNodeCommandResult(addr string, action types.CommandResult) {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()
	if nodes == nil {
		nodes = make(types.Nodes)
	}

	node, found := nodes[addr]
	if !found {
		node.Addr = addr
		node.FirstContact = time.Now()
		node.LastContact = time.Now()
	}

	node.Actions = append(node.Actions, action)
	nodes[addr] = node
}

// GetNodes return a copy of nodes
func GetNodes() types.Nodes {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()
	return nodes
}

// GetNode return an information about a single node
func GetNode(addr string) types.NodeDetails {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()
	if nodes == nil {
		return types.NodeDetails{}
	}
	return nodes[addr]
}
