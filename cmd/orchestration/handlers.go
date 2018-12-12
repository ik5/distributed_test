package main

import (
	"encoding/json"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ik5/distributed_test/actions"
	"github.com/ik5/distributed_test/clients"
	"github.com/ik5/distributed_test/log"
	"github.com/ik5/distributed_test/types"
	"github.com/spf13/viper"
)

func handleHeartbeat(conn net.Conn, proto types.RawAction) {
	var hb types.Heartbeat
	err := json.Unmarshal(proto.Body, &hb)
	if err != nil {
		log.Logger.Errorf("Unable to parse heartbeat: %s", err)
		return
	}
	addr := strings.Split(hb.From, ":")
	clients.UpdateNodeHeartbeat(addr[0], hb.Time)
	conn.Close()
}

func execRequest(addr string, actionList types.Action, refCounter *uint32) {
	defer func() { atomic.StoreUint32(refCounter, *refCounter-1) }()
	log.Logger.Tracef("Going to send action to %s", addr)
	conn, err := net.DialTimeout("tcp", addr, types.DefaultTimeout)
	if err != nil {
		log.Logger.Errorf("Error open connection for action: %s", err)
		return
	}
	defer conn.Close()

	buff, _ := json.Marshal(actionList)
	conn.SetWriteDeadline(time.Now().Add(types.DefaultWriteTimeout))
	_, err = conn.Write(actions.GenerateActionRequest(types.ActionRequest, buff))
	if err != nil {
		log.Logger.Errorf("Error sending action to connection: %s", err)
		return
	}
	log.Logger.Info("Action were sent")
	readBuff := make([]byte, types.MaxBufferSize)
	_, err = conn.Read(readBuff)
	if err != nil {
		log.Logger.Tracef("Unable to read readBuff: %s", err)
		return
	}
	protocol := actions.Parse(readBuff)
	if protocol.Type != types.ActionAnswer {
		log.Logger.Errorf("Unknown protocol type: %d", protocol.Type)
		return
	}
	var actions []types.CommandResult
	err = json.Unmarshal(protocol.Body, &actions)
	if err != nil {
		log.Logger.Errorf("unable to parse response: %s", err)
		return
	}
	conn.Write(readBuff)
	cleanAddr := strings.Split(addr, ":")
	for _, answer := range actions {
		clients.UpdateNodeCommandResult(cleanAddr[0], answer)
	}

	log.Logger.Infof("Updated result for client %s (uptime: %s)", cleanAddr[0], clients.NodeUptime(cleanAddr[0]))
}

func handleActionRequest(conn net.Conn, proto types.RawAction) {
	var action types.Action
	err := json.Unmarshal(proto.Body, &action)
	if err != nil {
		log.Logger.Errorf("Unable to parse action: %s", err)
		return
	}

	clientsList := viper.GetStringSlice("nodes")
	refCount := uint32(0)
	for _, clientAddr := range clientsList {
		cleanAddr := strings.Split(clientAddr, ":")
		clientRec := clients.GetNode(cleanAddr[0])
		if (clientRec.LastContact == time.Time{}) {
			log.Logger.Warnf("Client %s never contacted the server", cleanAddr[0])
			continue
		}

		lastContact := time.Now().Sub(clientRec.LastContact)
		if lastContact > (viper.GetDuration("node_uptime") * time.Second) {
			log.Logger.Warnf("Client %s last contact was too long (%s)", cleanAddr[0], lastContact)
			continue
		}
		atomic.AddUint32(&refCount, 1)
		go execRequest(clientAddr, action, &refCount)
	}

	go notifyParentWithRef(&refCount)
}

func notifyParent() {
	addr := viper.GetString("parent")
	_, err := actions.SendAggreatedActions(addr, clients.GetNodes())
	if err != nil {
		log.Logger.Errorf("Unable to notify parent %s: %s", addr, err)
		return
	}

	log.Logger.Infof("Parent %s was notified about actions", addr)
}

func notifyParentWithRef(refCount *uint32) {
	for {
		if *refCount == 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	notifyParent()
}

func handleActionAnswer(conn net.Conn, proto types.RawAction) {
	notifyParent()
}
