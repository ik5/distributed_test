package main

import (
	"net"

	"github.com/ik5/distributed_test/actions"
	"github.com/ik5/distributed_test/log"
	"github.com/ik5/distributed_test/server"
	"github.com/ik5/distributed_test/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handler(conn net.Conn, buf []byte) {
	log.Logger.WithFields(logrus.Fields{
		"buffer": string(buf),
	}).Trace("Got callback with buffer")

	protocol := actions.Parse(buf)
	log.Logger.Tracef("Protocol: %d, %s, %d %s",
		protocol.Length,
		protocol.Type,
		len(protocol.Body),
		string(protocol.Body),
	)
	switch protocol.Type {
	case types.ActionHeartbeat:
		handleHeartbeat(conn, protocol)
	case types.ActionRequest:
		handleActionRequest(conn, protocol)
	case types.ActionAnswer:
		handleActionAnswer(conn, protocol)
	default:
		log.Logger.Warnf("Unknown type: %d", protocol.Type)
		return
	}
}

func listen() {
	addr := viper.GetString("listen")
	log.Logger.Infof("Going to start a listening on: %s", addr)
	err := server.Listen(addr, handler)
	if err != nil {
		log.Logger.Errorf("Unable to listen: %s", err)
		return
	}
}
