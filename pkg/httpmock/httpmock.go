package httpmock

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct{
	Body io.ReadCloser
	StatusCode int
	Error error
	Count int // this property is useful to determine the number of calls
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	response := http.Response{
		Body: c.Body,
		StatusCode: c.StatusCode,
	}
	c.Count = c.Count + 1
	return &response, c.Error
}

func (c *Client) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	response := http.Response{
		Body: c.Body,
		StatusCode: c.StatusCode,
	}
	c.Count = c.Count + 1
	return &response, c.Error
}

func (c *Client) SetBody(body string) {
	c.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))
}

func (c *Client) SetError(err error) {
	c.Error = err
}

func (c *Client) ClearError() {
	c.Error = nil
}

func (c *Client) SetStatusCode(code int) {
	c.StatusCode = code
}