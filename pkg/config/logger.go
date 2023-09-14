package config

import (
	"log/slog"
	"os"
)

// InitLogger initialize the default app logger using the AppConfig to set the
// current log level
func InitLogger(config AppConfig) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.LogLevel)); err != nil {
		slog.Warn("Invalid log level configuration")
		level = slog.LevelWarn
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// InterceptorLogger adapts slog logger to interceptor logger.
// func InterceptorLogger(l *slog.Logger) logging.Logger {
// 	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
// 		l.Log(ctx, slog.Level(lvl), msg, fields...)
// 	})
// }
