package logdna

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Client used to make HTTP requests to the configuration api
type Client struct {
	ServiceKey string
	HTTPClient *http.Client
}

// ViewPayload contains view data (such as apps and categories) and is forwarded to the config-api
type ViewPayload struct {
	Apps     []string  `json:"apps,omitempty"`
	Category []string  `json:"category,omitempty"`
	Channels []Channel `json:"channels,omitempty"`
	Hosts    []string  `json:"hosts,omitempty"`
	Levels   []string  `json:"levels,omitempty"`
	Name     string    `json:"name,omitempty"`
	Query    string    `json:"query,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
}

// ViewResponsePayload contains view data returned from the config-api
type ViewResponsePayload struct {
	Apps     []string          `json:"apps,omitempty"`
	Category []string          `json:"category,omitempty"`
	Channels []ChannelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Hosts    []string          `json:"hosts,omitempty"`
	Levels   []string          `json:"levels,omitempty"`
	Name     string            `json:"name,omitempty"`
	Query    string            `json:"query,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	ViewID   string            `json:"viewID,omitempty"`
}

// AlertResponsePayload contains alert data returned from the config-api
type AlertResponsePayload struct {
	Channels []ChannelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Name     string            `json:"name,omitempty"`
	PresetID string            `json:"presetid,omitempty"`
}

// ChannelResponse contains channel data returned from the config-api
type ChannelResponse struct {
	AlertID         string            `json:"alertid,omitempty"`
	BodyTemplate    string            `json:"bodyTemplate,omitempty"`
	Emails          string            `json:"emails,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Immediate       bool              `json:"immediate,omitempty"`
	Integration     string            `json:"integration,omitempty"`
	Key             string            `json:"key,omitempty"`
	Method          string            `json:"method,omitempty"`
	Operator        string            `json:"operator,omitempty"`
	Terminal        bool              `json:"terminal,omitempty"`
	TriggerInterval int               `json:"triggerinterval,omitempty"`
	TriggerLimit    int               `json:"triggerlimit,omitempty"`
	Timezone        string            `json:"timezone,omitempty"`
	URL             string            `json:"url,omitempty"`
}

// MakeRequestAlert makes a HTTP request to the config-api with alert payload data and parses and returns the response
func MakeRequestAlert(c *Client, url string, urlsuffix string, method string, payload ViewPayload) (string, error) {
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

// MakeRequestView makes a HTTP request to the config-api with view payload data and parses and returns the response
func MakeRequestView(c *Client, url string, urlsuffix string, method string, payload ViewPayload) (string, error) {
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
		return "", fmt.Errorf(`Error with view: %s`, err)
	}
	defer resp.Body.Close()
	var result ViewResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(result.Error)
	}

	return result.ViewID, nil
}

// CreateAlert creates a Preset Alert with the provided payload
func (c *Client) CreateAlert(url string, payload ViewPayload) (string, error) {
	result, err := MakeRequestAlert(c, url, "/v1/config/alert", "POST", payload)
	return result, err
}

// UpdateAlert updates a Preset Alert with the provided presetID and payload
func (c *Client) UpdateAlert(url string, presetID string, payload ViewPayload) error {
	_, err := MakeRequestAlert(c, url, "/v1/config/alert/"+presetID, "PUT", payload)
	return err
}

// DeleteAlert deletes an alert with the provided presetID
func (c *Client) DeleteAlert(url, presetID string) error {
	_, err := MakeRequestAlert(c, url, "/v1/config/alert/"+presetID, "DELETE", ViewPayload{})
	return err
}

// CreateView creates a view with the provided payload
func (c *Client) CreateView(url string, payload ViewPayload) (string, error) {
	result, err := MakeRequestView(c, url, "/v1/config/view", "POST", payload)
	return result, err
}

// UpdateView updates the view with the provided viewID and payload
func (c *Client) UpdateView(url string, viewID string, payload ViewPayload) error {
	_, err := MakeRequestView(c, url, "/v1/config/view/"+viewID, "PUT", payload)
	return err

}

// DeleteView deletes a view with the provided viewID
func (c *Client) DeleteView(url, viewID string) error {
	_, err := MakeRequestView(c, url, "/v1/config/view/"+viewID, "DELETE", ViewPayload{})
	return err
}
