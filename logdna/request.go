package logdna

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpRequest func(string, string, io.Reader) (*http.Request, error)
type bodyReader func(io.Reader) ([]byte, error)
type jsonMarshal func(interface{}) ([]byte, error)
type httpClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}

// Configuration for the HTTP client used to make requests to remote resources
type requestConfig struct {
	serviceKey          string
	iamtoken            string
	cloud_resource_name string
	httpClient          httpClientInterface
	apiURL              string
	method              string
	body                interface{}
	httpRequest         httpRequest
	bodyReader          bodyReader
	jsonMarshal         jsonMarshal
}

// newRequestConfig abstracts the struct creation to allow for mocking
func newRequestConfig(pc *providerConfig, method string, uri string, body interface{}, mutators ...func(*requestConfig)) *requestConfig {
	rc := &requestConfig{
		serviceKey:          pc.serviceKey,
		iamtoken:            pc.iamtoken,
		cloud_resource_name: pc.cloud_resource_name,
		httpClient:          pc.httpClient,
		apiURL:              fmt.Sprintf("%s%s", pc.baseURL, uri), // uri should have a preceding slash (/)
		method:              method,
		body:                body,
		httpRequest:         http.NewRequest,
		bodyReader:          io.ReadAll,
		jsonMarshal:         json.Marshal,
	}

	// Used during testing only; Allow mutations passed in by tests
	for _, mutator := range mutators {
		mutator(rc)
	}
	return rc
}

func (c *requestConfig) MakeRequest() ([]byte, error) {
	payloadBuf := bytes.NewBuffer([]byte{})
	if c.body != nil {
		pbytes, err := c.jsonMarshal(c.body)
		if err != nil {
			return nil, err
		}
		payloadBuf = bytes.NewBuffer(pbytes)
	}

	req, err := c.httpRequest(c.method, c.apiURL, payloadBuf)
	if err != nil {
		return nil, err
	}
	if payloadBuf.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set the correct authorization headers depending on what has been passed in
	// the provider config
	if c.serviceKey != "" {
		req.Header.Set("servicekey", c.serviceKey)
	} else if c.iamtoken != "" && c.cloud_resource_name != "" {
		req.Header.Set("Authorization", "Bearer "+c.iamtoken)
		req.Header.Set("cloud-resource-name", c.cloud_resource_name)
	} else {
		err := fmt.Errorf("expected either servicekey or iamtoken to be set")
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during HTTP request: %s", err)
	}
	defer res.Body.Close()

	body, err := c.bodyReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP response: %s, %s", err, string(body))
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s, status %d NOT OK! %s", c.method, c.apiURL, res.StatusCode, string(body))
	}
	return body, err
}
