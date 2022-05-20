package logdna

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestStreamConfig_expectInvalidURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamConfig(`
					brokers = ["broker-1.example.org:9090"]
					topic = "test-topic"
					user = "test-user"
					password = "test-password"
				`, "http://api.logdna.co"),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"http://api.logdna.co/v1/config/stream\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestStreamConfig_expectInvalidBrokerError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamConfig(`
					brokers = ["broker-1.example.org:9090"]
					topic = "test-topic"
					user = "test-user"
					password = "test-password"
				`, apiHostUrl),
				ExpectError: regexp.MustCompile(`Failed to connect to Kafka broker`),
			},
		},
	})
}

func TestStreamConfig_expectInvalidConfigError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamConfig(`
					brokers = ["broker-1.example.org:9090"]
					topic = ""
					user = ""
					password = ""
				`, apiHostUrl),
				ExpectError: regexp.MustCompile(`\\"topic\\" is not allowed to be empty.*\\"user\\" is not allowed to be empty.*\\"password\\" is not allowed to be empty`),
			},
		},
	})
}

func TestStreamConfig_basic(t *testing.T) {
	assert := assert.New(t)
	brokers := []string{
		"broker-1.example.org:9090",
		"broker-2.example.org:9090",
	}
	topic := "test-topic"
	user := "test-user"
	password := "test-password"

	// This resource requires a valid Kafka broker, so this test runs against a mock server/
	// This TestCase is doing a Create and Update, both of which perform a Read at the
	// end of the routine and a refresh between steps. The test server response is
	// tailored to the request by the request number to return what is expected from the
	// step's config. This is a hacky workaround but avoids the need for real Kafka
	// infrastructure for the validation that occurs in this endpoint.
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		t := topic
		if count > 4 {
			t = "updated"
		}
		err := json.NewEncoder(w).Encode(streamConfig{
			Brokers: brokers,
			Topic:   t,
			User:    user,
			Status:  "active",
		})
		assert.Nil(err, "No errors")
	}))
	defer ts.Close()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamConfig(fmt.Sprintf(`
					brokers = [
						"%s",
						"%s"
					]
					topic = "%s"
					user = "%s"
					password = "%s"
				`, brokers[0], brokers[1], topic, user, password,
				), ts.URL),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("stream_config", "stream"),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "topic", topic),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "user", user),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "status", "active"),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "brokers.#", "2"),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "brokers.0", brokers[0]),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "brokers.1", brokers[1]),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "password", password),
				),
			},
			{
				Config: testStreamConfig(fmt.Sprintf(`
					brokers = [
						"%s",
						"%s"
					]
					topic = "updated"
					user = "%s"
					password = "%s"
				`, brokers[0], brokers[1], user, password,
				), ts.URL),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("stream_config", "stream"),
					resource.TestCheckResourceAttr("logdna_stream_config.stream", "topic", "updated"),
				),
			},
			{
				ResourceName:      "logdna_stream_config.stream",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testStreamConfig(fields string, url string) string {
	uc := ""
	if url != "" {
		uc = fmt.Sprintf(`url = "%s"`, url)
	}

	return fmt.Sprintf(`
		provider "logdna" {
			servicekey = "%s"
			%s
		}
		
		resource "logdna_stream_config" "stream" {
			%s
		}
	`, serviceKey, uc, fields)
}
