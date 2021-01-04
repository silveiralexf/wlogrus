/*Package wlogrus is a simple wrapper for the incredible logrus package to facilitate formatting and forwarding
structured logs into a central stash or to the standard output
*/
package wlogrus

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Message holds the existing log message fields
type Message struct {
	Severity string
	Tag      string
	Body     interface{}
	Location string
}

type logFormatter struct {
	logrus.TextFormatter
}

// Formats log output with timestamp, colors and proper severity
func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int

	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 44 // white letters + purple bg
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel:
		levelColor = 31 // red
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 41 // white letters + red bg
	default:
		levelColor = 36 // blue
	}
	return []byte(
		fmt.Sprintf(
			"%s [\x1b[%dm%s\x1b[0m] %s\n",
			entry.Time.Format(f.TimestampFormat),
			levelColor,
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		)), nil
}

// Info will log informational messages to stdout
func Info(tag string, message interface{}) {
	consoleOut(Message{Severity: "INFO", Tag: tag, Body: message})
}

// Warn will log warning messages to stdout
func Warn(tag string, message interface{}) {
	consoleOut(Message{Severity: "WARNING", Tag: tag, Body: message})
}

// Error will log Error messages to stdout and logs and provide information on runtime caller,
// preferably through 'wlogrus.CallerInfo()' function
func Error(tag string, err interface{}, callerInfo string) {
	consoleOut(Message{Severity: "ERROR", Tag: tag, Body: err, Location: callerInfo})
}

// Fatal will log Fatal messages to stdout and logs and provide information on runtime caller,
// preferably through 'wlogrus.CallerInfo()' function before exiting with an error
func Fatal(tag string, err interface{}, callerInfo string) {
	consoleOut(Message{Severity: "FATAL", Tag: tag, Body: err, Location: callerInfo})
}

// Debug will log informational additional messages to stdout and logs when debug flag is enabled
func Debug(tag string, message interface{}, callerInfo string) {
	consoleOut(Message{Severity: "DEBUG", Tag: tag, Body: message, Location: callerInfo})
}

func consoleOut(m Message) {
	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout),
		Level: logrus.InfoLevel,
		Formatter: &logFormatter{
			logrus.TextFormatter{
				FullTimestamp:          true,
				TimestampFormat:        "2006-01-02 15:04:05",
				ForceColors:            true,
				DisableLevelTruncation: true,
			},
		},
	}

	if os.Getenv("WLOGRUS_JSON") == "true" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	if os.Getenv("WLOGRUS_DEBUG") == "true" {
		logger.Level = logrus.DebugLevel
	}

	ctx := context.Background()
	entry := fmt.Sprintf("[%v] %v", m.Tag, m.Body)

	switch m.Severity {
	case "DEBUG":
		logger.WithFields(
			logrus.Fields{"severity": m.Severity, "tag": m.Tag, "body": m.Body, "caller": m.Location},
		).WithContext(ctx).Debug(entry)
	case "INFO":
		logger.WithFields(
			logrus.Fields{"severity": m.Severity, "tag": m.Tag, "body": m.Body},
		).WithContext(ctx).Info(entry)
	case "WARNING":
		logger.WithFields(
			logrus.Fields{"severity": m.Severity, "tag": m.Tag, "body": m.Body},
		).WithContext(ctx).Warning(entry)
	case "ERROR":
		entry := fmt.Sprintf("%v [%v]", entry, m.Location)
		logger.WithFields(
			logrus.Fields{"severity": m.Severity, "tag": m.Tag, "body": m.Body, "caller": m.Location},
		).WithContext(ctx).Error(entry)
	case "FATAL":
		entry := fmt.Sprintf("%v [%v]", entry, m.Location)
		logger.WithFields(
			logrus.Fields{"severity": m.Severity, "tag": m.Tag, "body": m.Body, "caller": m.Location},
		).WithContext(ctx).Fatal(entry)
	default:
		logger.WithContext(ctx).Info(entry)
	}
}

// CallerInfo returns function name, file and line number which invoked it
func CallerInfo(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s.%s:%d", filepath.Base(file), filepath.Base(runtime.FuncForPC(function).Name()), line)
}
