package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cockroachdb/errors"
)

const (
	modNum         = 0644
	stackTraceSize = 10
)

type Logger interface {
	Debug(message string, v ...interface{})
	Info(message string, v ...interface{})
	Warn(message string, v ...interface{})
	Error(message string, err error, v ...interface{})
}

type jsonLogger struct {
	logger *slog.Logger
}

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

var logger Logger

func Debug(message string, v ...interface{}) {
	getLogger().Debug(message, v...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func NewLogger() error {
	var err error
	logger, err = newJSONLogger()
	return err
}

func (l *jsonLogger) Debug(message string, v ...interface{}) {
	l.write(DEBUG, message, v...)
}

func (l *jsonLogger) Info(message string, v ...interface{}) {
	l.write(INFO, message, v...)
}

func (l *jsonLogger) Warn(message string, v ...interface{}) {
	l.write(WARN, message, v...)
}

func (l *jsonLogger) Error(message string, err error, v ...interface{}) {
	var errs []interface{}
	if err != nil {
		errs = append([]interface{}{err}, v...)
	} else {
		errs = v
	}
	l.write(ERROR, message, errs...)
}

func (l *jsonLogger) write(level LogLevel, message string, v ...interface{}) {
	if len(v) > 0 {
		if err, ok := v[0].(error); ok {

			rest := v[1:]
			if stackTrace := errors.GetReportableStackTrace(err); stackTrace != nil {
				frames := stackTrace.Frames
				num := len(frames)
				if num > stackTraceSize {
					frames = frames[num-stackTraceSize : num]
				}
				v = append([]interface{}{"stack_trace", frames}, rest...)
			} else {

				v = append([]interface{}{"error", err.Error()}, rest...)
			}
		}
	}

	switch level {
	case DEBUG:
		l.logger.Debug(message, v...)
	case INFO:
		l.logger.Info(message, v...)
	case WARN:
		l.logger.Warn(message, v...)
	case ERROR:
		l.logger.Error(message, v...)
	}
}

func newJSONLogger() (Logger, error) {
	var handler slog.Handler
	if os.Getenv("LOG_TYPE") == "file" {
		logFilePath := os.Getenv("LOG_FILE_PATH")
		if logFilePath == "" {
			logFilePath = "logs/gin.log"
		}
		// Create or open log file
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, modNum)
		if err != nil {
			return nil, fmt.Errorf("Failed to open log file: %v", err)
		}
		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	return &jsonLogger{logger: slog.New(handler)}, nil
}

func getLogger() Logger {
	if logger == nil {
		if err := NewLogger(); err != nil {
			panic("Failed to initialize the log.")
		}
	}
	return logger
}
