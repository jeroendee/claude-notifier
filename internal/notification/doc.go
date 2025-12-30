// Package notification provides notification storage and sound playback.
//
// # Store
//
// [Store] manages notifications with thread-safe operations. Multiple
// goroutines may safely call Store methods simultaneously.
//
// The store maintains a maximum of 50 notifications using FIFO eviction:
// when the limit is reached, the oldest notification is removed to make
// room for new ones.
//
// Change callbacks can be registered via [Store.SetOnChange] to receive
// notifications when the store contents change.
//
// # Sound Playback
//
// [SoundPlayer] provides non-blocking sound playback on macOS using the
// afplay command. Playback errors are silently ignored to avoid blocking
// the notification flow.
package notification
