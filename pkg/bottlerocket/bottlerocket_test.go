package bottlerocket

import (
	_ "embed"
	"fmt"
	"github.com/pipetail/bottlerocket-updater/pkg/httpmock"
	"testing"
)

//go:embed assets/idle.json
var idleResponseString string

func TestActivateUpdate(t *testing.T) {
	client := &httpmock.Client{}
	client.SetStatusCode(204)

	// pass
	err := ActivateUpdate(client)
	if err != nil {
		t.Errorf("error was not expected here: %s", err.Error())
	}

	// set different status code
	client.SetStatusCode(500)
	err = ActivateUpdate(client)
	if err == nil {
		t.Errorf("an error was expected here but got nil")
	}

	// simulate the case when client returns non-nil error
	client.SetStatusCode(500)
	client.SetError(fmt.Errorf("random error"))
	if err == nil {
		t.Errorf("an error was expected here but got nil")
	}
}

func TestGetUpdatesStatus(t *testing.T) {
	client := &httpmock.Client{}
	client.SetStatusCode(200)
	client.SetBody(idleResponseString)

	// test idle
	status, err := GetUpdatesStatus(client)
	if err != nil {
		t.Errorf("an error was na expected here: %s", err.Error())
	}

	if status.State != "Idle" {
		t.Errorf("Idle State was expected here but got %s", status.State)
	}
}
