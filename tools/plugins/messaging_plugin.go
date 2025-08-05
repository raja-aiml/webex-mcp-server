package plugins

import "github.com/raja-aiml/webex-mcp-server-go/tools"

// MessagingPlugin provides messaging-related tools
type MessagingPlugin struct{}

func NewMessagingPlugin() tools.ToolPlugin {
	return &MessagingPlugin{}
}

func (p *MessagingPlugin) Name() string    { return "messaging" }
func (p *MessagingPlugin) Version() string { return "1.0.0" }

func (p *MessagingPlugin) Register(registry *tools.Registry) error {
	messagingTools := []tools.Tool{
		tools.NewListMessagesTool(),
		tools.NewCreateMessageTool(),
		tools.NewDeleteMessageTool(),
		tools.NewEditMessageTool(),
		tools.NewGetMessageDetailsTool(),
		tools.NewListDirectMessagesTool(),
	}

	for _, tool := range messagingTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}
