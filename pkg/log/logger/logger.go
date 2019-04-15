package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/tokopedia/tdk/x/go/errors"
)

// Level type
type Level int

// level of log
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

// Logger struct
type Logger struct {
	logger zerolog.Logger
	config *Config
}

// Config of logger
type Config struct {
	Level   Level
	LogFile string
	Caller  bool
}

// New logger
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = &Config{
			Level:  InfoLevel,
			Caller: false,
		}
	}

	logger, err := newLogger(config)
	if err != nil {
		return nil, err
	}
	l := Logger{
		logger: logger,
		config: config,
	}
	return &l, nil
}

// DefaultLogger return default value of logger
func DefaultLogger() *Logger {
	defaultConfig := &Config{Level: InfoLevel, LogFile: ""}
	// do not check the error as error won't happen
	// error is only for checking file output
	l := Logger{
		logger: defaultLogger(defaultConfig.Level),
		config: defaultConfig,
	}
	return &l
}

func defaultLogger(level Level) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	logger = setLevel(logger, level)

	return logger
}

func newLogger(config *Config) (zerolog.Logger, error) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerSkipFrameCount = 4
	writers := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr})

	// set writer to file if config.LogFile is not empty
	if config.LogFile != "" {
		err := os.MkdirAll(filepath.Dir(config.LogFile), 0750)
		if err != nil {
			return zerolog.Logger{}, err
		}
		file, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0655)
		if err != nil {
			return zerolog.Logger{}, err
		}
		writers = zerolog.MultiLevelWriter(writers, file)
	}

	logger := zerolog.New(writers)
	logger = setLevel(logger, config.Level)
	if config.Caller {
		logger = logger.With().Caller().Logger()
	}

	return logger, nil
}

func setLevel(logger zerolog.Logger, level Level) zerolog.Logger {
	switch level {
	case DebugLevel:
		logger = logger.Level(zerolog.DebugLevel)
	case InfoLevel:
		logger = logger.Level(zerolog.InfoLevel)
	case WarnLevel:
		logger = logger.Level(zerolog.WarnLevel)
	case ErrorLevel:
		logger = logger.Level(zerolog.ErrorLevel)
	case FatalLevel:
		logger = logger.Level(zerolog.FatalLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}
	return logger
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level Level) {
	l.logger = setLevel(l.logger, level)
}

// SetLevelString set level using string instead of level
func (l *Logger) SetLevelString(level string) {
	l.logger = setLevel(l.logger, StringToLevel(level))
}

// StringToLevel convert string to log level
func StringToLevel(s string) Level {
	switch strings.ToLower(s) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		// TODO: make this more informative when happened
		return InfoLevel
	}
}

// LevelToString convert log level to readable string
func LevelToString(l Level) string {
	switch l {
	case DebugLevel:
		return DebugLevelString
	case InfoLevel:
		return InfoLevelString
	case WarnLevel:
		return WarnLevelString
	case ErrorLevel:
		return ErrorLevelString
	case FatalLevel:
		return FatalLevelString
	default:
		return InfoLevelString
	}
}

// DebugEvent function
func (l *Logger) DebugEvent() *zerolog.Event {
	return l.logger.Debug().Timestamp()
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug().Timestamp().Msg(fmt.Sprint(args...))
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Timestamp().Msgf(format, v...)
}

// InfoEvent function
func (l *Logger) InfoEvent() *zerolog.Event {
	return l.logger.Info().Timestamp()
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info().Timestamp().Msg(fmt.Sprint(args...))
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Timestamp().Msgf(format, v...)
}

// WarnEvent function
func (l *Logger) WarnEvent() *zerolog.Event {
	return l.logger.Warn().Timestamp()
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn().Timestamp().Msg(fmt.Sprint(args...))
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Timestamp().Msgf(format, v...)
}

// ErrorEvent function
func (l *Logger) ErrorEvent() *zerolog.Event {
	return l.logger.Error().Timestamp()
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error().Timestamp().Msg(fmt.Sprint(args...))
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Timestamp().Msgf(format, v...)
}

// Errors function to log errors package
func (l *Logger) Errors(err error) {
	switch err.(type) {
	case *errors.Error:
		e := err.(*errors.Error)

		fields := e.GetFields()
		if fields == nil {
			fields = make(errors.Fields)
		}
		fields["operations"] = e.OpTraces

		l.logger.Error().Timestamp().Fields(fields).Msg(e.Error())
	case error:
		l.logger.Error().Timestamp().Msg(err.Error())
	}
}

// FatalEvent function
func (l *Logger) FatalEvent() *zerolog.Event {
	return l.logger.Fatal().Timestamp()
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal().Timestamp().Msg(fmt.Sprint(args...))
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Timestamp().Msgf(format, v...)
}

// With function for logging with logging context
func With(event *zerolog.Event, msg string, keyValues map[string]interface{}) {
	for k, v := range keyValues {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}
