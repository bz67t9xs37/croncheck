# croncheck

Lightweight daemon that monitors cron job execution and alerts on missed or failed runs via webhook.

## Installation

```bash
go install github.com/yourusername/croncheck@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/croncheck.git && cd croncheck && go build -o croncheck .
```

## Usage

Start the daemon with a config file:

```bash
croncheck --config /etc/croncheck/config.yaml
```

Example `config.yaml`:

```yaml
webhook_url: "https://hooks.slack.com/services/your/webhook/url"
jobs:
  - name: "daily-backup"
    schedule: "0 2 * * *"
    timeout: 30m
    grace_period: 5m
  - name: "hourly-sync"
    schedule: "0 * * * *"
    timeout: 5m
```

Wrap your cron job to report execution status:

```bash
# In your crontab
0 2 * * * croncheck exec --job daily-backup /usr/local/bin/backup.sh
```

croncheck will alert via webhook if a job:
- Does not start within the expected schedule window
- Exceeds its configured timeout
- Exits with a non-zero status code

## Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `webhook_url` | Webhook endpoint for alerts | required |
| `schedule` | Cron expression for expected run time | required |
| `timeout` | Max allowed execution duration | `1h` |
| `grace_period` | Extra time before a missed run triggers an alert | `5m` |

## License

MIT © 2024 croncheck contributors