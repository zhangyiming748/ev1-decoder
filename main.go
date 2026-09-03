package main

import "github.com/zhangyiming748/ev1-decoder/cmd"

// version is injected at build time via -ldflags "-X main.version={{.Version}}"
// (see .goreleaser.yml). Local builds keep the "dev" default.
var version = "dev"

func main() {
	cmd.Execute(version)
}
