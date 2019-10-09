package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"chainup.dev/lib/log"
	"github.com/pkg/errors"
)

var (
	headers map[string]string
	baseURL string
)

func init() {
	headers = make(map[string]string)
	headers["Content-Type"] = "application/json"
}

// SendGet is a shorthand for sending a GET HTTP request.
func SendGet(url string, expectCode int, dest interface{}) error {
	req, err := http.NewRequest("GET", withBaseURL(url), nil)
	if err != nil {
		return errors.Wrap(err, "create get http request")
	}

	return sendRequest(req, expectCode, dest)
}

// SendPost is a shorthand for sending a POST HTTP request.
func SendPost(url string, src interface{}, expectCode int, dest interface{}) error {
	body, err := RequestBody(src)
	if err != nil {
		return errors.Wrap(err, "prepare http request body")
	}

	req, err := http.NewRequest("POST", withBaseURL(url), body)
	if err != nil {
		return errors.Wrap(err, "create post http request")
	}

	return sendRequest(req, expectCode, dest)
}

// sendRequest is a utility function for sending an HTTP request.
func sendRequest(req *http.Request, expectCode int, dest interface{}) error {
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send http request")
	}
	defer log.Closer(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read http response")
	}

	if gotCode := resp.StatusCode; gotCode != expectCode {
		return errors.Errorf("unexpected http status code: want %d, got %d, body '%s'",
			expectCode, gotCode, strings.Trim(string(body), " \n\t"))
	}

	if dest == nil {
		return nil
	}

	err = json.Unmarshal(body, dest)
	if err != nil {
		return errors.Wrapf(err, "decode http response: err %s, body '%s'", err, strings.Trim(string(body), " \n\t"))
	}

	return nil
}

// RequestBody marshals the request body into a format suitable
// for using in a http.NewRequest().
func RequestBody(src interface{}) (io.Reader, error) {
	data, err := json.Marshal(src)
	if err != nil {
		return nil, errors.Wrap(err, "encode http request body")
	}

	return bytes.NewBuffer(data), nil
}

// SetHeader to be sent with future requests.
func SetHeader(key string, value string) {
	headers[key] = value
}

// SetBaseURL configures the base URL to be used with future requests.
func SetBaseURL(url string) {
	baseURL = url
}

func withBaseURL(url string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(url, "/")
}
