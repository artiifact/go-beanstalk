package beanstalk

import "fmt"

type LogLevel int

const (
	DebugLogLevel LogLevel = iota + 1
	InfoLogLevel
	WarningLogLevel
	ErrorLogLevel
	PanicLogLevel
	FatalLogLevel
)

func (l LogLevel) String() string {
	switch l {
	case DebugLogLevel:
		return "debug"
	case InfoLogLevel:
		return "info"
	case WarningLogLevel:
		return "warning"
	case ErrorLogLevel:
		return "error"
	case PanicLogLevel:
		return "panic"
	case FatalLogLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

type Logger interface {
	Log(level LogLevel, msg string, args map[string]interface{})
}

// nop logger

type nopLogger struct{}

func (l *nopLogger) Log(LogLevel, string, map[string]interface{}) {}

var NopLogger = &nopLogger{}
