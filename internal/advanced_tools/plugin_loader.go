package advanced_tools

import (
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

// LoadAdvancedPlugins loads all advanced tool plugins
// These are tools not required for basic conversation functionality
func LoadAdvancedPlugins(manager *tools.PluginManager) {
	// Register advanced plugins
	manager.RegisterPlugin(&advancedRoomsPlugin{})
	manager.RegisterPlugin(&advancedPeoplePlugin{})
	manager.RegisterPlugin(&advancedMembershipPlugin{})
	manager.RegisterPlugin(&advancedTeamsPlugin{})
	manager.RegisterPlugin(&advancedMiscPlugin{})
}

// advancedRoomsPlugin provides advanced room management tools
type advancedRoomsPlugin struct{}

func (p *advancedRoomsPlugin) Name() string    { return "advanced-rooms" }
func (p *advancedRoomsPlugin) Version() string { return "1.0.0" }

func (p *advancedRoomsPlugin) Register(registry *tools.Registry) error {
	toolList := []tools.Tool{
		NewCreateRoomTool(),
		NewGetRoomDetailsTool(),
		NewUpdateRoomTool(),
		NewDeleteRoomTool(),
		NewGetRoomMeetingDetailsTool(),
	}

	for _, tool := range toolList {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// advancedPeoplePlugin provides advanced people management tools
type advancedPeoplePlugin struct{}

func (p *advancedPeoplePlugin) Name() string    { return "advanced-people" }
func (p *advancedPeoplePlugin) Version() string { return "1.0.0" }

func (p *advancedPeoplePlugin) Register(registry *tools.Registry) error {
	toolList := []tools.Tool{
		NewListPeopleTool(),
		NewCreatePersonTool(),
		NewGetPersonDetailsTool(),
		NewUpdatePersonTool(),
		NewDeletePersonTool(),
	}

	for _, tool := range toolList {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// advancedMembershipPlugin provides membership and team membership tools
type advancedMembershipPlugin struct{}

func (p *advancedMembershipPlugin) Name() string    { return "advanced-membership" }
func (p *advancedMembershipPlugin) Version() string { return "1.0.0" }

func (p *advancedMembershipPlugin) Register(registry *tools.Registry) error {
	toolList := []tools.Tool{
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

	for _, tool := range toolList {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// advancedTeamsPlugin provides team management tools
type advancedTeamsPlugin struct{}

func (p *advancedTeamsPlugin) Name() string    { return "advanced-teams" }
func (p *advancedTeamsPlugin) Version() string { return "1.0.0" }

func (p *advancedTeamsPlugin) Register(registry *tools.Registry) error {
	toolList := []tools.Tool{
		NewListTeamsTool(),
		NewCreateTeamTool(),
		NewGetTeamDetailsTool(),
		NewUpdateTeamTool(),
		NewDeleteTeamTool(),
	}

	for _, tool := range toolList {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// advancedMiscPlugin provides miscellaneous advanced tools
type advancedMiscPlugin struct{}

func (p *advancedMiscPlugin) Name() string    { return "advanced-misc" }
func (p *advancedMiscPlugin) Version() string { return "1.0.0" }

func (p *advancedMiscPlugin) Register(registry *tools.Registry) error {
	toolList := []tools.Tool{
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
	}

	for _, tool := range toolList {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}
