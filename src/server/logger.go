package server

import (
	"os"
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func InitLogger() {

	log.Formatter = new(logrus.TextFormatter) // default

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.Level = logrus.DebugLevel
}
