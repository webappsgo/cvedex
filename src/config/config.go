package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the resolved runtime configuration for cvedex.
type Config struct {
	Version      string
	CommitID     string
	BuildDate    string
	OfficialSite string

	Mode  string
	Debug bool

	Address string
	Port    int

	ConfigDir string
	DataDir   string
	LogDir    string
	CacheDir  string
}

// IsDebug returns true when debug mode is active.
func (c *Config) IsDebug() bool {
	return c.Debug
}

// Sanitized returns a copy of Config with sensitive values redacted.
func (c *Config) Sanitized() map[string]any {
	return map[string]any{
		"version":       c.Version,
		"commit_id":     c.CommitID,
		"build_date":    c.BuildDate,
		"official_site": c.OfficialSite,
		"mode":          c.Mode,
		"debug":         c.Debug,
		"address":       c.Address,
		"port":          c.Port,
		"config_dir":    c.ConfigDir,
		"data_dir":      c.DataDir,
	}
}

// Load builds a Config from environment variables and CLI flags.
// Build-time variables are injected via -ldflags at compile time.
func Load(version, commitID, buildDate, officialSite string) (*Config, error) {
	cfg := &Config{
		Version:      version,
		CommitID:     commitID,
		BuildDate:    buildDate,
		OfficialSite: officialSite,
		Mode:         "production",
		Debug:        false,
		Address:      "0.0.0.0",
		Port:         64580,
		ConfigDir:    defaultConfigDir(),
		DataDir:      defaultDataDir(),
		LogDir:       defaultLogDir(),
		CacheDir:     defaultCacheDir(),
	}

	if m := os.Getenv("MODE"); m != "" {
		cfg.Mode = strings.ToLower(m)
	}

	if d := os.Getenv("DEBUG"); d != "" {
		cfg.Debug = IsTruthy(d)
	}

	if addr := os.Getenv("ADDRESS"); addr != "" {
		cfg.Address = addr
	}

	if portStr := os.Getenv("PORT"); portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err == nil && p > 0 && p < 65536 {
			cfg.Port = p
		}
	}

	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		cfg.ConfigDir = dir
	}

	if dir := os.Getenv("DATA_DIR"); dir != "" {
		cfg.DataDir = dir
	}

	if dir := os.Getenv("LOG_DIR"); dir != "" {
		cfg.LogDir = dir
	}

	if dir := os.Getenv("CACHE_DIR"); dir != "" {
		cfg.CacheDir = dir
	}

	return cfg, nil
}

func defaultConfigDir() string {
	if os.Getuid() == 0 {
		return "/etc/casapps/cvedex"
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg + "/casapps/cvedex"
	}
	home, _ := os.UserHomeDir()
	return home + "/.config/casapps/cvedex"
}

func defaultDataDir() string {
	if os.Getuid() == 0 {
		return "/var/lib/casapps/cvedex"
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return xdg + "/casapps/cvedex"
	}
	home, _ := os.UserHomeDir()
	return home + "/.local/share/casapps/cvedex"
}

func defaultLogDir() string {
	if os.Getuid() == 0 {
		return "/var/log/casapps/cvedex"
	}
	return defaultDataDir() + "/log"
}

func defaultCacheDir() string {
	if os.Getuid() == 0 {
		return "/var/cache/casapps/cvedex"
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg + "/casapps/cvedex"
	}
	home, _ := os.UserHomeDir()
	return home + "/.cache/casapps/cvedex"
}
