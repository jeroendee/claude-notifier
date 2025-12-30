// Package app provides application orchestration for claude-notifier.
//
// It coordinates the lifecycle of all components: HTTP server, notification
// store, system tray, and menu. The [App] type uses a builder pattern for
// dependency injection, allowing components to be set individually.
//
// # Component Interfaces
//
// The package defines interfaces for each component:
//
//   - [Store] for notification storage
//   - [Server] for HTTP request handling
//   - [Systray] for system tray management
//   - [Menu] for menu UI updates
//
// # Signal Handling
//
// [App.Run] sets up handlers for SIGTERM and SIGINT, triggering graceful
// shutdown of all components when received.
//
// # Configuration
//
// [LoadConfig] reads configuration from environment variables:
//
//	CLAUDE_NOTIFIER_PORT        HTTP server port (default: 19199)
//	CLAUDE_NOTIFIER_SOUND       Path to notification sound file
//	CLAUDE_NOTIFIER_MAX_HISTORY Maximum stored notifications (default: 50)
package app
