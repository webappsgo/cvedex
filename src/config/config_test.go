package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfigIsDebug(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
	}{
		{name: "disabled", debug: false},
		{name: "enabled", debug: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Debug: tt.debug}
			if got := cfg.IsDebug(); got != tt.debug {
				t.Errorf("IsDebug() = %t, want %t", got, tt.debug)
			}
		})
	}
}

func TestConfigSanitized(t *testing.T) {
	cfg := &Config{
		Version:      "1.2.3",
		CommitID:     "abc1234",
		BuildDate:    "2026-11-23T12:00:00Z",
		OfficialSite: "https://example.invalid",
		Mode:         "development",
		Debug:        true,
		Address:      "127.0.0.1",
		Port:         8080,
		ConfigDir:    "/config",
		DataDir:      "/data",
		LogDir:       "/private/logs",
		CacheDir:     "/private/cache",
	}

	want := map[string]any{
		"version":      "1.2.3",
		"commit_id":    "abc1234",
		"build_date":   "2026-11-23T12:00:00Z",
		"official_site": "https://example.invalid",
		"mode":         "development",
		"debug":        true,
		"address":      "127.0.0.1",
		"port":         8080,
		"config_dir":   "/config",
		"data_dir":     "/data",
	}

	got := cfg.Sanitized()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sanitized() = %#v, want %#v", got, want)
	}
}

func TestLoadDefaults(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_CACHE_HOME", "")

	cfg, err := Load("1.2.3", "abc1234", "2026-11-23T12:00:00Z", "https://example.invalid")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{
		Version:      "1.2.3",
		CommitID:     "abc1234",
		BuildDate:    "2026-11-23T12:00:00Z",
		OfficialSite: "https://example.invalid",
		Mode:         "production",
		Debug:        false,
		Address:      "0.0.0.0",
		Port:         64580,
		ConfigDir:    defaultConfigDir(),
		DataDir:      defaultDataDir(),
		LogDir:       defaultLogDir(),
		CacheDir:     defaultCacheDir(),
	}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %#v, want %#v", cfg, want)
	}
}

func TestLoadEnvironmentOverrides(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("MODE", "DEVELOPMENT")
	t.Setenv("DEBUG", " yes ")
	t.Setenv("ADDRESS", "127.0.0.1")
	t.Setenv("PORT", "8443")
	t.Setenv("CONFIG_DIR", "/custom/config")
	t.Setenv("DATA_DIR", "/custom/data")
	t.Setenv("LOG_DIR", "/custom/log")
	t.Setenv("CACHE_DIR", "/custom/cache")

	cfg, err := Load("devel", "deadbee", "N/A", "")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Mode != "development" || !cfg.Debug || cfg.Address != "127.0.0.1" || cfg.Port != 8443 {
		t.Errorf("Load() basic overrides = (%q, %t, %q, %d), want (development, true, 127.0.0.1, 8443)", cfg.Mode, cfg.Debug, cfg.Address, cfg.Port)
	}
	if cfg.ConfigDir != "/custom/config" || cfg.DataDir != "/custom/data" || cfg.LogDir != "/custom/log" || cfg.CacheDir != "/custom/cache" {
		t.Errorf("Load() directory overrides = (%q, %q, %q, %q), want custom paths", cfg.ConfigDir, cfg.DataDir, cfg.LogDir, cfg.CacheDir)
	}
	if cfg.Version != "devel" || cfg.CommitID != "deadbee" || cfg.BuildDate != "N/A" {
		t.Errorf("Load() build info = (%q, %q, %q), want (devel, deadbee, N/A)", cfg.Version, cfg.CommitID, cfg.BuildDate)
	}
}

func TestLoadDebugAndPortFallbacks(t *testing.T) {
	tests := []struct {
		name      string
		debug     string
		port      string
		wantDebug bool
		wantPort  int
	}{
		{name: "unrecognized debug and empty port", debug: "sometimes", port: "", wantPort: 64580},
		{name: "false debug and zero port", debug: "no", port: "0", wantPort: 64580},
		{name: "negative port", debug: "true", port: "-1", wantDebug: true, wantPort: 64580},
		{name: "port above range", debug: "true", port: "65536", wantDebug: true, wantPort: 64580},
		{name: "non-numeric port", debug: "true", port: "invalid", wantDebug: true, wantPort: 64580},
		{name: "minimum port", debug: "false", port: "1", wantPort: 1},
		{name: "maximum port", debug: "false", port: "65535", wantPort: 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearConfigEnv(t)
			t.Setenv("DEBUG", tt.debug)
			t.Setenv("PORT", tt.port)

			cfg, err := Load("", "", "", "")
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.Debug != tt.wantDebug {
				t.Errorf("Debug = %t, want %t", cfg.Debug, tt.wantDebug)
			}
			if cfg.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", cfg.Port, tt.wantPort)
			}
		})
	}
}

func TestDefaultDirectories(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/xdg/config")
	t.Setenv("XDG_DATA_HOME", "/xdg/data")
	t.Setenv("XDG_CACHE_HOME", "/xdg/cache")

	configDir := "/xdg/config/casapps/cvedex"
	dataDir := "/xdg/data/casapps/cvedex"
	logDir := "/xdg/data/casapps/cvedex/log"
	cacheDir := "/xdg/cache/casapps/cvedex"
	if os.Getuid() == 0 {
		configDir = "/etc/casapps/cvedex"
		dataDir = "/var/lib/casapps/cvedex"
		logDir = "/var/log/casapps/cvedex"
		cacheDir = "/var/cache/casapps/cvedex"
	}

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "config", got: defaultConfigDir(), want: configDir},
		{name: "data", got: defaultDataDir(), want: dataDir},
		{name: "logs", got: defaultLogDir(), want: logDir},
		{name: "cache", got: defaultCacheDir(), want: cacheDir},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("default directory = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"MODE", "DEBUG", "ADDRESS", "PORT", "CONFIG_DIR", "DATA_DIR", "LOG_DIR", "CACHE_DIR",
	} {
		t.Setenv(name, "")
	}
}
