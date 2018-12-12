package config

import (
	"log/syslog"

	"github.com/spf13/viper"
)

// Init initialize the config file usage
func Init(confType, name string, path []string) error {
	viper.SetConfigType(confType)
	viper.SetConfigName(name)
	for _, loc := range path {
		viper.AddConfigPath(loc)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// SyslogLevel calculate the priority and facility
func SyslogLevel() syslog.Priority {
	data := map[string]syslog.Priority{
		"emerg":    syslog.LOG_EMERG,
		"alert":    syslog.LOG_ALERT,
		"crit":     syslog.LOG_CRIT,
		"err":      syslog.LOG_ERR,
		"warning":  syslog.LOG_WARNING,
		"notice":   syslog.LOG_NOTICE,
		"info":     syslog.LOG_INFO,
		"debug":    syslog.LOG_DEBUG,
		"kern":     syslog.LOG_KERN,
		"user":     syslog.LOG_USER,
		"mail":     syslog.LOG_MAIL,
		"daemon":   syslog.LOG_DAEMON,
		"auth":     syslog.LOG_AUTH,
		"syslog":   syslog.LOG_SYSLOG,
		"lpr":      syslog.LOG_LPR,
		"news":     syslog.LOG_NEWS,
		"uucp":     syslog.LOG_UUCP,
		"cron":     syslog.LOG_CRON,
		"authpriv": syslog.LOG_AUTHPRIV,
		"ftp":      syslog.LOG_FTP,
		"local0":   syslog.LOG_LOCAL0,
		"local1":   syslog.LOG_LOCAL1,
		"local2":   syslog.LOG_LOCAL2,
		"local3":   syslog.LOG_LOCAL3,
		"local4":   syslog.LOG_LOCAL4,
		"local5":   syslog.LOG_LOCAL5,
		"local6":   syslog.LOG_LOCAL6,
		"local7":   syslog.LOG_LOCAL7,
	}

	return data[viper.GetString("syslog_level")] | data[viper.GetString("syslog_facility")]
}
