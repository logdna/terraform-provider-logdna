package logdna

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_View(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ViewResponsePayload{ViewID: "test123456"})
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test", Query: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	result, err := client.CreateView(ts.URL, payload)
	assert.Equal(t, nil, err)
	assert.Len(t, result, 10)
}

func TestClient_Alert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(AlertResponsePayload{PresetID: "test123456"})
	}))
	defer ts.Close()
	var channels []Channel
	channels = append(channels, Channel{Integration: "pagerduty", Key: "Your PagerDuty API key goes here", TriggerInterval: "15m", TriggerLimit: 20})
	payload := ViewPayload{Name: "test", Channels: channels}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	result, err := client.CreateAlert(ts.URL, payload)
	assert.Equal(t, nil, err)
	assert.Len(t, result, 10)
}

func TestClient_ViewRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	payload := ViewPayload{Name: "test", Query: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	err := client.UpdateView(ts.URL, "\r\n", payload)
	assert.Error(t, err)
}

func TestClient_ViewResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ this is invalid JSON }"))
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test", Query: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateView(ts.URL, payload)
	assert.Error(t, err)
}

func TestClient_AlertRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	payload := ViewPayload{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	err := client.UpdateAlert(ts.URL, "\r\n", payload)
	assert.Error(t, err)
}

func TestClient_AlertResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ this is invalid JSON }"))
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateAlert(ts.URL, payload)
	assert.Error(t, err)
}

func TestClient_ViewURLError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Invalid URL Error"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test", Query: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateView(ts.URL, payload)
	assert.Equal(t, errors.New("Invalid URL Error"), err)
}

func TestClient_AlertURLError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Invalid URL Error"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateAlert(ts.URL, payload)
	assert.Equal(t, errors.New("Invalid URL Error"), err)
}

func TestClient_ViewResourceNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Resource Not Found"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test", Query: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateView(ts.URL, payload)
	assert.Equal(t, errors.New("Resource Not Found"), err)
}

func TestClient_AlertResourceNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Resource Not Found"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	payload := ViewPayload{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateAlert(ts.URL, payload)
	assert.Equal(t, errors.New("Resource Not Found"), err)
}
