// Package main is the entry point for cvedex-cli — the CLI client for cvedex.
// See AI.md PART 33 for client details.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
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
	fmt.Fprintf(os.Stderr, "cvedex-cli %s — not yet implemented\n", Version)
	os.Exit(1)
}
