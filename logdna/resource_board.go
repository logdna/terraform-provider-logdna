package logdna

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBoardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	board := boardRequest{}

	if diags = board.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/board",
		board,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdBoard := boardResponse{}
	err = json.Unmarshal(body, &createdBoard)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

	d.SetId(createdBoard.boardid)
}

func resourceBoard() *schema.Resource {
	return &schema.Resource{}
}
