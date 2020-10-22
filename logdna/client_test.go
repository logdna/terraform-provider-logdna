package logdna

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestClient_GetViewURLError(t *testing.T) {
	client := Client{ServiceKey: servicekey, HTTPClient: &http.Client{Timeout: 15 * time.Second}}
	_, err := client.GetView("invalid url", "1234")
	assert.Equal(t, errors.New("Error with view: Get \"invalid%20url/v1/config/view/1234\": unsupported protocol scheme \"\""), err)
}

func TestClient_GetAlertURLError(t *testing.T) {
	client := Client{ServiceKey: servicekey, HTTPClient: &http.Client{Timeout: 15 * time.Second}}
	_, err := client.GetAlert("invalid url", "1234")
	assert.Equal(t, errors.New("Error with alert: Get \"invalid%20url/v1/config/presetalert/1234\": unsupported protocol scheme \"\""), err)
}

func TestClient_GetViewResourceNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Resource Not Found"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.GetView(ts.URL, "1234")
	assert.Equal(t, errors.New("Resource Not Found"), err)
}

func TestClient_GetAlertResourceNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Resource Not Found"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.GetAlert(ts.URL, "1234")
	assert.Equal(t, errors.New("Resource Not Found"), err)
}

func TestClient_GetViewRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.GetView(ts.URL, "\r\n")
	assert.Error(t, err)
}

func TestClient_GetViewResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ this is invalid JSON }"))
	}))
	defer ts.Close()
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.GetView(ts.URL, "1234")
	assert.Error(t, err)
}

func TestClient_GetAlertRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.GetAlert(ts.URL, "\r\n")
	assert.Error(t, err)
}

func TestClient_GetAlertResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ this is invalid JSON }"))
	}))
	defer ts.Close()
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.GetAlert(ts.URL, "1234")
	assert.Error(t, err)
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
