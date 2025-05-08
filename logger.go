package joylogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// Logger is a wrapper around zap.Logger to provide additional functionality.
type Logger struct {
	zap     *zap.Logger
	sugared *zap.SugaredLogger
}

// New creates a new Logger instance with configurable settings.
func New(prod bool, level zapcore.Level, logToFile bool, filePath string) (*Logger, error) {
	var cfg zap.Config
	if prod {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// Set the log level
	cfg.Level.SetLevel(level)

	// Configure log output destinations
	var writer zapcore.WriteSyncer
	if logToFile {
		// Ensure the directory exists for the log file
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}

		// Open or create the log file
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(file)
	} else {
		writer = zapcore.AddSync(os.Stdout)
	}

	// Set encoder configuration (we use production encoder here)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 format for timestamps
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create core and logger instance using zapcore
	core := zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(level))

	// Create a logger with the desired configuration
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		zap:     logger,
		sugared: logger.Sugar(),
	}, nil
}

// With adds context fields to the logger for structured logging.
func (l *Logger) With(fields ...zap.Field) *Logger {
	newZap := l.zap.With(fields...)
	return &Logger{
		zap:     newZap,
		sugared: newZap.Sugar(),
	}
}

// Info logs an info-level message.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

// Debug logs a debug-level message.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

// Error logs an error-level message.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

// Warn logs a warning-level message.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

// Fatal logs a fatal-level message and terminates the application.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// Sync ensures that any buffered log entries are written.
func (l *Logger) Sync() error {
	return l.zap.Sync()
}
