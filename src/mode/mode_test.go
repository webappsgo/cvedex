package mode

import (
	"os"
	"testing"
)

func resetMode(t *testing.T) {
	t.Helper()
	SetAppMode("production")
	SetDebugEnabled(false)
	t.Cleanup(func() {
		SetAppMode("production")
		SetDebugEnabled(false)
	})
}

func TestAppModeString(t *testing.T) {
	tests := []struct {
		name string
		mode AppMode
		want string
	}{
		{name: "production", mode: Production, want: "production"},
		{name: "development", mode: Development, want: "development"},
		{name: "unknown mode defaults to production", mode: AppMode(99), want: "production"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetAppMode(t *testing.T) {
	tests := []struct {
		name string
		input string
		want  AppMode
	}{
		{name: "dev shortcut", input: "dev", want: Development},
		{name: "case insensitive development", input: "DEVELOPMENT", want: Development},
		{name: "production", input: "production", want: Production},
		{name: "unknown defaults to production", input: "staging", want: Production},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetMode(t)
			SetAppMode(tt.input)
			if got := GetCurrentAppMode(); got != tt.want {
				t.Errorf("GetCurrentAppMode() = %v, want %v", got, tt.want)
			}
			if got := IsAppModeDev(); got != (tt.want == Development) {
				t.Errorf("IsAppModeDev() = %v, want %v", got, tt.want == Development)
			}
			if got := IsAppModeProd(); got != (tt.want == Production) {
				t.Errorf("IsAppModeProd() = %v, want %v", got, tt.want == Production)
			}
		})
	}
}

func TestSetDebugEnabledAndGetAppModeString(t *testing.T) {
	resetMode(t)
	SetAppMode("development")
	SetDebugEnabled(true)
	if !IsDebugEnabled() {
		t.Fatal("IsDebugEnabled() = false after enabling debug")
	}
	if got, want := GetAppModeString(), "development [debugging]"; got != want {
		t.Errorf("GetAppModeString() = %q, want %q", got, want)
	}

	SetDebugEnabled(false)
	if IsDebugEnabled() {
		t.Fatal("IsDebugEnabled() = true after disabling debug")
	}
	if got, want := GetAppModeString(), "development"; got != want {
		t.Errorf("GetAppModeString() = %q, want %q", got, want)
	}
}

func TestFromEnv(t *testing.T) {
	tests := []struct {
		name       string
		initial    AppMode
		initialDbg bool
		mode       string
		debug      string
		wantMode   AppMode
		wantDebug  bool
	}{
		{
			name:      "development and truthy debug",
			mode:      "development",
			debug:     " yes ",
			wantMode:  Development,
			wantDebug: true,
		},
		{
			name:      "empty mode and false debug preserve defaults",
			mode:      "",
			debug:     "false",
			wantMode:  Production,
			wantDebug: false,
		},
		{
			name:       "false debug does not clear an existing debug state",
			initial:    Development,
			initialDbg: true,
			mode:       "production",
			debug:      "no",
			wantMode:   Production,
			wantDebug:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetMode(t)
			if tt.initial == Development {
				SetAppMode("development")
			}
			SetDebugEnabled(tt.initialDbg)
			t.Setenv("MODE", tt.mode)
			t.Setenv("DEBUG", tt.debug)

			FromEnv()
			if got := GetCurrentAppMode(); got != tt.wantMode {
				t.Errorf("GetCurrentAppMode() = %v, want %v", got, tt.wantMode)
			}
			if got := IsDebugEnabled(); got != tt.wantDebug {
				t.Errorf("IsDebugEnabled() = %v, want %v", got, tt.wantDebug)
			}
		})
	}
}

func TestFromEnvReadsCurrentEnvironment(t *testing.T) {
	resetMode(t)
	t.Setenv("MODE", "dev")
	t.Setenv("DEBUG", "1")
	FromEnv()
	if os.Getenv("MODE") != "dev" || !IsAppModeDev() || !IsDebugEnabled() {
		t.Fatal("FromEnv() did not apply MODE and DEBUG")
	}
}
