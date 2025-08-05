package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Event Tools - Read operations for events

// ListEventsTool lists events
type ListEventsTool struct {
	ToolBase
}

func NewListEventsTool() *ListEventsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"resource": StringProperty("Limit events to a specific resource type."),
		"type":     StringProperty("Limit events to a specific event type."),
		"actorId":  StringProperty("Limit events to those performed by a specific person."),
		"from":     StringProperty("Limit events to those that occurred after this date and time."),
		"to":       StringProperty("Limit events to those that occurred before this date and time."),
		"max":      IntegerProperty("Limit the maximum number of events in the response."),
	}, []string{})

	return &ListEventsTool{
		ToolBase: NewToolBase("list_events", "List events", schema),
	}
}

func (t *ListEventsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		Resource string `json:"resource,omitempty"`
		Type     string `json:"type,omitempty"`
		ActorId  string `json:"actorId,omitempty"`
		From     string `json:"from,omitempty"`
		To       string `json:"to,omitempty"`
		Max      int    `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	queryParams := make(map[string]string)
	if params.Resource != "" {
		queryParams["resource"] = params.Resource
	}
	if params.Type != "" {
		queryParams["type"] = params.Type
	}
	if params.ActorId != "" {
		queryParams["actorId"] = params.ActorId
	}
	if params.From != "" {
		queryParams["from"] = params.From
	}
	if params.To != "" {
		queryParams["to"] = params.To
	}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/events", queryParams)
}

func (t *ListEventsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// GetEventDetailsTool gets event details
type GetEventDetailsTool struct {
	ToolBase
}

func NewGetEventDetailsTool() *GetEventDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"eventId": StringProperty("The unique identifier for the event."),
	}, []string{"eventId"})

	return &GetEventDetailsTool{
		ToolBase: NewToolBase("get_event_details", "Get event details", schema),
	}
}

func (t *GetEventDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		EventId string `json:"eventId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.EventId == "" {
		return nil, fmt.Errorf("eventId is required")
	}

	return t.client.Get(fmt.Sprintf("/events/%s", params.EventId), nil)
}

func (t *GetEventDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}