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
	serviceKey    string
	enterpriseKey string
	httpClient    httpClientInterface
	apiURL        string
	method        string
	body          interface{}
	httpRequest   httpRequest
	bodyReader    bodyReader
	jsonMarshal   jsonMarshal
}

// newRequestConfig abstracts the struct creation to allow for mocking
func newRequestConfig(pc *providerConfig, method string, uri string, body interface{}, mutators ...func(*requestConfig)) *requestConfig {
	serviceKey := ""
	enterpriseKey := ""
	switch pc.orgType {
	case OrgTypeRegular:
		serviceKey = pc.serviceKey
	case OrgTypeEnterprise:
		enterpriseKey = pc.serviceKey
	}
	rc := &requestConfig{
		serviceKey:    serviceKey,
		enterpriseKey: enterpriseKey,
		httpClient:    pc.httpClient,
		apiURL:        fmt.Sprintf("%s%s", pc.baseURL, uri), // uri should have a preceding slash (/)
		method:        method,
		body:          body,
		httpRequest:   http.NewRequest,
		bodyReader:    io.ReadAll,
		jsonMarshal:   json.Marshal,
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

	if c.serviceKey != "" {
		req.Header.Set("servicekey", c.serviceKey)
	}
	if c.enterpriseKey != "" {
		req.Header.Set("enterprise-servicekey", c.enterpriseKey)
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
