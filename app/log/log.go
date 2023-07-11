package log

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

type Fields map[string]interface{}


func init() {
	var loglevel string

	flag.StringVar(&loglevel, "log level", "Trace", "Specify desired log level")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableQuote:  true,
		FullTimestamp: true})
	log.SetOutput(os.Stdout)

	switch loglevel {
	case "Trace":
		log.SetLevel(log.TraceLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	case "Fatal":
		log.SetLevel(log.FatalLevel)
	case "Panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func WithFields(fields Fields) *log.Entry {
	return log.WithFields(log.Fields(fields))
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	log.Trace(args...)
}

func Tracef(string string, args ...interface{}) {
	log.Tracef(string, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(string string, args ...interface{}) {
	log.Debugf(string, args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	log.Print(args...)
}

func Printf(string string, args ...interface{}) {
	log.Printf(string, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(string string, args ...interface{}) {
	log.Infof(string, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

func DWarnf(string string, args ...interface{}) {
	log.Warnf(string, args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	log.Warning(args...)
}

func Warningf(string string, args ...interface{}) {
	log.Warningf(string, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(string string, args ...interface{}) {
	log.Errorf(string, args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	log.Panic(args...)
}

func DPanicf(string string, args ...interface{}) {
	log.Panicf(string, args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(string string, args ...interface{}) {
	log.Fatalf(string, args...)
}
