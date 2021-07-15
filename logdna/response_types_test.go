package logdna

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
)

func TestResponseTypes_mapAllChannelsToSchema(t *testing.T) {
	assert := assert.New(t)

	t.Run("Inserts a Diagnostics error for an unknown integration type in the response", func(t *testing.T) {
		channels := []channelResponse{
			{Integration: "NOPE"},
		}
		channelIntegrations, diags := mapAllChannelsToSchema("view", &channels)

		expected := map[string][]interface{}{
			EMAIL:     make([]interface{}, 0),
			PAGERDUTY: make([]interface{}, 0),
			WEBHOOK:   make([]interface{}, 0),
		}
		assert.Equal(expected, channelIntegrations, "Nothing was returned")
		assert.Len(*diags, 1, "There was 1 diags error")

		err := (*diags)[0]
		assert.Equal(diag.Warning, err.Severity, "The level is Warning")
		assert.Equal("The remote view resource contains an unsupported integration: NOPE", err.Summary, "Summary")
		assert.Equal("NOPE integration ignored since it does not map to the schema", err.Detail, "Detail")
	})
}

func TestResponseTypes_appendError(t *testing.T) {
	assert := assert.New(t)

	t.Run("appendError puts an error on the diags reference", func(t *testing.T) {
		var diags diag.Diagnostics

		err := errors.New("Some Error")

		appendError(err, &diags)

		assert.Len(diags, 1, "There was 1 diags error")

		result := diags[0]
		assert.Equal(diag.Error, result.Severity, "The level is Error")
		assert.Equal("There was a problem setting the schema", result.Summary, "Summary")
		assert.Equal("Some Error", result.Detail, "Detail")
	})
}
