package logdna

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
)

func TestRequestTypes_iterateIntegrationType(t *testing.T) {
	assert := assert.New(t)

	t.Run("Inserts a Diagnostics error for an unknown integration type", func(t *testing.T) {
		var diags diag.Diagnostics
		nonEmptyEntry := []interface{}{
			map[string]interface{}{
				"nah": "will not work",
			},
		}

		channelRequests := iterateIntegrationType(nonEmptyEntry, "NOPE", &diags)

		assert.Empty(*channelRequests, "Nothing was returned")
		assert.Len(diags, 1, "There was 1 error")
		assert.True(diags.HasError(), "The message is of type `Error`")

		err := diags[0]
		assert.Equal("Cannot format integration channel for outbound request", err.Summary, "Summary")
		assert.Equal("Unrecognized integration: NOPE", err.Detail, "Detail")
	})
}
