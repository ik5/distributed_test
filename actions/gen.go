package actions

import (
	"fmt"

	"github.com/ik5/distributed_test/types"
)

// GenerateActionRequest generates a full header for connection request
// Body is the generated JSON content
func GenerateActionRequest(t types.ActionType, body []byte) []byte {
	l := len(body)
	result := make([]byte, l)
	result = []byte(fmt.Sprintf("l: %d\r\n", l))
	result = append(result, []byte(fmt.Sprintf("t: %d\r\n", t))...)
	result = append(result, '\r', '\n')
	result = append(result, body...)
	return result
}
