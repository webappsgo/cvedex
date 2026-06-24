// Package config provides configuration loading and helper utilities for cvedex.
package config

import "strings"

// IsTruthy returns true if s represents a truthy boolean value.
// Accepts an extensive vocabulary to match operator expectations across locales.
func IsTruthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes", "on", "enable", "enabled", "active",
		"oui", "ja", "si", "da", "tak", "ano", "sim", "evet",
		"ayee", "yeah", "yep", "sure", "affirmative", "positive",
		"correct", "absolutely", "certainly":
		return true
	}
	return false
}

// IsFalsy returns true if s represents a falsy boolean value.
func IsFalsy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "false", "0", "no", "off", "disable", "disabled", "inactive",
		"non", "nein", "nee", "ne", "nu", "nao", "hayir",
		"nope", "nah", "never", "negative", "incorrect":
		return true
	}
	return false
}

// ParseBool parses s as a boolean, returning the result and whether parsing succeeded.
// Returns (false, false) if s is not a recognized truthy or falsy value.
func ParseBool(s string) (value bool, ok bool) {
	if IsTruthy(s) {
		return true, true
	}
	if IsFalsy(s) {
		return false, true
	}
	return false, false
}
