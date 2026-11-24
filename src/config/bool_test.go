package config

import "testing"

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "true", input: "true"},
		{name: "numeric one", input: "1"},
		{name: "yes with whitespace and mixed case", input: "  YeS\t"},
		{name: "on", input: "on"},
		{name: "enable", input: "enable"},
		{name: "enabled", input: "enabled"},
		{name: "active", input: "active"},
		{name: "oui", input: "oui"},
		{name: "ja", input: "ja"},
		{name: "si", input: "si"},
		{name: "da", input: "da"},
		{name: "tak", input: "tak"},
		{name: "ano", input: "ano"},
		{name: "sim", input: "sim"},
		{name: "evet", input: "evet"},
		{name: "ayee", input: "ayee"},
		{name: "yeah", input: "yeah"},
		{name: "yep", input: "yep"},
		{name: "sure", input: "sure"},
		{name: "affirmative", input: "affirmative"},
		{name: "positive", input: "positive"},
		{name: "correct", input: "correct"},
		{name: "absolutely", input: "absolutely"},
		{name: "certainly", input: "certainly"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsTruthy(tt.input) {
				t.Errorf("IsTruthy(%q) = false, want true", tt.input)
			}
		})
	}

	for _, input := range []string{"", "false", "0", "maybe", "enabled?"} {
		t.Run("rejects "+input, func(t *testing.T) {
			if IsTruthy(input) {
				t.Errorf("IsTruthy(%q) = true, want false", input)
			}
		})
	}
}

func TestIsFalsy(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "false", input: "false"},
		{name: "numeric zero", input: "0"},
		{name: "no with whitespace and mixed case", input: " No "},
		{name: "off", input: "off"},
		{name: "disable", input: "disable"},
		{name: "disabled", input: "disabled"},
		{name: "inactive", input: "inactive"},
		{name: "non", input: "non"},
		{name: "nein", input: "nein"},
		{name: "nee", input: "nee"},
		{name: "ne", input: "ne"},
		{name: "nu", input: "nu"},
		{name: "nao", input: "nao"},
		{name: "hayir", input: "hayir"},
		{name: "nope", input: "nope"},
		{name: "nah", input: "nah"},
		{name: "never", input: "never"},
		{name: "negative", input: "negative"},
		{name: "incorrect", input: "incorrect"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsFalsy(tt.input) {
				t.Errorf("IsFalsy(%q) = false, want true", tt.input)
			}
		})
	}

	for _, input := range []string{"", "true", "1", "maybe", "disabled?"} {
		t.Run("rejects "+input, func(t *testing.T) {
			if IsFalsy(input) {
				t.Errorf("IsFalsy(%q) = true, want false", input)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue bool
		wantOK    bool
	}{
		{name: "truthy", input: "YES", wantValue: true, wantOK: true},
		{name: "falsy", input: "OFF", wantValue: false, wantOK: true},
		{name: "unrecognized", input: "sometimes", wantValue: false, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOK := ParseBool(tt.input)
			if gotValue != tt.wantValue || gotOK != tt.wantOK {
				t.Errorf("ParseBool(%q) = (%t, %t), want (%t, %t)", tt.input, gotValue, gotOK, tt.wantValue, tt.wantOK)
			}
		})
	}
}
