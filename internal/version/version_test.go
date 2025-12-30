package version

import "testing"

func TestGet_ReturnsInfo(t *testing.T) {
	t.Parallel()

	info := Get()

	if info.Version != Version {
		t.Errorf("Get().Version = %q, want %q", info.Version, Version)
	}
	if info.Commit != Commit {
		t.Errorf("Get().Commit = %q, want %q", info.Commit, Commit)
	}
	if info.BuildDate != BuildDate {
		t.Errorf("Get().BuildDate = %q, want %q", info.BuildDate, BuildDate)
	}
}

func TestInfo_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		info Info
		want string
	}{
		{
			name: "dev version",
			info: Info{Version: "dev"},
			want: "claude-notifier dev",
		},
		{
			name: "release version",
			info: Info{Version: "v1.0.0"},
			want: "claude-notifier v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.info.String(); got != tt.want {
				t.Errorf("Info.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
