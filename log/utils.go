package log

import (
	"context"
)

type contextKey string

const (
	ctxKeyLogger contextKey = "logger"
)

// CtxWithLogger injects a logger into the context
func CtxWithLogger(parentCtx context.Context, logger *Logger) context.Context {
	return context.WithValue(parentCtx, ctxKeyLogger, logger)
}

// LoggerFromCtx returns the logger injected to context. If there is no logger, return the default logger
func LoggerFromCtx(ctx context.Context) *Logger {
	res := ctx.Value(ctxKeyLogger)
	if res == nil {
		return getDefaultLog()
	}
	if logger, ok := res.(*Logger); ok {
		return logger
	}
	return getDefaultLog()
}
