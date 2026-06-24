// Package main is the entry point for cvedex — an authoritative DNS server
// that exposes the full CVE corpus as queryable DNS names under the .cve TLD.
//
// This software is licensed under the MIT License.
// See LICENSE.md for details.
package main

import (
	"fmt"
	"os"

	"github.com/casapps/cvedex/src/config"
	"github.com/casapps/cvedex/src/mode"
	"github.com/casapps/cvedex/src/service"
	"github.com/casapps/cvedex/src/signal"
)

// Build info — set via -ldflags at build time.
var (
	Version      = "devel"
	CommitID     = "N/A"
	BuildDate    = "N/A"
	OfficialSite = ""
)

func main() {
	cfg, err := config.Load(Version, CommitID, BuildDate, OfficialSite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cvedex: config error: %v\n", err)
		os.Exit(1)
	}

	mode.FromEnv()
	if cfg.Mode != "" {
		mode.SetAppMode(cfg.Mode)
	}
	if cfg.Debug {
		mode.SetDebugEnabled(true)
	}

	svc, err := service.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cvedex: init error: %v\n", err)
		os.Exit(1)
	}

	signal.Handle(svc)

	if err := svc.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cvedex: %v\n", err)
		os.Exit(1)
	}
}
