package updater

import (
	"github.com/pipetail/bottlerocket-updater/internal/common"
	"github.com/pipetail/bottlerocket-updater/pkg/bottlerocket"
	"log"
	"time"
)

type Config struct {
	SocketPath string
}

func RealMain(config Config) {
	// prepare HTTP client with the special UDS configuration
	client := common.GetHTTPClient(config.SocketPath)

	for range time.Tick(time.Minute * 10) {
		status, err := bottlerocket.GetUpdatesStatus(client)
		if err != nil {
			log.Printf("could not get status: %s", err.Error())
			refresh(client)
			continue
		}

		if status.State == "Available" {
			log.Println("preparing the update")
			// prepare update
			refresh(client)
			continue
		}

		if status.State == "Staged" {
			log.Println("activating the update")
			// activate update
			refresh(client)
			continue
		}

		if status.State == "Ready" {
			log.Printf("update already prepared")
			refresh(client)
			continue
		}

		log.Println("unknown state")
	}
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