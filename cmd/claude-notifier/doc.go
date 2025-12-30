// claude-notifier is a macOS menu bar application that displays notifications
// from Claude Code.
//
// It runs as a system tray application, listening for webhook requests
// and displaying notification history in a dropdown menu.
//
// # Usage
//
// Start claude-notifier:
//
//	claude-notifier
//
// The application listens on localhost:19199 by default. Claude Code
// sends POST requests to /notify with a JSON body containing a message.
//
// # Configuration
//
// Environment variables:
//
//	CLAUDE_NOTIFIER_PORT        HTTP server port (default: 19199)
//	CLAUDE_NOTIFIER_SOUND       Path to notification sound file
//	CLAUDE_NOTIFIER_MAX_HISTORY Maximum stored notifications (default: 50)
//
// # Endpoints
//
// The server exposes:
//
//	POST /notify    Add notification and play sound
//	GET  /health    Health check
//	POST /clear     Clear all notifications
package main
