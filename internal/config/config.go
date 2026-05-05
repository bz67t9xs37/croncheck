package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// JobConfig holds the configuration for a single monitored cron job.
type JobConfig struct {
	Name         string        `json:"name"`
	Schedule     string        `json:"schedule"`
	GracePeriod  time.Duration `json:"grace_period"`
	WebhookURL   string        `json:"webhook_url"`
}

// Config holds the top-level application configuration.
type Config struct {
	ListenAddr  string        `json:"listen_addr"`
	WebhookURL  string        `json:"webhook_url"`
	Jobs        []JobConfig   `json:"jobs"`
}

// UnmarshalJSON implements custom unmarshalling to support duration strings.
func (j *JobConfig) UnmarshalJSON(data []byte) error {
	type Alias struct {
		Name        string `json:"name"`
		Schedule    string `json:"schedule"`
		GracePeriod string `json:"grace_period"`
		WebhookURL  string `json:"webhook_url"`
	}
	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	j.Name = a.Name
	j.Schedule = a.Schedule
	j.WebhookURL = a.WebhookURL
	if a.GracePeriod != "" {
		d, err := time.ParseDuration(a.GracePeriod)
		if err != nil {
			return fmt.Errorf("invalid grace_period %q: %w", a.GracePeriod, err)
		}
		j.GracePeriod = d
	}
	return nil
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":8080"
	}
	return &cfg, nil
}
