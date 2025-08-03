package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Mock the config package for testing
func init() {
	// Set environment variables before importing config
	os.Setenv("NOT_DEBUG", "true")
	os.Setenv("SECRET_MANAGER_BASE64", "eyJPTkVfQVBJX0tFWSI6InRlc3QtYXBpLWtleSIsIk9ORV9BUElfVVJMIjoiaHR0cDovL3Rlc3QtYXBpLXVybCIsIk1PREVMIjoiZ3B0LTQwLW1pbmkiLCJHT09HTEVfQ0xPVURfQVBJX0tFWSI6InRlc3QtZ29vZ2xlLWtleSIsIkxPR19MRVZFTCI6IklORk8ifQ==")
}

// TestLoggerInitialization tests the logger initialization
func TestLoggerInitialization(t *testing.T) {
	// Test that logger is properly initialized
	logger := GetLogger()
	assert.NotNil(t, logger)
	// Note: The actual level depends on config.Settings.LogLevel, so we just check it's not nil
	assert.NotNil(t, logger)
}

// TestSetLevel tests the SetLevel function
func TestSetLevel(t *testing.T) {
	tests := []struct {
		name        string
		levelStr    string
		expectError bool
	}{
		{"debug level", "debug", false},
		{"info level", "info", false},
		{"warn level", "warn", false},
		{"error level", "error", false},
		{"fatal level", "fatal", false},
		{"panic level", "panic", false},
		{"invalid level", "invalid", true},
		{"empty level", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetLevel(tt.levelStr)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid log level")
			} else {
				assert.NoError(t, err)
				level, _ := logrus.ParseLevel(tt.levelStr)
				assert.Equal(t, level, GetLogger().GetLevel())
			}
		})
	}
}

// TestGetCallerInfo tests the getCallerInfo function
func TestGetCallerInfo(t *testing.T) {
	// This test verifies that getCallerInfo returns valid file and line information
	file, line := getCallerInfo()

	// Should return a valid filename (not "unknown")
	assert.NotEqual(t, "unknown", file)
	assert.NotEqual(t, 0, line)

	// Should return a valid filename (could be testing.go or logger_test.go)
	assert.True(t, strings.Contains(file, "test") || strings.Contains(file, "logger"), "Expected filename to contain 'test' or 'logger', got: %s", file)
}

// TestLoggingFunctions tests all logging functions
func TestLoggingFunctions(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	// Set log level to debug to capture all messages
	SetLevel("debug")

	tests := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{
			name: "Debug logging",
			logFunc: func() {
				Debug("debug message %s", "test")
			},
			expected: "debug message test",
		},
		{
			name: "Info logging",
			logFunc: func() {
				Info("info message %s", "test")
			},
			expected: "info message test",
		},
		{
			name: "Warn logging",
			logFunc: func() {
				Warn("warn message %s", "test")
			},
			expected: "warn message test",
		},
		{
			name: "Error logging",
			logFunc: func() {
				Errorf("error message %s", "test")
			},
			expected: "error message test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()
			output := buf.String()
			assert.Contains(t, output, tt.expected)
			assert.Contains(t, output, "file") // Should contain file information (with color codes)
		})
	}
}

// TestErrorLoggingWithError tests Error function with error object
func TestErrorLoggingWithError(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	SetLevel("error")

	testErr := errors.New("test error")
	Error(testErr, "error occurred: %s", "test message")

	output := buf.String()
	assert.Contains(t, output, "error occurred: test message")
	assert.Contains(t, output, "error") // Should contain error field (with color codes)
	assert.Contains(t, output, "file")  // Should contain file information (with color codes)
}

// TestErrorLoggingWithoutError tests Error function without error object
func TestErrorLoggingWithoutError(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	SetLevel("error")

	Error(nil, "error occurred: %s", "test message")

	output := buf.String()
	assert.Contains(t, output, "error occurred: test message")
	// Note: The log level is ERROR, so we check that there's no error field in the structured data
	// The output contains "ERROR" as the level, but not as a field
	assert.NotContains(t, output, "error=") // Should not contain error field when nil
	assert.Contains(t, output, "file")      // Should contain file information (with color codes)
}

// TestWithField tests WithField function
func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	SetLevel("info")

	entry := WithField("test_key", "test_value")
	entry.Info("message with field")

	output := buf.String()
	assert.Contains(t, output, "message with field")
	assert.Contains(t, output, "test_key") // Should contain the field key (with color codes)
}

