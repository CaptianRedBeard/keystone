package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerWritesToFile(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "keystone_test.log")
	defer os.Remove(tmp)

	if err := Init(tmp, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	Info("hello world", false)

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("failed to read logfile: %v", err)
	}

	log := string(data)
	if !strings.Contains(log, "hello world") {
		t.Fatalf("log entry missing message: %s", log)
	}
	if !strings.Contains(log, "✅") {
		t.Fatalf("log entry missing info symbol: %s", log)
	}
}

func TestVerboseWritesToProvidedWriter(t *testing.T) {
	var buf bytes.Buffer
	tmp := filepath.Join(os.TempDir(), "keystone_verbose.log")
	defer os.Remove(tmp)

	if err := InitWithWriter(tmp, true, &buf); err != nil {
		t.Fatalf("InitWithWriter failed: %v", err)
	}
	defer Close()

	Warn("something happened", false)

	output := buf.String()
	if !strings.Contains(output, "something happened") {
		t.Fatalf("expected verbose output, got: %s", output)
	}
	if !strings.Contains(output, "⚠️") {
		t.Fatalf("expected warning symbol, got: %s", output)
	}
}

func TestErrorSymbol(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "keystone_error.log")
	defer os.Remove(tmp)

	if err := Init(tmp, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	Error("boom", false)

	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "❌") {
		t.Fatal("missing error symbol in output")
	}
}

func TestCloseDoesNotPanic(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "keystone_close.log")
	defer os.Remove(tmp)

	if err := Init(tmp, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Should not panic even if called twice
	Close()
	Close()
}
