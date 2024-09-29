package helpers

import (
	"io"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func InitializeLogs() *logrus.Logger {
	log := logrus.New()

	logFile := &lumberjack.Logger{
		Filename:   "logs/order-service.log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     1,
		Compress:   false,
	}

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	log.SetFormatter(&logrus.JSONFormatter{})

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch logLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
