package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to read recent log file contents
func readLogFile(agentID string) (string, error) {
	data, err := os.ReadFile(filepath.Join("logs", agentID+".log"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestInitAndGet(t *testing.T) {
	Init(false)

	agentID := "test_agent"
	logger := Get(agentID)
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	// Confirm log file is created
	logPath := filepath.Join("logs", agentID+".log")
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("expected log file %s to exist: %v", logPath, err)
	}
}

func TestLogWritesToFile(t *testing.T) {
	Init(false)

	agentID := "file_test_agent"
	Log(agentID, "hello file logger")

	content, err := readLogFile(agentID)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(content, "hello file logger") {
		t.Errorf("expected log content to contain message, got: %s", content)
	}
}

func TestLogVerboseOutputsToConsole(t *testing.T) {
	Init(true)

	// Capture stdout
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Log("verbose_agent", "this should appear in console")

	w.Close()
	os.Stdout = orig
	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "this should appear in console") {
		t.Errorf("expected console output when verbose, got: %s", string(out))
	}
}

func TestLoggerReuse(t *testing.T) {
	Init(false)

	a1 := Get("reuse_agent")
	a2 := Get("reuse_agent")

	if a1 != a2 {
		t.Error("expected same logger instance for same agentID")
	}
}
