package main

import (
	"flag"
	"github.com/pipetail/bottlerocket-updater/internal/updater"
)

var socketPath *string

func init() {
	socketPath = flag.String("socket", "/run/api.sock", "path to bottlerocket socket")
	flag.Parse()
}

func main() {
	config := updater.Config{
		SocketPath: *socketPath,
	}
	updater.RealMain(config)
}
