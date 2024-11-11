package interceptors

import (
	"context"

	"github.com/elliot14A/jondev/infrastructure/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

// interceptorLogger adapts our Logger to the logging.Logger interface
func interceptorLogger(l *logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, fields...)
		case logging.LevelInfo:
			l.Info(msg, fields...)
		case logging.LevelWarn:
			l.Warn(msg, fields...)
		case logging.LevelError:
			l.Error(msg, fields...)
		default:
			l.Error("unknown log level", "level", lvl, "msg", msg)
		}
	})
}

// logTraceID extracts trace ID from context if present

// SetupGRPCServerWithLogging creates a new gRPC server with logging and tracing middleware
func SetupGRPCServerWithLogging(logger *logger.Logger) *grpc.Server {
	// Create server with logging middleware
	server := grpc.NewServer(
		// Add OpenTelemetry gRPC instrumentation
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(
				interceptorLogger(logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(
				interceptorLogger(logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
		),
	)

	return server
}
