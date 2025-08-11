package slog

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewLogger() *SlogLogger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format("2006/01/02 15:04:05")),
				}
			}
			return a
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	return &SlogLogger{logger: logger}
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
