package main

import (
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/ik5/distributed_test/actions"
	"github.com/ik5/distributed_test/clients"
	"github.com/ik5/distributed_test/log"
	"github.com/ik5/distributed_test/types"
	"github.com/sirupsen/logrus"
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
	log.Logger.Tracef("Updaeting %s heartbeat", hb.From)
	clients.UpdateOrchHeartbeat(addr[0], hb.Time)
	conn.Close()
}

func execRequest(addr string, actionList types.Action) {
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
	log.Logger.Info("Action was sent")
}

func handleActionRequest(conn net.Conn, proto types.RawAction) {
	var action types.Action
	err := json.Unmarshal(proto.Body, &action)
	if err != nil {
		log.Logger.Errorf("Unable to parse action: %s", err)
		return
	}

	clientsList := viper.GetStringSlice("orchestrations")
	for _, clientAddr := range clientsList {
		cleanAddr := strings.Split(clientAddr, ":")
		clientRec := clients.GetOrch(cleanAddr[0])
		if (clientRec.LastContact == time.Time{}) {
			log.Logger.Warnf("Client %s never contacted the server", cleanAddr[0])
			continue
		}

		lastContact := time.Now().Sub(clientRec.LastContact)
		if lastContact > (viper.GetDuration("orchestration_uptime") * time.Second) {
			log.Logger.Warnf("Client %s last contact was too long (%s)", cleanAddr[0], lastContact)
			continue
		}
		go execRequest(clientAddr, action)
	}
}

func handleActionAnswer(conn net.Conn, proto types.RawAction) {
	orchAddr := conn.RemoteAddr()
	addr := strings.Split(orchAddr.String(), ":")[0]
	var nodes types.Nodes
	err := json.Unmarshal(proto.Body, &nodes)
	if err != nil {
		log.Logger.Tracef("Unable to parse body: %s", err)
		return
	}
	log.Logger.WithFields(logrus.Fields{
		"nodes": nodes,
	}).Trace("Have action answer nodes")
	clients.UpdateOrchNodes(addr, nodes)
	log.Logger.Infof("Updated action answers from %s", addr)
}
