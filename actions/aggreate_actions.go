package actions

import (
	"encoding/json"
	"net"
	"time"

	"github.com/ik5/distributed_test/types"
)

// SendAggreatedActions send to an address the given actions
func SendAggreatedActions(addr string, nodes types.Nodes) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, types.DefaultTimeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	buff, _ := json.Marshal(nodes)
	conn.SetWriteDeadline(time.Now().Add(types.DefaultWriteTimeout))
	_, err = conn.Write(GenerateActionRequest(types.ActionAnswer, buff))
	if err != nil {
		return false, err
	}
	return true, nil
}
