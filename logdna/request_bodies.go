package logdna

// This separation of concerns between request and response bodies is only due
// to inconsistencies in the API data types returned by the PUT versus the ones
// returned by the GET. In a perfect world, they would use the same types.

type ViewRequest struct {
	Apps     []string  `json:"apps,omitempty"`
	Category []string  `json:"category,omitempty"`
	Channels []ChannelRequest `json:"channels,omitempty"`
	Hosts    []string  `json:"hosts,omitempty"`
	Levels   []string  `json:"levels,omitempty"`
	Name     string    `json:"name,omitempty"`
	Query    string    `json:"query,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
}

type ChannelRequest struct {
	BodyTemplate    map[string]interface{} `json:"bodyTemplate,omitempty"`
	Emails          []string               `json:"emails,omitempty"`
	Headers         map[string]string      `json:"headers,omitempty"`
	Immediate       string                 `json:"immediate,omitempty"`
	Integration     string                 `json:"integration,omitempty"`
	Key             string                 `json:"key,omitempty"`
	Method          string                 `json:"method,omitempty"`
	Operator        string                 `json:"operator,omitempty"`
	Terminal        string                 `json:"terminal,omitempty"`
	TriggerInterval string                 `json:"triggerinterval,omitempty"`
	TriggerLimit    int                    `json:"triggerlimit,omitempty"`
	Timezone        string                 `json:"timezone,omitempty"`
	URL             string                 `json:"url,omitempty"`
}
