package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Webhook Tools - CRUD operations for webhooks

// ListWebhooksTool lists all webhooks
type ListWebhooksTool struct {
	ToolBase
}

func NewListWebhooksTool() *ListWebhooksTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"max": IntegerProperty("Limit the maximum number of webhooks."),
	}, []string{})

	return &ListWebhooksTool{
		ToolBase: NewToolBase("list_webhooks", "List all webhooks", schema),
	}
}

func (t *ListWebhooksTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		Max int `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	queryParams := make(map[string]string)
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/webhooks", queryParams)
}

func (t *ListWebhooksTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// CreateWebhookTool creates a webhook
type CreateWebhookTool struct {
	ToolBase
}

func NewCreateWebhookTool() *CreateWebhookTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"name":      StringProperty("A user-friendly name for the webhook."),
		"targetUrl": StringProperty("The URL that receives POST requests for each event."),
		"resource":  StringProperty("The resource type for the webhook."),
		"event":     StringProperty("The event type for the webhook."),
		"filter":    StringProperty("The filter that defines the webhook scope."),
		"secret":    StringProperty("The secret used to generate payload signature."),
	}, []string{"name", "targetUrl", "resource", "event"})

	return &CreateWebhookTool{
		ToolBase: NewToolBase("create_a_webhook", "Create a webhook", schema),
	}
}

func (t *CreateWebhookTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required fields
	required := []string{"name", "targetUrl", "resource", "event"}
	for _, field := range required {
		if params[field] == nil || params[field] == "" {
			return nil, fmt.Errorf("%s is required", field)
		}
	}

	return t.client.Post("/webhooks", params)
}

func (t *CreateWebhookTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// GetWebhookDetailsTool gets webhook details
type GetWebhookDetailsTool struct {
	ToolBase
}

func NewGetWebhookDetailsTool() *GetWebhookDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"webhookId": StringProperty("The unique identifier for the webhook."),
	}, []string{"webhookId"})

	return &GetWebhookDetailsTool{
		ToolBase: NewToolBase("get_webhook_details", "Get webhook details", schema),
	}
}

func (t *GetWebhookDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		WebhookId string `json:"webhookId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.WebhookId == "" {
		return nil, fmt.Errorf("webhookId is required")
	}

	return t.client.Get(fmt.Sprintf("/webhooks/%s", params.WebhookId), nil)
}

func (t *GetWebhookDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// UpdateWebhookTool updates a webhook
type UpdateWebhookTool struct {
	ToolBase
}

func NewUpdateWebhookTool() *UpdateWebhookTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"webhookId": StringProperty("The unique identifier for the webhook."),
		"name":      StringProperty("A user-friendly name for the webhook."),
		"targetUrl": StringProperty("The URL that receives POST requests for each event."),
		"secret":    StringProperty("The secret used to generate payload signature."),
		"status":    StringProperty("The status of the webhook (active or inactive)."),
	}, []string{"webhookId"})

	return &UpdateWebhookTool{
		ToolBase: NewToolBase("update_a_webhook", "Update a webhook", schema),
	}
}

func (t *UpdateWebhookTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	webhookId, ok := params["webhookId"].(string)
	if !ok || webhookId == "" {
		return nil, fmt.Errorf("webhookId is required")
	}

	// Remove webhookId from params as it's in the URL
	delete(params, "webhookId")

	// At least one field to update must be specified
	if len(params) == 0 {
		return nil, fmt.Errorf("at least one field to update must be specified")
	}

	return t.client.Put(fmt.Sprintf("/webhooks/%s", webhookId), params)
}

func (t *UpdateWebhookTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// DeleteWebhookTool deletes a webhook
type DeleteWebhookTool struct {
	ToolBase
}

func NewDeleteWebhookTool() *DeleteWebhookTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"webhookId": StringProperty("The unique identifier for the webhook."),
	}, []string{"webhookId"})

	return &DeleteWebhookTool{
		ToolBase: NewToolBase("delete_a_webhook", "Delete a webhook", schema),
	}
}

func (t *DeleteWebhookTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		WebhookId string `json:"webhookId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.WebhookId == "" {
		return nil, fmt.Errorf("webhookId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/webhooks/%s", params.WebhookId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true, "message": "Webhook deleted"}, nil
}

func (t *DeleteWebhookTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}