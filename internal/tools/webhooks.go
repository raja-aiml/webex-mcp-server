package tools

import (
	"github.com/google/jsonschema-go/jsonschema"
)

// ListWebhooksParams defines the parameters for listing webhooks
type ListWebhooksParams struct {
	Max int `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreateWebhookParams defines the parameters for creating a webhook
type CreateWebhookParams struct {
	Name      string `json:"name" required:"true"`
	TargetUrl string `json:"targetUrl" required:"true"`
	Resource  string `json:"resource" required:"true"`
	Event     string `json:"event" required:"true"`
	Filter    string `json:"filter,omitempty"`
	Secret    string `json:"secret,omitempty"`
}

// UpdateWebhookParams defines the parameters for updating a webhook
type UpdateWebhookParams struct {
	WebhookId string `json:"webhookId" required:"true"`
	Name      string `json:"name,omitempty"`
	TargetUrl string `json:"targetUrl,omitempty"`
	Secret    string `json:"secret,omitempty"`
	Status    string `json:"status,omitempty"`
}

// NewListWebhooksTool lists webhooks
func NewListWebhooksTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"max": IntegerProperty("Limit the maximum number of webhooks in the response."),
	}

	return NewListTool[ListWebhooksParams](
		"list_webhooks",
		"List all of your webhooks.",
		"/webhooks",
		properties,
		[]string{}, // No required fields for listing webhooks
	)
}

// NewCreateWebhookTool creates a new webhook
func NewCreateWebhookTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"name":      StringProperty("A user-friendly name for the webhook."),
		"targetUrl": StringProperty("The URL that receives POST requests for each event."),
		"resource":  StringProperty("The resource type for the webhook. Possible values: messages, memberships, etc."),
		"event":     StringProperty("The event type for the webhook. Possible values: created, updated, deleted."),
		"filter":    StringProperty("The filter that defines the webhook scope."),
		"secret":    StringProperty("The secret used to generate payload signature."),
	}

	return NewCreateTool[CreateWebhookParams](
		"create_a_webhook",
		"Create a webhook.",
		"/webhooks",
		properties,
		[]string{"name", "targetUrl", "resource", "event"},
	)
}

// NewGetWebhookDetailsTool gets webhook details
func NewGetWebhookDetailsTool() Tool {
	return NewGetTool(
		"get_webhook_details",
		"Get details for a webhook by ID.",
		"/webhooks",
		"webhookId",
		"The unique identifier for the webhook.",
	)
}

// NewUpdateWebhookTool updates a webhook
func NewUpdateWebhookTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"webhookId": StringProperty("The unique identifier for the webhook."),
		"name":      StringProperty("A user-friendly name for the webhook."),
		"targetUrl": StringProperty("The URL that receives POST requests for each event."),
		"secret":    StringProperty("The secret used to generate payload signature."),
		"status":    StringProperty("The status of the webhook. Use 'active' to reactivate a disabled webhook."),
	}

	return NewUpdateTool[UpdateWebhookParams](
		"update_a_webhook",
		"Update a webhook by ID.",
		"/webhooks",
		"webhookId",
		properties,
		[]string{"webhookId"},
	)
}

// NewDeleteWebhookTool deletes a webhook
func NewDeleteWebhookTool() Tool {
	return NewDeleteTool(
		"delete_a_webhook",
		"Delete a webhook by ID.",
		"/webhooks",
		"webhookId",
		"The unique identifier for the webhook.",
	)
}
