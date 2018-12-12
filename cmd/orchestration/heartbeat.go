package main

import (
	"time"

	"github.com/ik5/distributed_test/heartbeat"
	"github.com/ik5/distributed_test/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func sendHeartbeat() {
	log.Logger.Trace("Initialize sendHeartbeat")
	delay := time.Duration(viper.GetInt("heartbeat_every")) * time.Millisecond

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	addr := viper.GetString("parent")
	log.Logger.WithFields(logrus.Fields{
		"delay": delay,
		"addr":  addr,
	}).Trace("Have needed fields")
	for t := range ticker.C {

		log.Logger.WithFields(logrus.Fields{
			"ticker": t,
		}).Tracef("Sending heartbeat to: %s", addr)
		result, err := heartbeat.SendHeartbeat(addr)
		if err != nil {
			// wait for some time before retry
			log.Logger.Errorf("Unable to send heartbeat: %s", err)
			time.Sleep(delay + (5 * time.Second))
			continue
		}
		log.Logger.Debugf("Sent heartbeat, status: %t", result)
	}
}
