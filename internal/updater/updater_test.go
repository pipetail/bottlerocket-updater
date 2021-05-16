package updater

import (
	_ "embed"
	"github.com/pipetail/bottlerocket-updater/pkg/bottlerocket"
	"github.com/pipetail/bottlerocket-updater/pkg/httpmock"
	"testing"
)

//go:embed assets/idle.json
var idleResponseString string

func TestHandleStatesIdle(t *testing.T) {

	client := &httpmock.Client{}
	client.SetStatusCode(200)
	client.SetBody(idleResponseString)

	status, err := bottlerocket.GetUpdatesStatus(client)
	if err != nil {
		t.Errorf("error was not expected here: %s", err.Error())
	}

	// reconfigure client for

	stage, stageStatus := HandleStates(client, status)
	if stage != "noop_idle" {
		t.Errorf("such stage title was not expected here: %s", stage)
	}

	if ! stageStatus {
		t.Errorf("successful status was expected here")
	}
}
