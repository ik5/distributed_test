package actions

import (
	"encoding/json"
	"testing"

	"github.com/ik5/distributed_test/types"
)

func TestParser(t *testing.T) {
	hbt := "t: 1\r\nl: 92\r\n\r\n{\"action\":\"heartbeat\",\"from\":\"127.0.0.1:46950\",\"time\":\"2018-12-01T21:22:22.680336737+02:00\"}"
	parsed := Parse([]byte(hbt))

	if parsed.Length != 92 {
		t.Errorf("parsed.Length == %d", parsed.Length)
	}

	if parsed.Type != types.ActionHeartbeat {
		t.Errorf("parsed.Type == %d", parsed.Type)
	}

	var hb types.Heartbeat
	err := json.Unmarshal(parsed.Body, &hb)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if hb.Action != "heartbeat" {
		t.Errorf("Invalid action: %s", hb.Action)
	}

	if hb.From != "127.0.0.1:46950" {
		t.Errorf("Invalid From: %s", hb.From)
	}

	if hb.Time.String() != "2018-12-01 21:22:22.680336737 +0200 IST" {
		t.Errorf("Invalid Time: %s", hb.Time)
	}
}
