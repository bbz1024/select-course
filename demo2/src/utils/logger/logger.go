package logger

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"select-course/demo2/src/constant/config"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()

	switch config.EnvCfg.LoggerLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN", "WARNING":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	}

	filePath := path.Join("/var", "log", "select-course", "select-course.log")
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	Logger = log.WithFields(log.Fields{
		"Hostname": hostname,
	})
}

var Logger *log.Entry

func LogService(name string) *log.Entry {
	return Logger.WithFields(log.Fields{
		"Service": name,
	})
}
