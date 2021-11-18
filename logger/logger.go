package logger

import (
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"golang.org/x/term"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

const (
	dirPermMode  = 0744 // rwxr--r--
	filePermMode = 0644 // rw-r--r--

	consoleTimeFormat = time.RFC3339
)

func init() {
	zerolog.TimestampFunc = utcNow
}

func utcNow() time.Time {
	return time.Now().UTC()
}

type resilientMultiWriter struct {
	writers []io.Writer
}

// This custom resilientMultiWriter is an alternative to zerolog's so that we can make it resilient to individual
// writer's errors. E.g., when running as a Windows service, the console writer fails, but we don't want to
// allow that to prevent all logging to fail due to breaking the for loop upon an error.
func (t resilientMultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		_, _ = w.Write(p)
	}
	return len(p), nil
}

var levelErrorLogged = false

func newZeroLog(cfg config.LogConfig) *zerolog.Logger {
	var writers []io.Writer

	if cfg.ConsoleLoggingEnabled {
		writers = append(writers, createConsoleLogger(cfg))
	}
	if cfg.FileLoggingEnabled {
		rollingLogger, err := createRollingLogger(cfg)
		if err != nil {

		}
		fileLogger, err := createFileWriter(cfg)
		if err != nil {

		}

		writers = append(writers, rollingLogger)
		writers = append(writers, fileLogger)
	}

	multi := resilientMultiWriter{writers}

	level, levelErr := zerolog.ParseLevel(cfg.Level)
	if levelErr != nil {
		level = zerolog.InfoLevel
	}
	log := zerolog.New(multi).With().Timestamp().Logger().Level(level)
	if !levelErrorLogged && levelErr != nil {
		log.Error().Msgf("Failed to parse log level %q, using %q instead", cfg.Level, level)
		levelErrorLogged = true
	}

	return &log
}

func NewLogger(cfg config.LogConfig) *zerolog.Logger {
	return newZeroLog(cfg)
}

func createConsoleLogger(_ config.LogConfig) io.Writer {
	consoleOut := os.Stderr
	return zerolog.ConsoleWriter{
		Out:        colorable.NewColorable(consoleOut),
		NoColor:    !term.IsTerminal(int(consoleOut.Fd())), //TODO: add noColor to config
		TimeFormat: consoleTimeFormat,
	}
}

type fileInitializer struct {
	once          sync.Once
	writer        io.Writer
	creationError error
}

var (
	singleFileInit   fileInitializer
	rotatingFileInit fileInitializer
)

func createFileWriter(cfg config.LogConfig) (io.Writer, error) {
	singleFileInit.once.Do(func() {

		var logFile io.Writer
		fullPath := filepath.Join(cfg.Directory, cfg.Filename)

		// Try to open the existing file
		logFile, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, filePermMode)
		if err != nil {
			// If the existing file wasn't found, or couldn't be opened, just ignore
			// it and recreate a new one.
			logFile, err = createDirFile(cfg)
			// If creating a new logfile fails, then we have no choice but to error out.
			if err != nil {
				singleFileInit.creationError = err
				return
			}
		}

		singleFileInit.writer = logFile
	})

	return singleFileInit.writer, singleFileInit.creationError
}

func createDirFile(cfg config.LogConfig) (io.Writer, error) {
	if cfg.Directory != "" {
		err := os.MkdirAll(cfg.Directory, dirPermMode)

		if err != nil {
			return nil, fmt.Errorf("unable to create directories for new logfile: %s", err)
		}
	}

	mode := os.FileMode(filePermMode)

	fullPath := filepath.Join(cfg.Directory, cfg.Filename)
	logFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, mode)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new logfile: %s", err)
	}

	return logFile, nil
}

func createRollingLogger(cfg config.LogConfig) (io.Writer, error) {
	rotatingFileInit.once.Do(func() {
		if err := os.MkdirAll(cfg.Directory, dirPermMode); err != nil {
			rotatingFileInit.creationError = err
			return
		}

		rotatingFileInit.writer = &lumberjack.Logger{
			Filename:   path.Join(cfg.Directory, cfg.Filename),
			MaxBackups: cfg.MaxBackups,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
		}
	})

	return rotatingFileInit.writer, rotatingFileInit.creationError
}
