package main

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/ik5/distributed_test/actions"
	"github.com/ik5/distributed_test/log"
	"github.com/ik5/distributed_test/types"
)

func doAction(taskName, path, typ, check string) (bool, error) {
	log.Logger.Debugf("Going over task %s", taskName)
	switch typ {
	case "file_contains":
		return actions.IsFileContains(path, check)
	case "file_exists":
		return actions.IsFileExists(path)
	case "proc_running":
		return actions.ProcExists(path)
	}
	return false, errors.New("Unsupported type")
}

func handleActionRequest(conn net.Conn, proto types.RawAction) {
	var action types.Action
	err := json.Unmarshal(proto.Body, &action)
	if err != nil {
		log.Logger.Errorf("Unable to parse action: %s", err)
		return
	}

	var results []types.CommandResult
	for name, task := range action {
		result, err := doAction(name, task.Path, task.Type, task.Check)
		if err != nil {
			log.Logger.Errorf("Unable to execute %s: %s", name, err)
			continue
		}

		log.Logger.Debugf("Action %s finished with %t ", name, result)

		cmdResult := types.CommandResult{
			Path:   task.Path,
			Type:   task.Type,
			Check:  task.Check,
			Result: result,
		}
		results = append(results, cmdResult)
	}

	buff, _ := json.Marshal(results)
	conn.SetWriteDeadline(time.Now().Add(types.DefaultWriteTimeout))
	_, err = conn.Write(actions.GenerateActionRequest(types.ActionAnswer, buff))
	if err != nil {
		log.Logger.Errorf("Unable to send answer: %s", err)
	}

}
