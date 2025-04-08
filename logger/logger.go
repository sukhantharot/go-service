package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Set output to stdout
	log.Out = os.Stdout

	// Set log level based on environment
	if os.Getenv("APP_ENV") == "production" {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	// Set formatter
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}

// Fields type for structured logging
type Fields map[string]interface{}

// Debug logs a debug message
func Debug(msg string, fields Fields) {
	log.WithFields(logrus.Fields(fields)).Debug(msg)
}

// Info logs an info message
func Info(msg string, fields Fields) {
	log.WithFields(logrus.Fields(fields)).Info(msg)
}

// Warn logs a warning message
func Warn(msg string, fields Fields) {
	log.WithFields(logrus.Fields(fields)).Warn(msg)
}

// Error logs an error message
func Error(msg string, err error, fields Fields) {
	if fields == nil {
		fields = Fields{}
	}
	fields["error"] = err.Error()
	log.WithFields(logrus.Fields(fields)).Error(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error, fields Fields) {
	if fields == nil {
		fields = Fields{}
	}
	fields["error"] = err.Error()
	log.WithFields(logrus.Fields(fields)).Fatal(msg)
} 