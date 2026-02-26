package logger

import (
	"os"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	tmpFile := "/tmp/test_logger.log"
	defer os.Remove(tmpFile)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid path", tmpFile, false},
		{"invalid path", "/nonexistent/path/to/file.log", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLog(t *testing.T) {
	tmpFile := "/tmp/test_log_msg.log"
	defer os.Remove(tmpFile)

	Init(tmpFile)

	tests := []struct {
		name string
		msg  string
	}{
		{"simple message", "test message"},
		{"empty message", ""},
		{"with spaces", "  spaced message  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.msg)
		})
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected log file to have content")
	}
}

func TestLogStart(t *testing.T) {
	tmpFile := "/tmp/test_log_start.log"
	defer os.Remove(tmpFile)

	Init(tmpFile)

	LogStart("test_process", 12345)

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "test_process") {
		t.Error("expected log to contain process name")
	}
	if !strings.Contains(string(content), "12345") {
		t.Error("expected log to contain pid")
	}
	if !strings.Contains(string(content), "started") {
		t.Error("expected log to contain 'started'")
	}
}

func TestLogStop(t *testing.T) {
	tmpFile := "/tmp/test_log_stop.log"
	defer os.Remove(tmpFile)

	Init(tmpFile)

	LogStop("my_service")

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "my_service") {
		t.Error("expected log to contain service name")
	}
	if !strings.Contains(string(content), "stopped") {
		t.Error("expected log to contain 'stopped'")
	}
}

func TestLogDied(t *testing.T) {
	tmpFile := "/tmp/test_log_died.log"
	defer os.Remove(tmpFile)

	Init(tmpFile)

	LogDied("worker", 1, "SIGTERM")

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logStr := string(content)
	if !strings.Contains(logStr, "worker") {
		t.Error("expected log to contain process name")
	}
	if !strings.Contains(logStr, "died unexpectedly") {
		t.Error("expected log to contain 'died unexpectedly'")
	}
	if !strings.Contains(logStr, "exit code 1") {
		t.Error("expected log to contain exit code")
	}
	if !strings.Contains(logStr, "SIGTERM") {
		t.Error("expected log to contain signal")
	}
}

func TestLogRestart(t *testing.T) {
	tmpFile := "/tmp/test_log_restart.log"
	defer os.Remove(tmpFile)

	Init(tmpFile)

	LogRestart("app")

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "app") {
		t.Error("expected log to contain app name")
	}
	if !strings.Contains(string(content), "restarting") {
		t.Error("expected log to contain 'restarting'")
	}
}

func TestLogReload(t *testing.T) {
	tmpFile := "/tmp/test_log_reload.log"
	defer os.Remove(tmpFile)

	Init(tmpFile)

	LogReload()

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "config reloaded") {
		t.Error("expected log to contain 'config reloaded'")
	}
}
