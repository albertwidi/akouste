package log

import (
	"github.com/tokopedia/tdk/go/log/logger"
)

// KV is a type for logging with more information
// this used by with function
type KV map[string]interface{}

// level of log
const (
	DebugLevel = logger.DebugLevel
	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	FatalLevel = logger.FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

var defaultLogger *logger.Logger
var debugLogger *logger.Logger

func init() {
	defaultLogger = logger.DefaultLogger()
	debugLogger, _ = logger.New(&logger.Config{Level: logger.DebugLevel})
}

// Config of log
type Config struct {
	Level string
	// LogFile for log to file
	// this is not needed by default
	// application is expected to run in containerized environment
	LogFile   string
	DebugFile string
	// set true to log line numbers
	// make sure you understand the overhead when use this
	Caller bool
}

// SetConfig to the current logger
func SetConfig(config *Config) error {
	var err error

	loggerConfig := logger.Config{
		Level: logger.InfoLevel,
	}
	debugLoggerConfig := logger.Config{
		Level: logger.DebugLevel,
	}

	if config != nil {
		// level
		loggerConfig.Level = logger.StringToLevel(config.Level)
		debugLoggerConfig.Level = logger.StringToLevel(config.Level)
		// output file
		loggerConfig.LogFile = config.LogFile
		debugLoggerConfig.LogFile = config.DebugFile
		// runtime caller
		loggerConfig.Caller = config.Caller
		debugLoggerConfig.Caller = config.Caller
	}

	newLogger, err := logger.New(&loggerConfig)
	if err != nil {
		return err
	}
	// extra check because it is very difficult to debug if the log itself causes the panic
	if newLogger != nil {
		defaultLogger = newLogger
	}

	newDebugLogger, err := logger.New(&debugLoggerConfig)
	if err != nil {
		return err
	}
	if newDebugLogger != nil {
		debugLogger = newDebugLogger
	}

	return nil
}

// SetLevel of log
func SetLevel(level logger.Level) {
	setLevel(level)
}

// SetLevelString to set log level using string
func SetLevelString(level string) {
	setLevel(logger.StringToLevel(level))
}

// setLevel function set the log level to the desired level for defaultLogger and debugLogger
// debugLogger level can go to any level, but not with defaultLogger
// this to make sure debugLogger to be disabled when level is > debug
// and defaultLogger to not overlap with debugLogger
func setLevel(level logger.Level) {
	if level < InfoLevel {
		debugLogger.SetLevel(level)
	} else {
		defaultLogger.SetLevel(level)
		debugLogger.SetLevel(level)
	}
}

// Debug function
func Debug(args ...interface{}) {
	debugLogger.Debug(args...)
}

// Debugf function
func Debugf(format string, v ...interface{}) {
	debugLogger.Debugf(format, v...)
}

// Debugw function
func Debugw(msg string, keyValues KV) {
	logger.With(debugLogger.DebugEvent(), msg, keyValues)
}

// Print function
func Print(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Println function
func Println(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Printf function
func Printf(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Info function
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof function
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Infow function
func Infow(msg string, keyValues KV) {
	logger.With(defaultLogger.InfoEvent(), msg, keyValues)
}

// Warn function
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf function
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Warnw function
func Warnw(msg string, keyValues KV) {
	logger.With(defaultLogger.WarnEvent(), msg, keyValues)
}

// Error function
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf function
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Errorw function
func Errorw(msg string, keyValues KV) {
	logger.With(defaultLogger.ErrorEvent(), msg, keyValues)
}

// Errors function to log errors package
func Errors(err error) {
	defaultLogger.Errors(err)
}

// Fatal function
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf function
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

// Fatalw function
func Fatalw(msg string, keyValues KV) {
	logger.With(defaultLogger.FatalEvent(), msg, keyValues)
}
