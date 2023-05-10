package log

import "fmt"

type LevelLog uint32

const (
	PanicLevel LevelLog = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (level LevelLog) MarshalText() (string, error) {
	switch level {
	case TraceLevel:
		return "trace", nil
	case DebugLevel:
		return "debug", nil
	case InfoLevel:
		return "info", nil
	case WarnLevel:
		return "warning", nil
	case ErrorLevel:
		return "error", nil
	case FatalLevel:
		return "fatal", nil
	case PanicLevel:
		return "panic", nil
	}
	return "", fmt.Errorf("not a valid logrus level %d", level)
}
