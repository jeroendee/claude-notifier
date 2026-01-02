# Claude Notifier

A macOS menu bar notification app for Claude Code. Receive notifications when Claude completes tasks, with sound alerts and notification history accessible from your menu bar.

## Features

- **Menu bar notifications** - Persistent icon in macOS menu bar shows notification history
- **Dynamic icon state** - Icon changes from grey (idle) to orange when unread notifications exist
- **Sound alerts** - Plays a sound when notifications arrive (configurable)
- **Notification history** - Stores up to 50 notifications with read/unread status
- **HTTP API** - Receive notifications via simple HTTP POST requests
- **Claude Code integration** - Hook script for automatic task completion notifications
- **Single binary** - All assets embedded via `go:embed` for easy distribution

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

### Install Hook Scripts

```bash
make install-hook
```

This installs two hook scripts to `~/.claude/hooks/`:
- `claude-hook.sh` - For Stop events (task completions with SESSION summaries)
- `notification-hook.sh` - For idle prompts and permission requests

### Configure Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "Stop": [{
      "hooks": [{ "type": "command", "command": "~/.claude/hooks/claude-hook.sh" }]
    }],
    "Notification": [{
      "matcher": "idle_prompt",
      "hooks": [{ "type": "command", "command": "~/.claude/hooks/notification-hook.sh" }]
    }],
    "PermissionRequest": [{
      "matcher": "",
      "hooks": [{ "type": "command", "command": "~/.claude/hooks/notification-hook.sh" }]
    }]
  }
}
```

### Hook Types

| Hook | Trigger | Notification |
|------|---------|--------------|
| **Stop** | Claude finishes responding | `{time} {project} - {SESSION summary}` |
| **Notification** (`idle_prompt`) | Claude waiting for input 60+ seconds | `{time} {project} - Waiting for input` |
| **PermissionRequest** | Claude needs permission approval | `{time} {project} - Waiting for input` |

Now Claude Code will notify you when:
- Tasks complete (with context from SESSION markers)
- Claude has been waiting for your input for 60+ seconds
- Claude needs permission to proceed

## Summary Command

The `summary` subcommand extracts the last `[SESSION: ...]` marker from a Claude Code transcript file:

```bash
claude-notifier summary /path/to/transcript.jsonl
```

This outputs the last session marker found in assistant messages. The hook script uses this to add context to notifications.

**How it works:**
- Parses the JSONL transcript file
- Scans assistant messages for `[SESSION: ...]` markers
- Returns the last marker found (e.g., `[SESSION: Fixed bug in parser]` → `Fixed bug in parser`)

This marker-based approach provides instant extraction without API calls.

### Required: CLAUDE.md Configuration

For session summaries to work, Claude must be instructed to include `[SESSION: ...]` markers. Add this to your `~/.claude/CLAUDE.md` (or project-level CLAUDE.md):

```markdown
### Session Summary Protocol

**At the end of every task completion response**, include a summary line:

[SESSION: <5-7 word description of what was accomplished>]

**Subagent Delegation Rule:**
When work is delegated to subagents (via Task tool), the **main agent** is still
responsible for the session summary. Subagent responses are internal - the user only
sees the main agent's final response. The session summary must reflect the
completed work regardless of whether it was done directly or delegated.

**Format rules:**
- Exactly one `[SESSION: ...]` line per completion
- 5-7 words maximum
- Describes the outcome, not the process
- No special characters except hyphens

**Examples:**
- `[SESSION: Fixed transcript parser string handling]`
- `[SESSION: Added user authentication endpoint]`
- `[SESSION: Refactored database connection pool]`
```

Without this configuration, notifications will show the project name only.

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
│  • Store (history)    │    │  • Menu bar icon (dynamic state)   │
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
| Transcript | `internal/transcript` | Parse JSONL transcripts, extract headings |
| UI | `internal/ui` | macOS menu bar integration via systray |
| Assets | `internal/assets` | Embedded icons (normal and alert states) |
| Scripts | `internal/scripts` | Embedded hook script and LaunchAgent plist |

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
├── cmd/
│   ├── claude-notifier/     # Main entry point
│   ├── install-hook/        # Hook script installer
│   └── install-launchagent/ # LaunchAgent installer
├── internal/
│   ├── app/           # App orchestration and config
│   ├── assets/        # Embedded icons (go:embed)
│   ├── notification/  # Store and sound player
│   ├── scripts/       # Embedded hook script and plist (go:embed)
│   ├── server/        # HTTP server
│   ├── transcript/    # JSONL transcript parser for summaries
│   ├── ui/            # macOS menu bar (systray)
│   └── version/       # Version info (injected via ldflags)
└── Makefile
```

## Menu Bar Interface

The claude-notifier icon appears in your menu bar with two states:
- **Grey robot** - No unread notifications (idle)
- **Orange robot** - Unread notifications exist (alert)

When you click the icon, you see:

- **Header** - Shows unread count (e.g., "Notifications (3 unread)")
- **Notifications** - Listed newest first with indicators:
  - `●` = unread
  - `○` = read
- **Mark All Read** - Marks all notifications as read
- **Clear History** - Removes all notifications
- **About** - Shows version (e.g., "claude-notifier v1.0.0")
- **Quit** - Stops the claude-notifier

Clicking a notification marks it as read. When all notifications are read, the icon returns to grey.

## Requirements

- macOS (uses `afplay` for sounds and getlantern/systray for menu bar)
- Go 1.21+

## License

MIT
