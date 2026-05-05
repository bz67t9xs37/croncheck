package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertType represents the kind of alert being sent.
type AlertType string

const (
	AlertMissed AlertType = "missed"
	AlertFailed AlertType = "failed"
)

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	JobName   string    `json:"job_name"`
	Alert     AlertType `json:"alert"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Notifier sends alert payloads to a configured webhook URL.
type Notifier struct {
	URL        string
	HTTPClient *http.Client
}

// NewNotifier creates a Notifier with a sensible default HTTP client.
func NewNotifier(url string) *Notifier {
	return &Notifier{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send marshals the payload and POSTs it to the webhook URL.
func (n *Notifier) Send(p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := n.HTTPClient.Post(n.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
