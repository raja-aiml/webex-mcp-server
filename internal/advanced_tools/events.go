package advanced_tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type ListEventsParams struct {
	Resource string `json:"resource,omitempty" query:"resource"`
	Type     string `json:"type,omitempty" query:"type"`
	ActorId  string `json:"actorId,omitempty" query:"actorId"`
	From     string `json:"from,omitempty" query:"from"`
	To       string `json:"to,omitempty" query:"to"`
	Max      int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// NewListEventsTool lists events
func NewListEventsTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"resource": StringProperty("List events related to this resource. Possible values: messages, memberships, etc."),
		"type":     StringProperty("List events of this type. Possible values: created, updated, deleted."),
		"actorId":  StringProperty("List events performed by this person, by ID."),
		"from":     StringProperty("List events which occurred after this date and time (ISO8601 format)."),
		"to":       StringProperty("List events which occurred before this date and time (ISO8601 format)."),
		"max":      IntegerProperty("Limit the maximum number of events in the response."),
	}

	return tools.NewListTool[ListEventsParams](
		"list_events",
		"List events in your organization.",
		"/events",
		properties,
		[]string{}, // No required fields
	)
}

// NewGetEventDetailsTool gets event details
func NewGetEventDetailsTool() Tool {
	return tools.NewGetTool(
		"get_event_details",
		"Get details for an event by ID.",
		"/events",
		"eventId",
		"The unique identifier for the event.",
	)
}
