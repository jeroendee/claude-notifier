// Package ui provides the macOS menu bar user interface.
//
// It uses [fyne.io/systray] to display a system tray icon and dropdown
// menu showing notification history.
//
// # Components
//
// [Systray] manages the system tray icon and event loop. It loads an icon
// file and runs the systray event loop.
//
// [Menu] builds and maintains the dropdown menu from the notification store.
// It displays notifications newest-first with read status indicators:
//
//   - ● (filled circle) for unread notifications
//   - ○ (empty circle) for read notifications
//
// # Menu Structure
//
// The menu contains:
//
//   - Header showing unread count
//   - Notification items (click to mark as read)
//   - "Mark All Read" action
//   - "Clear History" action
//   - "Quit" to exit the application
//
// [fyne.io/systray]: https://pkg.go.dev/fyne.io/systray
package ui
