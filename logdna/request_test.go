package logdna

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type badClient struct{}

func (fc *badClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("FAKE ERROR calling httpClient.Do")
}

func setHTTPRequest(customReq httpRequest) func(*requestConfig) {
	return func(req *requestConfig) {
		req.httpRequest = customReq
	}
}

func setBodyReader(customReader bodyReader) func(*requestConfig) {
	return func(req *requestConfig) {
		req.bodyReader = customReader
	}
}

func setJSONMarshal(customMarshaller jsonMarshal) func(*requestConfig) {
	return func(req *requestConfig) {
		req.jsonMarshal = customMarshaller
	}
}

func TestRequest_MakeRequest(t *testing.T) {
	assert := assert.New(t)
	pc := providerConfig{serviceKey: "abc123", httpClient: &http.Client{Timeout: 15 * time.Second}}
	resourceID := "test123456"

	t.Run("Server receives proper method, URL, and headers", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal("GET", r.Method, "method is correct")
			assert.Equal(fmt.Sprintf("/someapi/%s", resourceID), r.URL.String(), "URL is correct")
			key, ok := r.Header["Servicekey"]
			assert.Equal(true, ok, "servicekey header exists")
			assert.Equal(1, len(key), "servicekey header is correct")
			key = r.Header["Content-Type"]
			assert.Equal("application/json", key[0], "content-type header is correct")
		}))
		defer ts.Close()

		pc.baseURL = ts.URL

		req := newRequestConfig(
			&pc,
			"GET",
			fmt.Sprintf("/someapi/%s", resourceID),
			nil,
		)

		_, err := req.MakeRequest()
		assert.Nil(err, "No errors")
	})

	t.Run("Reads and decodes response from the server", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(viewResponse{ViewID: "test123456"})
			assert.Nil(err, "No errors")
		}))
		defer ts.Close()

		pc.baseURL = ts.URL

		req := newRequestConfig(
			&pc,
			"GET",
			fmt.Sprintf("/someapi/%s", resourceID),
			nil,
		)

		body, err := req.MakeRequest()
		assert.Nil(err, "No errors")
		assert.Equal(
			`{"viewID":"test123456"}`,
			strings.TrimSpace(string(body)),
			"Returned body is correct",
		)
	})

	t.Run("Successfully marshals a provided body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			postedBody, _ := ioutil.ReadAll(r.Body)
			assert.Equal(
				`{"name":"Test View"}`,
				strings.TrimSpace(string(postedBody)),
				"Body got marshalled and sent correctly",
			)
		}))
		defer ts.Close()

		pc.baseURL = ts.URL

		req := newRequestConfig(
			&pc,
			"POST",
			"/someapi",
			viewRequest{
				Name: "Test View",
			},
		)

		_, err := req.MakeRequest()
		assert.Nil(err, "No errors")
	})

	t.Run("Handles errors when marshalling JSON", func(t *testing.T) {
		const ERROR = "FAKE ERROR during json.Marshal"
		req := newRequestConfig(
			&pc,
			"POST",
			"/will/not/work",
			viewRequest{Name: "NOPE"},
			setJSONMarshal(func(interface{}) ([]byte, error) {
				return nil, errors.New(ERROR)
			}),
		)
		body, err := req.MakeRequest()
		assert.Nil(body, "No body due to error")
		assert.Error(err, "Expected error")
		assert.Equal(
			ERROR,
			err.Error(),
			"Expected error message",
		)
	})

	t.Run("Handles errors when creating a new HTTP request", func(t *testing.T) {
		const ERROR = "FAKE ERROR for http.NewRequest"
		req := newRequestConfig(
			&pc,
			"GET",
			"/will/not/work",
			nil,
			setHTTPRequest(func(string, string, io.Reader) (*http.Request, error) {
				return nil, errors.New(ERROR)
			}),
		)
		body, err := req.MakeRequest()
		assert.Nil(body, "No body due to error")
		assert.Error(err, "Expected error")
		assert.Equal(
			ERROR,
			err.Error(),
			"Expected error message",
		)
	})

	t.Run("Handles errors during the HTTP request", func(t *testing.T) {
		req := newRequestConfig(
			&pc,
			"GET",
			"/will/not/work",
			nil,
			func(req *requestConfig) {
				req.httpClient = &badClient{}
			},
		)

		body, err := req.MakeRequest()
		assert.Nil(body, "No body due to error")
		assert.Error(err, "Expected error")
		assert.Equal(
			true,
			strings.Contains(err.Error(), "error during HTTP request: FAKE ERROR calling httpClient.Do"),
			"Expected error message",
		)
	})

	t.Run("Throws non-200 errors returned by the server", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
		}))
		defer ts.Close()

		pc.baseURL = ts.URL

		req := newRequestConfig(
			&pc,
			"POST",
			fmt.Sprintf("/someapi/%s", resourceID),
			nil,
		)

		_, err := req.MakeRequest()
		assert.Error(err, "Expected error")
		assert.Equal(
			true,
			strings.Contains(err.Error(), "status 400 NOT OK!"),
			"Expected error message",
		)
	})

	t.Run("Handles errors when creating a new HTTP request", func(t *testing.T) {
		const ERROR = "FAKE ERROR for body reader"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(viewResponse{ViewID: "test123456"})
			assert.Nil(err, "No errors")
		}))
		defer ts.Close()

		pc.baseURL = ts.URL
		req := newRequestConfig(
			&pc,
			"GET",
			fmt.Sprintf("/someapi/%s", resourceID),
			nil,
			setBodyReader(func(io.Reader) ([]byte, error) {
				return nil, errors.New(ERROR)
			}),
		)
		body, err := req.MakeRequest()
		assert.Nil(body, "No body due to error")
		assert.Error(err, "Expected error")
		assert.Equal(
			true,
			strings.Contains(err.Error(), "error parsing HTTP response: "+ERROR),
			"Expected error message",
		)
	})
}
