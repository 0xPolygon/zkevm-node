package log

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/hermeznetwork/tracerr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogEnvironment represents the possible log environments.
type LogEnvironment string

const (
	// EnvironmentProduction production log environment.
	EnvironmentProduction = LogEnvironment("production")
	// EnvironmentDevelopment development log environment.
	EnvironmentDevelopment = LogEnvironment("development")
)

// Logger is a wrapper providing logging facilities.
type Logger struct {
	x *zap.SugaredLogger
}

// root logger
var log atomic.Pointer[Logger]

func getDefaultLog() *Logger {
	l := log.Load()
	if l != nil {
		return l
	}
	// default level: debug
	zapLogger, _, err := NewLogger(Config{
		Environment: EnvironmentDevelopment,
		Level:       "debug",
		Outputs:     []string{"stderr"},
	})
	if err != nil {
		panic(err)
	}
	log.Store(&Logger{x: zapLogger})
	return log.Load()
}

// GetDefaultLog returns the default logger
func GetDefaultLog() *Logger {
	return getDefaultLog()
}

// Init the logger with defined level. outputs defines the outputs where the
// logs will be sent. By default outputs contains "stdout", which prints the
// logs at the output of the process. To add a log file as output, the path
// should be added at the outputs array. To avoid printing the logs but storing
// them on a file, can use []string{"pathtofile.log"}
func Init(cfg Config) {
	zapLogger, _, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	log.Store(&Logger{x: zapLogger})
}

// NewLogger creates the logger with defined level. outputs defines the outputs where the
// logs will be sent. By default, outputs contains "stdout", which prints the
// logs at the output of the process. To add a log file as output, the path
// should be added at the outputs array. To avoid printing the logs but storing
// them on a file, can use []string{"pathtofile.log"}
func NewLogger(cfg Config) (*zap.SugaredLogger, *zap.AtomicLevel, error) {
	var level zap.AtomicLevel
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, nil, fmt.Errorf("error on setting log level: %s", err)
	}

	var zapCfg zap.Config

	switch cfg.Environment {
	case EnvironmentProduction:
		zapCfg = zap.NewProductionConfig()
	default:
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	zapCfg.Level = level
	zapCfg.OutputPaths = cfg.Outputs
	zapCfg.InitialFields = map[string]interface{}{
		"version": zkevm.Version,
		"pid":     os.Getpid(),
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, nil, err
	}
	defer logger.Sync() //nolint:gosec,errcheck

	// skip 2 callers: one for our wrapper methods and one for the package functions
	withOptions := logger.WithOptions(zap.AddCallerSkip(2)) //nolint:gomnd
	return withOptions.Sugar(), &level, nil
}

// WithFields returns a new Logger (derived from the root one) with additional
// fields as per keyValuePairs.  The root Logger instance is not affected.
func WithFields(keyValuePairs ...interface{}) *Logger {
	l := getDefaultLog().WithFields(keyValuePairs...)

	// since we are returning a new instance, remove one caller from the
	// stack, because we'll be calling the retruned Logger methods
	// directly, not the package functions.
	x := l.x.WithOptions(zap.AddCallerSkip(-1))
	l.x = x
	return l
}

// WithFields returns a new Logger with additional fields as per keyValuePairs.
// The original Logger instance is not affected.
func (l *Logger) WithFields(keyValuePairs ...interface{}) *Logger {
	return &Logger{
		x: l.x.With(keyValuePairs...),
	}
}

func sprintStackTrace(st []tracerr.Frame) string {
	builder := strings.Builder{}
	// Skip deepest frame because it belongs to the go runtime and we don't
	// care about it.
	if len(st) > 0 {
		st = st[:len(st)-1]
	}
	for _, f := range st {
		builder.WriteString(fmt.Sprintf("\n%s:%d %s()", f.Path, f.Line, f.Func))
	}
	builder.WriteString("\n")
	return builder.String()
}

// appendStackTraceMaybeArgs will append the stacktrace to the args
func appendStackTraceMaybeArgs(args []interface{}) []interface{} {
	for i := range args {
		if err, ok := args[i].(error); ok {
			err = tracerr.Wrap(err)
			st := tracerr.StackTrace(err)
			return append(args, sprintStackTrace(st))
		}
	}
	return args
}

// Debug calls log.Debug
func (l *Logger) Debug(args ...interface{}) {
	l.x.Debug(args...)
}

// Info calls log.Info
func (l *Logger) Info(args ...interface{}) {
	l.x.Info(args...)
}

