// Package mode manages the application runtime mode (production/development) and debug flag.
package mode

import (
	"os"
	"runtime"
	"strings"

	"github.com/casapps/cvedex/src/config"
)

var (
	currentMode  = Production
	debugEnabled = false
)

// AppMode represents the application operational mode.
type AppMode int

const (
	// Production is the default mode — strict security, minimal logging.
	Production AppMode = iota
	// Development enables verbose logging and relaxed security for local dev.
	Development
)

func (m AppMode) String() string {
	switch m {
	case Development:
		return "development"
	default:
		return "production"
	}
}

// SetAppMode sets the application mode from a string shortcut.
func SetAppMode(m string) {
	switch strings.ToLower(m) {
	case "dev", "development":
		currentMode = Development
	default:
		currentMode = Production
	}
	updateProfilingSettings()
}

// SetDebugEnabled enables or disables debug mode.
func SetDebugEnabled(enabled bool) {
	debugEnabled = enabled
	updateProfilingSettings()
}

// updateProfilingSettings enables or disables runtime profiling based on the debug flag.
func updateProfilingSettings() {
	if debugEnabled {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
	} else {
		runtime.SetBlockProfileRate(0)
		runtime.SetMutexProfileFraction(0)
	}
}

// GetCurrentAppMode returns the current application mode.
func GetCurrentAppMode() AppMode {
	return currentMode
}

// IsAppModeDev returns true if in development mode.
func IsAppModeDev() bool {
	return currentMode == Development
}

// IsAppModeProd returns true if in production mode.
func IsAppModeProd() bool {
	return currentMode == Production
}

// IsDebugEnabled returns true if debug mode is enabled (--debug or DEBUG=true).
func IsDebugEnabled() bool {
	return debugEnabled
}

// GetAppModeString returns the mode string with debug suffix when enabled.
func GetAppModeString() string {
	s := currentMode.String()
	if debugEnabled {
		s += " [debugging]"
	}
	return s
}

// FromEnv sets mode and debug from environment variables.
func FromEnv() {
	if m := os.Getenv("MODE"); m != "" {
		SetAppMode(m)
	}
	if config.IsTruthy(os.Getenv("DEBUG")) {
		SetDebugEnabled(true)
	}
}
