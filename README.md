# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected listeners with configurable rules.

---

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file:

```bash
portwatch --config portwatch.yaml
```

Example `portwatch.yaml`:

```yaml
interval: 30s
allowed_ports:
  - 22
  - 80
  - 443
alert:
  type: log
  path: /var/log/portwatch.log
```

portwatch will poll open ports at the configured interval and log (or notify) whenever an unexpected listener is detected.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `portwatch.yaml` | Path to config file |
| `--interval` | `30s` | Polling interval |
| `--verbose` | `false` | Enable verbose output |

---

## How It Works

portwatch reads the list of currently open TCP/UDP ports at each interval, compares them against your allowlist, and fires an alert for any port not explicitly permitted.

---

## License

MIT © yourname