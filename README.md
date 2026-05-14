# logslice

A command-line log parser that extracts and filters structured log entries by time range or pattern.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

## Usage

```bash
# Filter logs by time range
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log

# Filter logs by pattern
logslice --pattern "ERROR" app.log

# Combine time range and pattern
logslice --from "2024-01-15T08:00:00" --pattern "timeout" app.log

# Read from stdin
cat app.log | logslice --pattern "WARN"
```

### Flags

| Flag        | Description                          |
|-------------|--------------------------------------|
| `--from`    | Start of time range (RFC3339 format) |
| `--to`      | End of time range (RFC3339 format)   |
| `--pattern` | Regex pattern to match log entries   |
| `--format`  | Log format: `json`, `logfmt`, `text` (default: `json`) |

## Example Output

```
2024-01-15T08:23:11Z [ERROR] connection timeout: host=db.internal retries=3
2024-01-15T08:45:02Z [ERROR] failed to write record: disk quota exceeded
```

## License

MIT — see [LICENSE](LICENSE) for details.