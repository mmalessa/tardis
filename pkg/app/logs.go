package app

import (
	log "github.com/sirupsen/logrus"
)

func InitLogs(environment string) error {
	switch environment {
	case "dev":
		return initLogsEnvDev()
	default:
		return initLogsEnvProd()
	}
	return nil
}

func initLogsEnvDev() error {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
	return nil
}

func initLogsEnvProd() error {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:     false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.InfoLevel)
	return nil
}
