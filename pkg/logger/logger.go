package logger

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

var (
	hostname string
	once     sync.Once
	log      *Log
)

type Log struct {
	logger zerolog.Logger
}

const (
	logLevelEnvName      string = "LOG_LEVEL"          // Default: DEBUG
	enableConsoleEnvName string = "ENABLE_CONSOLE_LOG" // Default: true

	logLevelTraceStr string = "TRACE"
	logLevelDebugStr string = "DEBUG"
	logLevelInfoStr  string = "INFO"
	logLevelWarnStr  string = "WARN"
	logLevelErrorStr string = "ERROR"

	logLevelTrace int = -1
	logLevelDebug int = 0
	logLevelInfo  int = 1
	logLevelWarn  int = 2
	logLevelError int = 3

	milliseconds float32 = 1000000.0
	logSize      int     = 1000
	pollInterval int64   = 10
)

func NewLog() *Log {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	logLevel, logLevelStr := getLogLevel()
	enableConsoleLogs, _ := strconv.ParseBool(getEnvVar(enableConsoleEnvName, "true"))
	var outputs []io.Writer
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	if enableConsoleLogs {
		wr := diode.NewWriter(os.Stdout, logSize, time.Duration(pollInterval)*time.Millisecond,
			func(missed int) {
				fmt.Printf("Logger dropped %d messages", missed)
			})
		outputs = append(outputs, wr)
	}

	multi := zerolog.MultiLevelWriter(outputs...)
	Log := &Log{logger: zerolog.New(multi).With().Timestamp().Logger()}
	Log.Debugf("Log level set to %s", logLevelStr)
	return Log
}

func Get() *Log {
	once.Do(func() {
		log = NewLog()
	})

	return log
}

func getLogLevel() (logLevel int, logLevelString string) {
	logLevelStr := getEnvVar(logLevelEnvName, logLevelDebugStr)
	switch logLevelStr {
	case logLevelTraceStr:
		return logLevelTrace, logLevelTraceStr
	case logLevelInfoStr:
		return logLevelInfo, logLevelInfoStr
	case logLevelWarnStr:
		return logLevelWarn, logLevelWarnStr
	case logLevelErrorStr:
		return logLevelError, logLevelErrorStr
	default:
		return logLevelDebug, logLevelDebugStr
	}
}

func getEnvVar(env, defaultValue string) string {
	value, found := os.LookupEnv(env)
	if !found || value == "" {
		return defaultValue
	}
	return value
}

func (l *Log) Infof(format string, args ...any) {
	l.logger.Info().Msgf(format, args...)
}

func (l *Log) Request(c echo.Context, start time.Time) {
	if c.Get("response-error") == nil {
		if c.Request().RequestURI != "/metrics" {
			if c.Get("response-body") == nil {
				l.logger.Info().
					Str("method", c.Request().Method).
					Int("status", c.Response().Status).
					Str("request", c.Request().RequestURI).
					Float32("timestamp", timestamp(start)).
					Msg("")
			} else {
				l.logger.Info().
					Str("method", c.Request().Method).
					Int("status", c.Response().Status).
					Str("request", c.Request().RequestURI).
					Interface("responseBody", c.Get("response-body")).
					Float32("timestamp", timestamp(start)).
					Msg("")
			}
		}
	} else {
		l.logger.Error().
			Str("method", c.Request().Method).
			Int("status", c.Response().Status).
			Str("request", c.Request().RequestURI).
			Interface("responseBody", c.Get("response-body")).
			Str("error", c.Get("response-error").(string)).
			Float32("timestamp", timestamp(start)).
			Msg("")
	}
}

func (l *Log) Debugf(format string, args ...any) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *Log) FatalIfError(moduleName, functionName string, errs ...error) {
	var sb strings.Builder
	for _, err := range errs {
		if err != nil {
			sb.WriteString(err.Error())
			sb.WriteString("\t")
		}
	}
	if sb.Len() > 0 {
		l.Fatalf(moduleName, functionName, sb.String())
	}
}

func (l *Log) Fatalf(moduleName, functionName, format string, args ...any) {
	l.logger.Fatal().
		Str("module", moduleName).
		Str("function", functionName).
		Msgf(format, args...)
}

func timestamp(start time.Time) float32 {
	return float32(time.Since(start).Nanoseconds()) / milliseconds
}
