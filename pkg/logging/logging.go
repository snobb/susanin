package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel constants
const (
	Info = iota
	Warn
	Error
	Debug
)

// LogLevel for log messages
type LogLevel int

// Meta stores the log meta data
type Meta map[string]interface{}

// Logger interface to represent a logging vacility
type Logger interface {
	NewSubLogger(name string) Logger
	Log(level LogLevel, meta Meta)
	Info(meta Meta)
	Error(meta Meta)
}

// DefaultLogger is a stdout implementation
type DefaultLogger struct {
	name string
	Out  *log.Logger
}

// New DefaultLogger instance
func New(name string, writer io.Writer) Logger {
	return DefaultLogger{
		name: strings.ToUpper(name),
		Out:  log.New(writer, "", 0),
	}
}

// NewSubLogger - create a sublogger
func (dl DefaultLogger) NewSubLogger(name string) Logger {
	subName := fmt.Sprintf("%s::%s", dl.name, strings.ToUpper(name))
	return New(subName, dl.Out.Writer())
}

// Log - universal log function
func (dl DefaultLogger) Log(logLevel LogLevel, meta Meta) {
	var level string

	switch logLevel {
	case Info:
		level = "info"
	case Warn:
		level = "warn"
	case Error:
		level = "error"
	case Debug:
		level = "debug"
	default:
		level = "info"
	}

	logMeta := map[string]interface{}{
		"time":  time.Now(),
		"name":  dl.name,
		"pid":   os.Getpid(),
		"level": level,
	}

	for k, v := range meta {
		logMeta[k] = v
	}

	logStr, err := json.Marshal(logMeta)
	if err != nil {
		dl.Out.Printf("Unable to marshal %v\n", logMeta)
	}

	dl.Out.Println(string(logStr))
}

// Info logger
func (dl DefaultLogger) Info(meta Meta) {
	dl.Log(Info, meta)
}

// Error logger
func (dl DefaultLogger) Error(meta Meta) {
	dl.Log(Error, meta)
}
