// Package main is the entry point for cvedex-cli — the CLI client for cvedex.
// See AI.md PART 33 for client details.
package main

import (
	"fmt"
	"os"
)

// Build info — set via -ldflags at build time.
var (
	Version      = "devel"
	CommitID     = "N/A"
	BuildDate    = "N/A"
	OfficialSite = ""
)

func main() {
	fmt.Fprintf(os.Stderr, "cvedex-cli %s — not yet implemented\n", Version)
	os.Exit(1)
}
