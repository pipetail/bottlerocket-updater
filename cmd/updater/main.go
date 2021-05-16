package main

import (
	"flag"
	"github.com/pipetail/bottlerocket-updater/internal/updater"
)

var socketPath *string
var oneTime *bool

func init() {
	socketPath = flag.String("socket", "/run/api.sock", "path to bottlerocket socket")
	oneTime = flag.Bool("one-tine", false, "operation mode of the application")
	flag.Parse()
}

func main() {
	config := updater.Config{
		SocketPath: *socketPath,
	}
	if *oneTime {
		updater.OneTime(config)
	} else {
		updater.RealMain(config)
	}
}
