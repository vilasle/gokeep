package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func InterceptorLogger() logging.Logger {
	expectedFields := map[string]struct{}{
		"grpc.service":    {}, //server
		"grpc.method":     {}, //proto.MetricService
		"grpc.address":    {},
		"grpc.start_time": {},
		"grpc.time_ms":    {},
	}
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]any, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			if _, ok := expectedFields[key.(string)]; !ok {
				continue
			}

			f = append(f, key.(string), value)
		}
		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg, f...)
		case logging.LevelInfo:
			logger.Info(msg, f...)
		case logging.LevelWarn:
			logger.Warn(msg, f...)
		case logging.LevelError:
			logger.Error(msg, f...)
		default:
			logger.Error("unknown level", "level", lvl)
		}
	})
}
