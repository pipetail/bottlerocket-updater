package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UnmarshalResponse creates struct directly from the HTTP response
func UnmarshalResponse(response *http.Response, v interface{}) error {

	// get the body as []byte
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read the body: %s", err.Error())
	}

	// try to Unmarshal to struct
	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("could not unmarshal: %s", err.Error())
	}

	return nil
}
