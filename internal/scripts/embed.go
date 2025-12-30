package scripts

import _ "embed"

//go:embed claude-hook.sh
var HookScript []byte

//go:embed com.dee.claude-notifier.plist
var PlistTemplate []byte
