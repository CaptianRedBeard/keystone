package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	globalVerbose bool
	globalLogger  *log.Logger
	logFile       *os.File
	mu            sync.Mutex
)

const defaultLogPath = "./logs/keystone.log"

// InitDefault initializes logger using the default dev path
func InitDefault(verbose bool) error {
	const defaultPath = "logs/keystone.log"
	return Init(defaultPath, verbose)
}

// Init initializes the logger system with verbosity control.
func Init(logPath string, verbose bool) error {
	return InitWithWriter(logPath, verbose, nil)
}

// InitWithWriter allows using a custom writer (useful in tests)
func InitWithWriter(logPath string, verbose bool, out io.Writer) error {
	mu.Lock()
	defer mu.Unlock()

	globalVerbose = verbose

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	logFile = file

	var writer io.Writer
	if out != nil {
		writer = io.MultiWriter(file, out)
	} else {
		writer = file
	}
	globalLogger = log.New(writer, "", 0)
	return nil
}

// Close should be called on program exit (or in tests)
func Close() {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		_ = logFile.Close()
	}
}

// logEntry handles the actual log formatting and optional console output.
func logEntry(symbol, message string, jsonMode bool) {
	mu.Lock()
	defer mu.Unlock()

	if globalLogger == nil {
		globalLogger = log.New(os.Stderr, "", 0)
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] %s %s\n", ts, symbol, message)

	globalLogger.Print(entry)

	if globalVerbose && !jsonMode {
		fmt.Print(entry)
	}
}

func Info(message string, jsonMode bool)  { logEntry("✅", message, jsonMode) }
func Warn(message string, jsonMode bool)  { logEntry("⚠️", message, jsonMode) }
func Error(message string, jsonMode bool) { logEntry("❌", message, jsonMode) }
