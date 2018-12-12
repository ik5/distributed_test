package log

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Logger is the instance for the log system
var Logger *logrus.Logger

// InitLog initialize logging facility
func InitLog(socketType, address, tag, level string, useSyslog bool, priority syslog.Priority, debug bool) {
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{PrettyPrint: false}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.SetLevel(logrus.TraceLevel)
	} else {
		Logger.SetLevel(logLevel)
	}
	Logger.SetReportCaller(debug)

	if !useSyslog {
		return
	}
	hook, err := logrus_syslog.NewSyslogHook(socketType, address, priority, tag)
	if err != nil {
		Logger.Errorf("Unable to use syslog: %s", err)
	} else {
		Logger.AddHook(hook)
	}
}
