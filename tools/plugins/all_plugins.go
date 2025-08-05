package plugins

import "github.com/raja-aiml/webex-mcp-server-go/tools"

// LoadAllPlugins registers all available plugins
func LoadAllPlugins(manager *tools.PluginManager) {
	// Register all tool plugins
	manager.RegisterPlugin(NewMessagingPlugin())
	manager.RegisterPlugin(NewRoomsPlugin())
	manager.RegisterPlugin(NewPeoplePlugin())
	manager.RegisterPlugin(NewMembershipPlugin())
	manager.RegisterPlugin(NewTeamsPlugin())
	manager.RegisterPlugin(NewWebhooksPlugin())
	manager.RegisterPlugin(NewMiscPlugin())
}

// MembershipPlugin provides membership-related tools
type MembershipPlugin struct{}

func NewMembershipPlugin() tools.ToolPlugin {
	return &MembershipPlugin{}
}

func (p *MembershipPlugin) Name() string    { return "membership" }
func (p *MembershipPlugin) Version() string { return "1.0.0" }

func (p *MembershipPlugin) Register(registry *tools.Registry) error {
	membershipTools := []tools.Tool{
		tools.NewListMembershipsTool(),
		tools.NewCreateMembershipTool(),
		tools.NewGetMembershipDetailsTool(),
		tools.NewUpdateMembershipTool(),
		tools.NewDeleteMembershipTool(),
		tools.NewListTeamMembershipsTool(),
		tools.NewCreateTeamMembershipTool(),
		tools.NewGetTeamMembershipDetailsTool(),
		tools.NewUpdateTeamMembershipTool(),
		tools.NewDeleteTeamMembershipTool(),
	}

	for _, tool := range membershipTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// TeamsPlugin provides team-related tools
type TeamsPlugin struct{}

func NewTeamsPlugin() tools.ToolPlugin {
	return &TeamsPlugin{}
}

func (p *TeamsPlugin) Name() string    { return "teams" }
func (p *TeamsPlugin) Version() string { return "1.0.0" }

func (p *TeamsPlugin) Register(registry *tools.Registry) error {
	teamTools := []tools.Tool{
		tools.NewListTeamsTool(),
		tools.NewCreateTeamTool(),
		tools.NewGetTeamDetailsTool(),
		tools.NewUpdateTeamTool(),
		tools.NewDeleteTeamTool(),
	}

	for _, tool := range teamTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// WebhooksPlugin provides webhook-related tools
type WebhooksPlugin struct{}

func NewWebhooksPlugin() tools.ToolPlugin {
	return &WebhooksPlugin{}
}

func (p *WebhooksPlugin) Name() string    { return "webhooks" }
func (p *WebhooksPlugin) Version() string { return "1.0.0" }

func (p *WebhooksPlugin) Register(registry *tools.Registry) error {
	webhookTools := []tools.Tool{
		tools.NewListWebhooksTool(),
		tools.NewCreateWebhookTool(),
		tools.NewGetWebhookDetailsTool(),
		tools.NewUpdateWebhookTool(),
		tools.NewDeleteWebhookTool(),
	}

	for _, tool := range webhookTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// MiscPlugin provides miscellaneous tools
type MiscPlugin struct{}

func NewMiscPlugin() tools.ToolPlugin {
	return &MiscPlugin{}
}

func (p *MiscPlugin) Name() string    { return "misc" }
func (p *MiscPlugin) Version() string { return "1.0.0" }

func (p *MiscPlugin) Register(registry *tools.Registry) error {
	miscTools := []tools.Tool{
		// Room tabs
		tools.NewListRoomTabsTool(),
		tools.NewCreateRoomTabTool(),
		tools.NewGetRoomTabDetailsTool(),
		tools.NewUpdateRoomTabTool(),
		tools.NewDeleteRoomTabTool(),
		// Attachments
		tools.NewCreateAttachmentActionTool(),
		tools.NewGetAttachmentActionDetailsTool(),
		// Events
		tools.NewListEventsTool(),
		tools.NewGetEventDetailsTool(),
		// ECM
		tools.NewCreateECMFolderConfigurationTool(),
		tools.NewGetECMFolderDetailsTool(),
		tools.NewListECMFolderTool(),
		tools.NewUpdateECMLinkedFolderTool(),
		tools.NewUnlinkECMLinkedFolderTool(),
	}

	for _, tool := range miscTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}