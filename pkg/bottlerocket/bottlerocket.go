package bottlerocket

import (
	"bytes"
	"fmt"
	"github.com/pipetail/bottlerocket-updater/pkg/json"
	"io"
	"net/http"
)

// Image contains details about the image staged in the primary or
// staging partition
type Image struct {
	Arch string `json:"arch"`
	Version string `json:"version"`
	Variant string `json:"variant"`
}

// Partition contains information about the primary or staging partition
type Partition struct {
	Image Image `json:"image"`
	NextToBoot bool `json:"next_to_boot"`
}

// UpdateStatus contains all information as reference in the
// official documentation https://github.com/bottlerocket-os/bottlerocket/blob/develop/sources/api/openapi.yaml
type UpdateStatus struct {
	State string `json:"update_state"`
	AvailableUpdates []string `json:"available_updates"`
	ChosenUpdate Image `json:"chosen_update"`
	ActivePartition Partition `json:"active_partition"`
	StagingPartition Partition `json:"staging_partition"`
}

// HTTPClient contains all the methods we need for the communication
// with bottlerocket API
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// GetUpdatesStatus get the bottlerocket update status
// please note that HTTP client needs to be constructed with
// Unix domain socket configuration
func GetUpdatesStatus(client HTTPClient) (UpdateStatus, error) {
	status := UpdateStatus{}
	res, err := client.Get("http://unix/updates/status")
	if err != nil {
		return status, fmt.Errorf("could not get status: %s", err.Error())
	}

	if res.StatusCode != 200 {
		return status, fmt.Errorf("api returned unexpected code %d", res.StatusCode)
	}

	err = json.UnmarshalResponse(res, &status)
	if err != nil {
		return status, fmt.Errorf("could not obtain status from the body: %s", err.Error())
	}

	return status, nil
}

// RefreshUpdates issues refresh requests to the update api as per official documentation
// https://github.com/bottlerocket-os/bottlerocket/blob/develop/sources/api/openapi.yaml#L328
func RefreshUpdates(client HTTPClient) error {
	// send the empty body
	body := bytes.NewBuffer([]byte{})
	res, err := client.Post("http://unix/actions/refresh-updates", "", body)
	if err != nil {
		return fmt.Errorf("could not refresh updates: %s", err.Error())
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("api returned unexpected code %d", res.StatusCode)
	}

	return nil
}

// PrepareUpdate installs updates to the staging partition as per
// https://github.com/bottlerocket-os/bottlerocket/blob/develop/sources/api/openapi.yaml#L340
func PrepareUpdate(client HTTPClient) error {
	// send the empty body
	body := bytes.NewBuffer([]byte{})
	res, err := client.Post("http://unix/actions/prepare-update", "", body)
	if err != nil {
		return fmt.Errorf("could not prepare update: %s", err.Error())
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("api returned unexpected code %d", res.StatusCode)
	}

	return nil
}

// ActivateUpdate marks the partition with the prepared update
// as the next boot target as per
// https://github.com/bottlerocket-os/bottlerocket/blob/develop/sources/api/openapi.yaml#L356
func ActivateUpdate(client HTTPClient) error {
	// send the empty body
	body := bytes.NewBuffer([]byte{})
	res, err := client.Post("http://unix/actions/activate-update", "", body)
	if err != nil {
		return fmt.Errorf("could not activate update: %s", err.Error())
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("api returned unexpected code %d", res.StatusCode)
	}

	return nil
}


// Reboot reboots the system as per
// https://github.com/bottlerocket-os/bottlerocket/blob/develop/sources/api/openapi.yaml#L318
func Reboot(client HTTPClient) error {
	// send the empty body
	body := bytes.NewBuffer([]byte{})
	res, err := client.Post("http://unix/actions/reboot", "", body)
	if err != nil {
		return fmt.Errorf("could not reboot: %s", err.Error())
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("api returned unexpected code %d", res.StatusCode)
	}

	return nil
}


