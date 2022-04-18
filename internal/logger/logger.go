package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"io"
)

func InitLogger(logFile io.Writer) error {
	log.SetOutput(logFile)
	log.SetReportCaller(true)

	log.SetFormatter(&nested.Formatter{})

	return nil
}
