package main

import "github.com/mylinden-tech/linden-cli/internal/cli"

// version, commit, and date are injected at build time via -ldflags (see .goreleaser.yml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.Execute(cli.BuildInfo{Version: version, Commit: commit, Date: date})
}
