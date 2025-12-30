// Package scripts provides embedded script assets for the application.
//
// Scripts are embedded at compile time using Go's embed directive.
// Available assets:
//
//   - [HookScript] - Shell script for Claude notification hook
//   - [PlistTemplate] - LaunchAgent plist configuration template
//
// Extraction functions:
//
//   - [WriteHookScript] - Writes the hook script to a destination directory
//   - [WritePlist] - Writes the plist template with binary path substitution
package scripts
