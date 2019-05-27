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

// Logger interface to represent a logging facility
type Logger interface {
	NewSubLogger(name string) Logger
	Info(kv ...interface{})
	Warn(kv ...interface{})
	Error(kv ...interface{})
	Debug(kv ...interface{})
	Fatal(kv ...interface{})
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

// logMessage -- universal log function
func (dl *DefaultLogger) logMessage(level string, kv []interface{}) {
	if len(kv)%2 != 0 {
		dl.Error("error", "the number of keys/values should be even", "kvs", kv)
		return
	}

	if len(kv) == 0 {
		dl.Error("error", "at least 1 field pair must be provided")
		return
	}

	logMeta := make(map[string]interface{})

	for i := 0; i < len(kv); i += 2 {
		key := kv[i].(string)
		logMeta[key] = kv[i+1]
	}

	meta, err := json.Marshal(logMeta)
	if err != nil {
		dl.Error("msg", "unable to marshal", "error", err)
		return
	}

	header := fmt.Sprintf(`{"time":"%v","name":"%s","pid":%d,"level":"%s",`,
		time.Now().Format(time.UnixDate), dl.name, os.Getpid(), level)

	dl.Out.Println(header + string(meta)[1:])
}

// Info logger
func (dl DefaultLogger) Info(kv ...interface{}) {
	dl.logMessage("info", kv)
}

// Warn logger
func (dl DefaultLogger) Warn(kv ...interface{}) {
	dl.logMessage("warn", kv)
}

// Error logger
func (dl DefaultLogger) Error(kv ...interface{}) {
	dl.logMessage("error", kv)
}

// Fatal logger
func (dl DefaultLogger) Fatal(kv ...interface{}) {
	dl.logMessage("fatal", kv)
	os.Exit(1)
}

// Debug logger
func (dl DefaultLogger) Debug(kv ...interface{}) {
	dl.logMessage("debug", kv)
}
