package logdna

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"time"
)

type HttpRequest func(string, string, io.Reader) (*http.Request, error)
type BodyReader func(io.Reader) ([]byte, error)
type jsonMarshal func(interface{}) ([]byte, error)
type httpClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}

// Client used to make HTTP requests to the configuration api
type requestConfig struct {
	ServiceKey string
	HTTPClient httpClientInterface
	ApiUrl     string
	Method     string
	Body       interface{}
	HttpRequest HttpRequest
	BodyReader BodyReader
	jsonMarshal jsonMarshal
}

// AlertResponsePayload contains alert data returned from the config-api
type AlertResponsePayload struct {
	Channels []ChannelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Name     string            `json:"name,omitempty"`
	PresetID string            `json:"presetid,omitempty"`
}

func NewRequestConfig(pc *providerConfig, method string, uri string, body interface{}, mutators ...func(*requestConfig)) *requestConfig {
	rc := &requestConfig{
		ServiceKey: pc.ServiceKey,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		ApiUrl: fmt.Sprintf("%s/%s", pc.Host, uri),
		Method: method,
		Body: body,
		HttpRequest: http.NewRequest,
		BodyReader: ioutil.ReadAll,
		jsonMarshal: json.Marshal,
	}

	// Testing only; Allow mutations passed in by tests
	for _, mutator := range mutators {
		mutator(rc)
	}
	return rc
}

func (c *requestConfig) MakeRequest() ([]byte, error) {
	payloadBuf := bytes.NewBuffer([]byte{})
	if c.Body != nil {
		pbytes, err := c.jsonMarshal(c.Body)
		if err != nil {
			return nil, err
		}
		payloadBuf = bytes.NewBuffer(pbytes)
	}

	req, err := c.HttpRequest(c.Method, c.ApiUrl, payloadBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("servicekey", c.ServiceKey)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error during HTTP request: %s, %+v", err, c)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s, status NOT OK: %d", c.Method, c.ApiUrl, res.StatusCode)
	}
	defer res.Body.Close()

	body, err := c.BodyReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing HTTP response: %s, %+v", err, c)
	}

	return body, err
}

// MakeRequestAlert makes a HTTP request to the config-api with alert payload data and parses and returns the response
func MakeRequestAlert(c *requestConfig, url string, urlsuffix string, method string, payload ViewRequest) (string, error) {
	pbytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, url+urlsuffix, bytes.NewBuffer(pbytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("servicekey", c.ServiceKey)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf(`Error with alert: %s`, err)
	}
	defer resp.Body.Close()
	var result AlertResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(result.Error)
	}

	return result.PresetID, nil
}

// CreateAlert creates a Preset Alert with the provided payload
func (c *requestConfig) CreateAlert(url string, payload ViewRequest) (string, error) {
	result, err := MakeRequestAlert(c, url, "/v1/config/presetalert", "POST", payload)
	return result, err
}

// UpdateAlert updates a Preset Alert with the provided presetID and payload
func (c *requestConfig) UpdateAlert(url string, presetID string, payload ViewRequest) error {
	_, err := MakeRequestAlert(c, url, "/v1/config/presetalert/"+presetID, "PUT", payload)
	return err
}

// DeleteAlert deletes an alert with the provided presetID
func (c *requestConfig) DeleteAlert(url, presetID string) error {
	_, err := MakeRequestAlert(c, url, "/v1/config/presetalert/"+presetID, "DELETE", ViewRequest{})
	return err
}
