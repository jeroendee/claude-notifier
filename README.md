# Claude Notifier

A macOS menu bar notification app for Claude Code. Receive notifications when Claude completes tasks, with sound alerts and notification history accessible from your menu bar.

## Features

- **Menu bar notifications** - Persistent icon in macOS menu bar shows notification history
- **Sound alerts** - Plays a sound when notifications arrive (configurable)
- **Notification history** - Stores up to 50 notifications with read/unread status
- **HTTP API** - Receive notifications via simple HTTP POST requests
- **Claude Code integration** - Hook script for automatic task completion notifications

## Quick Start

```bash
# Build
make build

# Run
./bin/claude-notifier

# Check version
./bin/claude-notifier --version
```

The claude-notifier icon appears in your macOS menu bar (top-right, near clock/wifi/battery). Click it to see notifications.

## Installation

### Install Binary

```bash
make install
```

Installs `claude-notifier` to `GOPATH/bin` (or `GOBIN` if set).

### Auto-Start on Login

```bash
make install-launchagent
```

Installs a LaunchAgent to start the claude-notifier automatically when you log in.

### Uninstall

```bash
make uninstall-launchagent
make uninstall
```

## Configuration

Environment variables (with defaults):

| Variable | Default | Description |
|----------|---------|-------------|
| `CLAUDE_NOTIFIER_PORT` | `19199` | HTTP server port |
| `CLAUDE_NOTIFIER_SOUND` | `/System/Library/Sounds/Glass.aiff` | Sound file to play |
| `CLAUDE_NOTIFIER_MAX_HISTORY` | `50` | Maximum notifications stored |

Example:
```bash
CLAUDE_NOTIFIER_PORT=8080 CLAUDE_NOTIFIER_SOUND=/System/Library/Sounds/Ping.aiff ./bin/claude-notifier
```

## Claude Code Integration

### Install Hook Script

```bash
make install-hook
```

This installs a hook script to `~/.claude/hooks/claude-hook.sh`.

### Configure Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "Stop": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "~/.claude/hooks/claude-hook.sh"
          }
        ]
      }
    ]
  }
}
```

Now Claude Code will notify you when tasks complete.

## HTTP API

### Send Notification

```bash
curl -X POST http://localhost:19199/notify \
  -H "Content-Type: application/json" \
  -d '{"message":"Task completed"}'
```

Returns `201 Created` on success.

### Health Check

```bash
curl http://localhost:19199/health
```

Returns `{"status":"ok"}`.

### Clear Notifications

```bash
curl -X POST http://localhost:19199/clear
```

Returns `204 No Content`.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Claude Code                               │
│                    (via Stop hook)                               │
└─────────────────────┬───────────────────────────────────────────┘
                      │ HTTP POST /notify
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      HTTP Server                                 │
│                   (internal/server)                              │
│                                                                  │
│  Endpoints:                                                      │
│  • POST /notify  → receive notification                          │
│  • GET  /health  → health check                                  │
│  • POST /clear   → clear all notifications                       │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                         App                                      │
│                    (internal/app)                                │
│                                                                  │
│  • Loads configuration from environment                          │
│  • Coordinates all components                                    │
│  • Handles notification flow                                     │
└───────┬─────────────────────────────────┬───────────────────────┘
        │                                 │
        ▼                                 ▼
┌───────────────────────┐    ┌────────────────────────────────────┐
│    Notification       │    │           UI (systray)             │
│  (internal/notification)│   │          (internal/ui)             │
│                       │    │                                    │
│  • Store (history)    │    │  • Menu bar icon                   │
│  • Sound player       │    │  • Notification list               │
│  • Read/unread state  │    │  • Mark read / Clear actions       │
└───────────────────────┘    └────────────────────────────────────┘
```

### Data Flow

1. **Notification arrives** via HTTP POST to `/notify`
2. **Server** passes notification to App
3. **App** stores notification and triggers sound
4. **UI** refreshes menu bar to show new notification
5. **User** clicks menu bar icon to view/manage notifications

### Key Components

| Component | Package | Responsibility |
|-----------|---------|----------------|
| Config | `internal/app` | Load environment variables, provide defaults |
| Server | `internal/server` | HTTP API endpoints |
| Store | `internal/notification` | Notification storage, read/unread tracking |
| Sound | `internal/notification` | Play audio alerts via `afplay` |
| UI | `internal/ui` | macOS menu bar integration via systray |

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run Linter

```bash
make lint
```

Requires `golangci-lint` to be installed.

### Project Structure

```
claude-notifier/
├── cmd/claude-notifier/  # Main entry point
├── internal/
│   ├── app/           # App orchestration and config
│   ├── notification/  # Store and sound player
│   ├── server/        # HTTP server
│   └── ui/            # macOS menu bar (systray)
├── assets/            # Icon files
├── scripts/           # Hook script and LaunchAgent
└── Makefile
```

## Menu Bar Interface

When you click the claude-notifier icon in the menu bar, you see:

- **Header** - Shows unread count (e.g., "Notifications (3 unread)")
- **Notifications** - Listed newest first with indicators:
  - `●` = unread
  - `○` = read
- **Mark All Read** - Marks all notifications as read
- **Clear History** - Removes all notifications
- **About** - Shows version information
- **Quit** - Stops the claude-notifier

Clicking a notification marks it as read.

## Requirements

- macOS (uses `afplay` for sounds and getlantern/systray for menu bar)
- Go 1.21+

## License

MIT