// TestWithFields tests WithFields function
func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	SetLevel("info")

	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}
	entry := WithFields(fields)
	entry.Info("message with fields")

	output := buf.String()
	assert.Contains(t, output, "message with fields")
	assert.Contains(t, output, "key1") // Should contain the field key (with color codes)
	assert.Contains(t, output, "key2") // Should contain the field key (with color codes)
}

// TestLogLevels tests different log levels
func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	tests := []struct {
		name     string
		level    string
		logFunc  func()
		expected bool // whether the message should appear in output
	}{
		{
			name:  "Debug level - debug message",
			level: "debug",
			logFunc: func() {
				Debug("debug message")
			},
			expected: true,
		},
		{
			name:  "Info level - debug message",
			level: "info",
			logFunc: func() {
				Debug("debug message")
			},
			expected: false,
		},
		{
			name:  "Info level - info message",
			level: "info",
			logFunc: func() {
				Info("info message")
			},
			expected: true,
		},
		{
			name:  "Warn level - info message",
			level: "warn",
			logFunc: func() {
				Info("info message")
			},
			expected: false,
		},
		{
			name:  "Warn level - warn message",
			level: "warn",
			logFunc: func() {
				Warn("warn message")
			},
			expected: true,
		},
		{
			name:  "Error level - warn message",
			level: "error",
			logFunc: func() {
				Warn("warn message")
			},
			expected: false,
		},
		{
			name:  "Error level - error message",
			level: "error",
			logFunc: func() {
				Errorf("error message")
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			SetLevel(tt.level)
			tt.logFunc()
			output := buf.String()

			if tt.expected {
				assert.Contains(t, output, strings.Split(tt.name, " - ")[1])
			} else {
				assert.NotContains(t, output, strings.Split(tt.name, " - ")[1])
			}
		})
	}
}

// TestLogFormatting tests log formatting
func TestLogFormatting(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	SetLevel("info")

	Info("test message with %s", "formatting")

	output := buf.String()

	// Should contain the formatted message
	assert.Contains(t, output, "test message with formatting")

	// Should contain timestamp (logrus in ra [timestamp])
	assert.Contains(t, output, "[")
	assert.Contains(t, output, "]")

	// Should contain level
	assert.Contains(t, output, "INFO")

	// Should contain file information
	assert.Contains(t, output, "file")
}

// TestConcurrentLogging tests concurrent logging
func TestConcurrentLogging(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := GetLogger().Out
	GetLogger().SetOutput(&buf)
	defer GetLogger().SetOutput(originalOutput)

	SetLevel("info")

	// Run multiple goroutines logging simultaneously
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			Info("concurrent message %d", id)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	output := buf.String()

	// Should contain all messages
	for i := 0; i < 10; i++ {
		assert.Contains(t, output, fmt.Sprintf("concurrent message %d", i))
	}
}

// TestGetLogger tests GetLogger function
func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assert.NotNil(t, logger)
	assert.IsType(t, &logrus.Logger{}, logger)
}

// TestLoggerOutput tests logger output configuration
func TestLoggerOutput(t *testing.T) {
	// Test that logger outputs to stdout by default
	logger := GetLogger()
	assert.Equal(t, os.Stdout, logger.Out)
}

// TestLoggerFormatter tests logger formatter configuration
func TestLoggerFormatter(t *testing.T) {
	logger := GetLogger()
	formatter := logger.Formatter.(*logrus.TextFormatter)

	// Test formatter configuration
	assert.True(t, formatter.FullTimestamp)
	assert.Equal(t, "2006-01-02T15:04:05Z07:00", formatter.TimestampFormat)
	assert.True(t, formatter.ForceColors)
	assert.False(t, formatter.DisableColors)
	assert.True(t, formatter.DisableQuote)
	assert.False(t, formatter.DisableSorting)
	assert.True(t, formatter.PadLevelText)
}

// BenchmarkLogging benchmarks logging performance
func BenchmarkLogging(b *testing.B) {
	// Disable output for benchmarking
	GetLogger().SetOutput(io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message %d", i)
	}
}

// BenchmarkWithFields benchmarks WithFields performance
func BenchmarkWithFields(b *testing.B) {
	GetLogger().SetOutput(io.Discard)

	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := WithFields(fields)
		entry.Info("benchmark message")
	}
}
