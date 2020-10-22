package logdna

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type viewPayload struct {
	Apps     []string  `json:"apps,omitempty"`
	Category []string  `json:"category,omitempty"`
	Channels []channel `json:"channels,omitempty"`
	Hosts    []string  `json:"hosts,omitempty"`
	Levels   []string  `json:"levels,omitempty"`
	Name     string    `json:"name,omitempty"`
	Query    string    `json:"query,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
}

type viewResponsePayload struct {
	Apps     []string          `json:"apps,omitempty"`
	Category []string          `json:"category,omitempty"`
	Channels []channelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Hosts    []string          `json:"hosts,omitempty"`
	Levels   []string          `json:"levels,omitempty"`
	Name     string            `json:"name,omitempty"`
	Query    string            `json:"query,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	Viewid   string            `json:"viewid,omitempty"`
}

type alertResponsePayload struct {
	Channels []channelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Name     string            `json:"name,omitempty"`
	Presetid string            `json:"presetid,omitempty"`
}

type channelResponse struct {
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
	Triggerinterval int               `json:"triggerinterval,omitempty"`
	Triggerlimit    int               `json:"triggerlimit,omitempty"`
	Timezone        string            `json:"timezone,omitempty"`
	URL             string            `json:"url,omitempty"`
}

// Client used to make http requests to the configuration api
type Client struct {
	servicekey string
	httpClient *http.Client
}

// MakeRequestView makes a request to the config-api and parses and returns the response
func MakeRequestView(c *Client, url string, urlsuffix string, method string, payload viewPayload) (string, error) {
	pbytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(method, url+urlsuffix, bytes.NewBuffer(pbytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("servicekey", c.servicekey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf(`Error with view: %s`, err)
	}
	defer resp.Body.Close()
	var result viewResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(result.Error)
	}

	return result.Viewid, nil
}

// CreateView creates a view
func (c *Client) CreateView(url string, payload viewPayload) (string, error) {
	result, err := MakeRequestView(c, url, "/v1/config/view", "POST", payload)
	return result, err
}

// UpdateView updates a view
func (c *Client) UpdateView(url string, viewid string, payload viewPayload) error {
	_, err := MakeRequestView(c, url, "/v1/config/view/"+viewid, "PUT", payload)
	return err

}

// DeleteView deletes a view with a provided view id
func (c *Client) DeleteView(url, viewid string) error {
	_, err := MakeRequestView(c, url, "/v1/config/view/"+viewid, "DELETE", viewPayload{})
	return err
}
