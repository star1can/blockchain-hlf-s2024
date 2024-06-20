package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Entry {
	level, _ := logrus.ParseLevel("debug")
	logrus.SetLevel(level)

	return logrus.WithFields(logrus.Fields{
		"app_name": "saac-v2 chaincode",
		"version":  "2",
	})
}
