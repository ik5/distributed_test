package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ik5/distributed_test/config"
	"github.com/ik5/distributed_test/log"
	"github.com/spf13/viper"
)

var logFile *os.File

func initLogger() {
	debug := viper.GetString("env") == "debug"
	log.InitLog(viper.GetString("syslog_socket_type"),
		viper.GetString("syslog_address"),
		viper.GetString("syslog_tag"),
		viper.GetString("log_level"),
		viper.GetBool("use_syslog"),
		config.SyslogLevel(),
		debug,
	)

	var err error
	logFile, err = os.OpenFile(viper.GetString("log_file"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Logger.SetOutput(os.Stdout)
	} else {
		log.Logger.SetOutput(logFile)
	}
}

func initApp() {
	configPath := []string{"."}
	err := config.Init("yaml", "config", configPath)
	if err != nil {
		panic(err)
	}
	initLogger()
}

func main() {
	initApp()
	if logFile != nil {
		defer logFile.Close()
	}

	log.Logger.Info("Starting application")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stop := make(chan bool, 1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			log.Logger.Info("Shutdown")
			stop <- true
		case <-ctx.Done():
			log.Logger.Info("Shutdown")
			stop <- true
		}
	}()

	go listen()
	go sendHeartbeat()
	<-stop
	log.Logger.Info("Bye")
}
