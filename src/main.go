// Package main is the entry point for cvedex — an authoritative DNS server
// that exposes the full CVE corpus as queryable DNS names under the .cve TLD.
//
// This software is licensed under the MIT License.
// See LICENSE.md for details.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/casapps/cvedex/src/config"
	"github.com/casapps/cvedex/src/mode"
	"github.com/casapps/cvedex/src/service"
	"github.com/casapps/cvedex/src/signal"
)

// Build info — set via -ldflags at build time.
var (
	Version      = "devel"
	CommitID     = "N/A"
	// BuildDate is derived from BuildEpoch during initialization.
	BuildDate    = "N/A"
	// BuildEpoch is the Unix build timestamp set through -ldflags.
	BuildEpoch   = "0"
	OfficialSite = ""
)

func buildEpoch() int64 {
	epoch, err := strconv.ParseInt(BuildEpoch, 10, 64)
	if err != nil || epoch <= 0 {
		return 0
	}
	return epoch
}

func init() {
	setBuildDateFromEpoch()
}

func setBuildDateFromEpoch() {
	if epoch := buildEpoch(); epoch > 0 {
		BuildDate = time.Unix(epoch, 0).UTC().Format("2006-01-02T15:04:05Z")
	}
}

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
