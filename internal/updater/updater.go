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

		stage, stageStatus := HandleStates(client, status)
		log.Printf("stage '%s' executed with status: %v", stage, stageStatus)
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

// HandleStates check returned Update Status and executes the required
// action that needs to be executed afterwards
func HandleStates(client bottlerocket.HTTPClient, status bottlerocket.UpdateStatus) (string, bool) {
	var err error
	var action string
	if status.State == "Available" {
		action = "prepare"
		log.Println("preparing the update")
		err = bottlerocket.PrepareUpdate(client)
	}

	if status.State == "Staged" {
		action = "activate"
		log.Println("activating the update")
		err = bottlerocket.ActivateUpdate(client)
	}

	if status.State == "Ready" {
		action = "noop_ready"
		log.Println("update ready")
	}

	if status.State == "Idle" {
		action = "noop_idle"
		log.Println("update idle")
	}

	// always refresh the updates
	refresh(client)
	if err != nil {
		log.Printf("there was an error during '%s' stage: %s", action, err.Error())
		return action, false
	}

	return action, true
}