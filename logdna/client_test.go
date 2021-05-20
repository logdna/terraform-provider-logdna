package logdna

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"strings"

	"github.com/stretchr/testify/assert"
)

const SERVICE_KEY = "abc123"

func TestClient_MakeRequest(t *testing.T) {
  assert := assert.New(t)
  resourceId := "test123456"

	t.Run("Server receives proper method, URL, and headers", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    assert.Equal("GET", r.Method,  "Method is correct")
	    assert.Equal(fmt.Sprintf("/someapi/%s", resourceId), r.URL.String(), "URL is correct")
			key, ok := r.Header["Servicekey"]
	    assert.Equal(true, ok, "servicekey header exists")
	    assert.Equal(1, len(key), "servicekey header is correct")
			key = r.Header["Content-Type"]
	    assert.Equal("application/json", key[0], "content-type header is correct")
	  }))
	  defer ts.Close()

	  client := Client{
	    ServiceKey: SERVICE_KEY,
	    HTTPClient: ts.Client(),
	    ApiUrl: fmt.Sprintf("%s/someapi/%s", ts.URL, resourceId),
	    Method: "GET",
	  }

	  _, err := client.MakeRequest()
		assert.Nil(err, "No errors")
	})

	t.Run("Reads and decodes response from the server", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(ViewResponse{ViewID: "test123456"})
		}))
		defer ts.Close()

		client := Client{
			ServiceKey: SERVICE_KEY,
			HTTPClient: ts.Client(),
			ApiUrl: ts.URL,
			Method: "GET",
		}

		body, err := client.MakeRequest()
		assert.Nil(err, "No errors")
		assert.Equal(
			`{"viewID":"test123456"}`,
			strings.TrimSpace(string(body)),
			"Returned body is correct",
		)
	})

	t.Run("Handles errors on NewRequest", func(t *testing.T) {
		client := Client{
			ServiceKey: SERVICE_KEY,
			HTTPClient: &http.Client{},
			ApiUrl: "/uri/only/will/not/work",
			Method: "GET",
		}

		body, err := client.MakeRequest()
		assert.Nil(body, "No body due to error")
		assert.Error(err, "Expected error")
		assert.Equal(
			true,
			strings.Contains(err.Error(), "unsupported protocol scheme"),
			"Expected error message",
		)
	})

	t.Run("Throws non-200 errors returned by the server", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
		}))
		defer ts.Close()

		client := Client{
			ServiceKey: SERVICE_KEY,
			HTTPClient: ts.Client(),
			ApiUrl: ts.URL,
			Method: "POST",
		}

		_, err := client.MakeRequest()
		assert.Error(err, "Expected error")
		assert.Equal(
			true,
			strings.Contains(err.Error(), "status NOT OK: 400"),
			"Expected error message",
		)
	})
}

func TestClient_Alert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(AlertResponsePayload{PresetID: "test123456"})
	}))
	defer ts.Close()
	var channels []ChannelRequest
	channels = append(channels, ChannelRequest{Integration: "pagerduty", Key: "Your PagerDuty API key goes here", TriggerInterval: "15m", TriggerLimit: 20})
	payload := ViewRequest{Name: "test", Channels: channels}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	result, err := client.CreateAlert(ts.URL, payload)
	assert.Equal(t, nil, err)
	assert.Len(t, result, 10)
}

// func TestClient_ViewResponseError(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write([]byte("{ this is invalid JSON }"))
// 	}))
// 	defer ts.Close()
// 	payload := ViewRequest{Name: "test", Query: "test"}
// 	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
// 	_, err := client.CreateView(ts.URL, payload)
// 	assert.Error(t, err)
// }

func TestClient_AlertRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	payload := ViewRequest{Name: "test"}
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
	payload := ViewRequest{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateAlert(ts.URL, payload)
	assert.Error(t, err)
}

// func TestClient_ViewURLError(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		response := make(map[string]interface{})
// 		response["error"] = "Invalid URL Error"
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 	}))
// 	defer ts.Close()
// 	payload := ViewRequest{Name: "test", Query: "test"}
// 	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
// 	_, err := client.CreateView(ts.URL, payload)
// 	assert.Equal(t, errors.New("Invalid URL Error"), err)
// }

func TestClient_AlertURLError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Invalid URL Error"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	payload := ViewRequest{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateAlert(ts.URL, payload)
	assert.Equal(t, errors.New("Invalid URL Error"), err)
}

// func TestClient_ViewResourceNotFound(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		response := make(map[string]interface{})
// 		response["error"] = "Resource Not Found"
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 	}))
// 	defer ts.Close()
// 	payload := ViewRequest{Name: "test", Query: "test"}
// 	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
// 	_, err := client.CreateView(ts.URL, payload)
// 	assert.Equal(t, errors.New("Resource Not Found"), err)
// }

func TestClient_AlertResourceNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		response["error"] = "Resource Not Found"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()
	payload := ViewRequest{Name: "test"}
	client := Client{ServiceKey: servicekey, HTTPClient: ts.Client()}
	_, err := client.CreateAlert(ts.URL, payload)
	assert.Equal(t, errors.New("Resource Not Found"), err)
}
