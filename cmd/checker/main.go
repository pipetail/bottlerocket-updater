package main

import (
	"flag"
	"github.com/pipetail/bottlerocket-updater/internal/checker"
)

var socketPath *string

func init() {
	socketPath = flag.String("socket", "/run/api.sock", "path to bottlerocket socket")
	flag.Parse()
}

func main() {
	config := checker.Config{
		SocketPath: *socketPath,
	}
	checker.RealMain(config)
}
