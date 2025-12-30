// Package server provides an HTTP server for receiving Claude Code webhooks.
//
// The server listens on localhost only and provides endpoints for notification
// management. It integrates with a notification [Store] and [SoundPlayer] to
// store messages and play sounds when notifications arrive.
//
// # Endpoints
//
// The server exposes three endpoints:
//
//	POST /notify    Add a notification and play sound
//	GET  /health    Health check returning {"status": "ok"}
//	POST /clear     Clear all notifications
//
// # Request Format
//
// The /notify endpoint accepts JSON with a message field:
//
//	{"message": "Build completed successfully"}
//
// # Timeouts
//
// Read and write timeouts are set to 10 seconds. Graceful shutdown
// allows up to 5 seconds for in-flight requests to complete.
package server
