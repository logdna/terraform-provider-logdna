package logdna

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
)

// Client used to make HTTP requests to the configuration api
type Client struct {
	ServiceKey string
	HTTPClient *http.Client
	ApiUrl string
	Method string
	Body interface{}
}

// AlertResponsePayload contains alert data returned from the config-api
type AlertResponsePayload struct {
	Channels []ChannelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Name     string            `json:"name,omitempty"`
	PresetID string            `json:"presetid,omitempty"`
}

func (c *Client) MakeRequest() ([]byte, error) {
	payloadBuf := bytes.NewBuffer([]byte{})
	if c.Body != nil {
		pbytes, err := json.Marshal(c.Body)
		if err != nil {
			return nil, err
		}
		payloadBuf = bytes.NewBuffer(pbytes)
	}

	req, err := http.NewRequest(c.Method, c.ApiUrl, payloadBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("servicekey", c.ServiceKey)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error during HTTP request: %s, %+v", err, c)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s, status: %d, body: %s", c.Method, c.ApiUrl, res.StatusCode, body)
	}

	return body, err
}

// MakeRequestAlert makes a HTTP request to the config-api with alert payload data and parses and returns the response
func MakeRequestAlert(c *Client, url string, urlsuffix string, method string, payload ViewRequest) (string, error) {
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
func (c *Client) CreateAlert(url string, payload ViewRequest) (string, error) {
	result, err := MakeRequestAlert(c, url, "/v1/config/presetalert", "POST", payload)
	return result, err
}

// UpdateAlert updates a Preset Alert with the provided presetID and payload
func (c *Client) UpdateAlert(url string, presetID string, payload ViewRequest) error {
	_, err := MakeRequestAlert(c, url, "/v1/config/presetalert/"+presetID, "PUT", payload)
	return err
}

// DeleteAlert deletes an alert with the provided presetID
func (c *Client) DeleteAlert(url, presetID string) error {
	_, err := MakeRequestAlert(c, url, "/v1/config/presetalert/"+presetID, "DELETE", ViewRequest{})
	return err
}
