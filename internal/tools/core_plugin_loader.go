package tools

// LoadCorePlugins loads only the essential tools for conversation functionality
// Following KISS principle - Keep It Simple, only what's needed for bot conversations
func LoadCorePlugins(manager *PluginManager) {
	// Register only core plugins needed for conversation
	manager.RegisterPlugin(&coreMessagingPlugin{})
	manager.RegisterPlugin(&coreWebhooksPlugin{})
	manager.RegisterPlugin(&coreInfoPlugin{})
}

// coreMessagingPlugin provides only essential messaging tools for conversations
type coreMessagingPlugin struct{}

func (p *coreMessagingPlugin) Name() string    { return "core-messaging" }
func (p *coreMessagingPlugin) Version() string { return "1.0.0" }

func (p *coreMessagingPlugin) Register(registry *Registry) error {
	// YAGNI: Only tools needed for query-response conversations
	tools := []Tool{
		NewListMessagesTool(),  // Read incoming queries
		NewCreateMessageTool(), // Send responses
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// coreWebhooksPlugin provides webhook tools for receiving messages
type coreWebhooksPlugin struct{}

func (p *coreWebhooksPlugin) Name() string    { return "core-webhooks" }
func (p *coreWebhooksPlugin) Version() string { return "1.0.0" }

func (p *coreWebhooksPlugin) Register(registry *Registry) error {
	// Essential for receiving incoming messages
	tools := []Tool{
		NewCreateWebhookTool(), // Set up message reception
		NewListWebhooksTool(),  // Manage webhooks
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// coreInfoPlugin provides basic info tools
type coreInfoPlugin struct{}

func (p *coreInfoPlugin) Name() string    { return "core-info" }
func (p *coreInfoPlugin) Version() string { return "1.0.0" }

func (p *coreInfoPlugin) Register(registry *Registry) error {
	// Minimal tools for bot identity and room context
	tools := []Tool{
		NewListRoomsTool(),       // Find rooms to operate in
		NewGetMyOwnDetailsTool(), // Get bot identity
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}
