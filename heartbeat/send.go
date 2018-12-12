package heartbeat

import (
	"encoding/json"
	"net"
	"time"

	"github.com/ik5/distributed_test/actions"
	"github.com/ik5/distributed_test/types"
)

// SendHeartbeat sends a heartbeat request to a server
// addr should looke like so:  ip:port
// The function returns error if one found, otherwise, true if the server responded
func SendHeartbeat(addr string) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, types.DefaultTimeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	hb := types.Heartbeat{
		Action: "heartbeat",
		From:   conn.LocalAddr().String(),
		Time:   time.Now(),
	}

	buff, _ := json.Marshal(hb)
	conn.SetWriteDeadline(time.Now().Add(types.DefaultWriteTimeout))
	_, err = conn.Write(actions.GenerateActionRequest(types.ActionHeartbeat, buff))
	if err != nil {
		return false, err
	}

	return true, nil
}
