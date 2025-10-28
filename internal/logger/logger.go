package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	globalVerbose bool
	loggers       = make(map[string]*log.Logger)
	mu            sync.Mutex
)

// Init sets the global verbosity (from --verbose)
func Init(verbose bool) {
	globalVerbose = verbose
	_ = os.MkdirAll("logs", 0755)
}

func Get(agentID string) *log.Logger {
	mu.Lock()
	defer mu.Unlock()

	if logger, ok := loggers[agentID]; ok {
		return logger
	}

	logPath := filepath.Join("logs", fmt.Sprintf("%s.log", agentID))
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file for agent %s: %v\n", agentID, err)
		file = os.Stdout
	}

	logger := log.New(file, "", 0)
	loggers[agentID] = logger
	return logger
}

// Log writes a message to both file and console (if verbose)
func Log(agentID, message string) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] %s\n", ts, message)
	Get(agentID).Print(entry)

	if globalVerbose {
		fmt.Print(entry)
	}
}