// Warn calls log.Warn
func (l *Logger) Warn(args ...interface{}) {
	l.x.Warn(args...)
}

// Error calls log.Error
func (l *Logger) Error(args ...interface{}) {
	l.x.Error(args...)
}

// Fatal calls log.Fatal
func (l *Logger) Fatal(args ...interface{}) {
	l.x.Fatal(args...)
}

// Debugf calls log.Debugf
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.x.Debugf(template, args...)
}

// Infof calls log.Infof
func (l *Logger) Infof(template string, args ...interface{}) {
	l.x.Infof(template, args...)
}

// Warnf calls log.Warnf
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.x.Warnf(template, args...)
}

// Fatalf calls log.Fatalf
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.x.Fatalf(template, args...)
}

// Errorf calls log.Errorf and stores the error message into the ErrorFile
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.x.Errorf(template, args...)
}

// Debug calls log.Debug on the root Logger.
func Debug(args ...interface{}) {
	getDefaultLog().Debug(args...)
}

// Info calls log.Info on the root Logger.
func Info(args ...interface{}) {
	getDefaultLog().Info(args...)
}

// Warn calls log.Warn on the root Logger.
func Warn(args ...interface{}) {
	getDefaultLog().Warn(args...)
}

// Error calls log.Error on the root Logger.
func Error(args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLog().Error(args...)
}

// Fatal calls log.Fatal on the root Logger.
func Fatal(args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLog().Fatal(args...)
}

// Debugf calls log.Debugf on the root Logger.
func Debugf(template string, args ...interface{}) {
	getDefaultLog().Debugf(template, args...)
}

// Infof calls log.Infof on the root Logger.
func Infof(template string, args ...interface{}) {
	getDefaultLog().Infof(template, args...)
}

// Warnf calls log.Warnf on the root Logger.
func Warnf(template string, args ...interface{}) {
	getDefaultLog().Warnf(template, args...)
}

// Fatalf calls log.Fatalf on the root Logger.
func Fatalf(template string, args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLog().Fatalf(template, args...)
}

// Errorf calls log.Errorf on the root logger and stores the error message into
// the ErrorFile.
func Errorf(template string, args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLog().Errorf(template, args...)
}

// appendStackTraceMaybeKV will append the stacktrace to the KV
func appendStackTraceMaybeKV(msg string, kv []interface{}) string {
	for i := range kv {
		if i%2 == 0 {
			continue
		}
		if err, ok := kv[i].(error); ok {
			err = tracerr.Wrap(err)
			st := tracerr.StackTrace(err)
			return fmt.Sprintf("%v: %v%v\n", msg, err, sprintStackTrace(st))
		}
	}
	return msg
}

// Debugw calls log.Debugw
func (l *Logger) Debugw(msg string, kv ...interface{}) {
	l.x.Debugw(msg, kv...)
}

// Infow calls log.Infow
func (l *Logger) Infow(msg string, kv ...interface{}) {
	l.x.Infow(msg, kv...)
}

// Warnw calls log.Warnw
func (l *Logger) Warnw(msg string, kv ...interface{}) {
	l.x.Warnw(msg, kv...)
}

// Errorw calls log.Errorw
func (l *Logger) Errorw(msg string, kv ...interface{}) {
	l.x.Errorw(msg, kv...)
}

// Fatalw calls log.Fatalw
func (l *Logger) Fatalw(msg string, kv ...interface{}) {
	l.x.Fatalw(msg, kv...)
}

// Debugw calls log.Debugw on the root Logger.
func Debugw(msg string, kv ...interface{}) {
	getDefaultLog().Debugw(msg, kv...)
}

// Infow calls log.Infow on the root Logger.
func Infow(msg string, kv ...interface{}) {
	getDefaultLog().Infow(msg, kv...)
}

// Warnw calls log.Warnw on the root Logger.
func Warnw(msg string, kv ...interface{}) {
	getDefaultLog().Warnw(msg, kv...)
}

// Errorw calls log.Errorw on the root Logger.
func Errorw(msg string, kv ...interface{}) {
	msg = appendStackTraceMaybeKV(msg, kv)
	getDefaultLog().Errorw(msg, kv...)
}

// Fatalw calls log.Fatalw on the root Logger.
func Fatalw(msg string, kv ...interface{}) {
	msg = appendStackTraceMaybeKV(msg, kv)
	getDefaultLog().Fatalw(msg, kv...)
}
