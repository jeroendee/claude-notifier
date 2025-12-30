// Package version provides build-time version information.
//
// Version, commit hash, and build date are injected via ldflags at build time.
// Use [Get] to retrieve the current version info, and [Info.String] for a
// formatted display string.
//
// # Build Flags
//
// Set version info during build:
//
//	go build -ldflags "-X github.com/jeroendee/claude-notifier/internal/version.Version=1.0.0"
package version
