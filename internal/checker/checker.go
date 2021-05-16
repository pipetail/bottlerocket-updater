package checker

import (
	"fmt"
	"github.com/pipetail/bottlerocket-updater/internal/common"
	"github.com/pipetail/bottlerocket-updater/pkg/bottlerocket"
	"log"
	"os"
)

type Config struct {
	SocketPath string
}

func RealMain(config Config) {
	// prepare HTTP client with the special UDS configuration
	client := common.GetHTTPClient(config.SocketPath)

	// get status from bottlerocket api
	status, err := bottlerocket.GetUpdatesStatus(client)
	if err != nil {
		log.Printf("could not get update status %s", err.Error())
		refresh(client) // always issue the refresh command
		os.Exit(1)
	}

	if status.State == "Ready" {
		fmt.Println("update is ready")
		refresh(client) // always issue the refresh command
		os.Exit(0)
	}

	log.Printf("irelevant update status '%s'", status.State)
	refresh(client) // always issue the refresh command
	os.Exit(1)      // if not 'Ready' we can't exit with 0 since it would cause reboot
}

// refresh issues the refresh requests to the bottlerocket api
// We want to run this in best effort mode hence no error checks
// here
func refresh(client bottlerocket.HTTPClient) {
	err := bottlerocket.RefreshUpdates(client)
	if err != nil {
		log.Println("could not refresh update status")
	}
}
