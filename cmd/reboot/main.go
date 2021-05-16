package main

import (
	"flag"
	"github.com/pipetail/bottlerocket-updater/internal/reboot"
)

var socketPath *string

func init() {
	socketPath = flag.String("socket", "/run/api.sock", "path to bottlerocket socket")
	flag.Parse()
}

func main() {
	config := reboot.Config{
		SocketPath: *socketPath,
	}
	reboot.RealMain(config)
}
