package tools

// LoadDefaultPlugins loads all available tool plugins
// This function registers all the default plugins with the plugin manager
func LoadDefaultPlugins(manager *PluginManager) {
	// Register messaging plugin
	manager.RegisterPlugin(&messagingPlugin{})

	// Register rooms plugin
	manager.RegisterPlugin(&roomsPlugin{})

	// Register people plugin
	manager.RegisterPlugin(&peoplePlugin{})

	// Register membership plugin
	manager.RegisterPlugin(&membershipPlugin{})

	// Register teams plugin
	manager.RegisterPlugin(&teamsPlugin{})

	// Register webhooks plugin
	manager.RegisterPlugin(&webhooksPlugin{})

	// Register miscellaneous tools plugin
	manager.RegisterPlugin(&miscPlugin{})
}

// messagingPlugin provides messaging-related tools
type messagingPlugin struct{}

func (p *messagingPlugin) Name() string    { return "messaging" }
func (p *messagingPlugin) Version() string { return "1.0.0" }

func (p *messagingPlugin) Register(registry *Registry) error {
	tools := []Tool{
		// Original tools
		NewListMessagesTool(),
		NewCreateMessageTool(),
		NewDeleteMessageTool(),
		NewEditMessageTool(),
		NewGetMessageDetailsTool(),
		NewListDirectMessagesTool(),

		// Generic implementations (examples)
		// Uncomment to test generic versions
		// NewListMessagesToolGeneric(),
		// NewCreateMessageToolGeneric(),
		// NewDeleteMessageToolGeneric(),
		// NewGetMessageDetailsToolGeneric(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// roomsPlugin provides room-related tools
type roomsPlugin struct{}

func (p *roomsPlugin) Name() string    { return "rooms" }
func (p *roomsPlugin) Version() string { return "1.0.0" }

func (p *roomsPlugin) Register(registry *Registry) error {
	tools := []Tool{
		NewListRoomsTool(),
		NewCreateRoomTool(),
		NewGetRoomDetailsTool(),
		NewUpdateRoomTool(),
		NewDeleteRoomTool(),
		NewGetRoomMeetingDetailsTool(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// peoplePlugin provides people-related tools
type peoplePlugin struct{}

func (p *peoplePlugin) Name() string    { return "people" }
func (p *peoplePlugin) Version() string { return "1.0.0" }

func (p *peoplePlugin) Register(registry *Registry) error {
	tools := []Tool{
		NewListPeopleTool(),
		NewCreatePersonTool(),
		NewGetPersonDetailsTool(),
		NewUpdatePersonTool(),
		NewDeletePersonTool(),
		NewGetMyOwnDetailsTool(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// membershipPlugin provides membership-related tools
type membershipPlugin struct{}

func (p *membershipPlugin) Name() string    { return "membership" }
func (p *membershipPlugin) Version() string { return "1.0.0" }

func (p *membershipPlugin) Register(registry *Registry) error {
	tools := []Tool{
		NewListMembershipsTool(),
		NewCreateMembershipTool(),
		NewGetMembershipDetailsTool(),
		NewUpdateMembershipTool(),
		NewDeleteMembershipTool(),
		NewListTeamMembershipsTool(),
		NewCreateTeamMembershipTool(),
		NewGetTeamMembershipDetailsTool(),
		NewUpdateTeamMembershipTool(),
		NewDeleteTeamMembershipTool(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// teamsPlugin provides team-related tools
type teamsPlugin struct{}

func (p *teamsPlugin) Name() string    { return "teams" }
func (p *teamsPlugin) Version() string { return "1.0.0" }

func (p *teamsPlugin) Register(registry *Registry) error {
	tools := []Tool{
		NewListTeamsTool(),
		NewCreateTeamTool(),
		NewGetTeamDetailsTool(),
		NewUpdateTeamTool(),
		NewDeleteTeamTool(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// webhooksPlugin provides webhook-related tools
type webhooksPlugin struct{}

func (p *webhooksPlugin) Name() string    { return "webhooks" }
func (p *webhooksPlugin) Version() string { return "1.0.0" }

func (p *webhooksPlugin) Register(registry *Registry) error {
	tools := []Tool{
		NewListWebhooksTool(),
		NewCreateWebhookTool(),
		NewGetWebhookDetailsTool(),
		NewUpdateWebhookTool(),
		NewDeleteWebhookTool(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// miscPlugin provides miscellaneous tools
type miscPlugin struct{}

func (p *miscPlugin) Name() string    { return "misc" }
func (p *miscPlugin) Version() string { return "1.0.0" }

func (p *miscPlugin) Register(registry *Registry) error {
	tools := []Tool{
		// Room tabs
		NewListRoomTabsTool(),
		NewCreateRoomTabTool(),
		NewGetRoomTabDetailsTool(),
		NewUpdateRoomTabTool(),
		NewDeleteRoomTabTool(),
		// Attachments
		NewCreateAttachmentActionTool(),
		NewGetAttachmentActionDetailsTool(),
		// Events
		NewListEventsTool(),
		NewGetEventDetailsTool(),
		// ECM
		NewCreateECMFolderConfigurationTool(),
		NewGetECMFolderDetailsTool(),
		NewListECMFolderTool(),
		NewUpdateECMLinkedFolderTool(),
		NewUnlinkECMLinkedFolderTool(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}
