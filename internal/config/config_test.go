package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "croncheck-config-*.json")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `{
		"listen_addr": ":9090",
		"webhook_url": "https://example.com/hook",
		"jobs": [
			{"name": "backup", "schedule": "0 2 * * *", "grace_period": "10m"}
		]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddr != ":9090" {
		t.Errorf("expected listen_addr :9090, got %s", cfg.ListenAddr)
	}
	if len(cfg.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].GracePeriod != 10*time.Minute {
		t.Errorf("expected grace_period 10m, got %v", cfg.Jobs[0].GracePeriod)
	}
}

func TestLoad_DefaultListenAddr(t *testing.T) {
	path := writeTempConfig(t, `{"jobs": []}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddr != ":8080" {
		t.Errorf("expected default :8080, got %s", cfg.ListenAddr)
	}
}

func TestLoad_InvalidGracePeriod(t *testing.T) {
	path := writeTempConfig(t, `{
		"jobs": [{"name": "x", "schedule": "* * * * *", "grace_period": "notaduration"}]
	}`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid grace_period, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
