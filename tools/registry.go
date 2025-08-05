package tools

import (
	"encoding/json"
	"fmt"
)

type Tool interface {
	Name() string
	Description() string
	GetInputSchema() interface{}
	Execute(args json.RawMessage) (interface{}, error)
	ExecuteWithMap(args map[string]interface{}) (interface{}, error)
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(tool Tool) error {
	if _, exists := r.tools[tool.Name()]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name())
	}
	r.tools[tool.Name()] = tool
	return nil
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

func (r *Registry) GetTools() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

func LoadTools() (*Registry, error) {
	registry := NewRegistry()

	// Use plugin architecture for extensibility
	// This implements Open/Closed Principle - open for extension, closed for modification
	manager := NewPluginManager()
	
	// Load all available plugins
	// To add new tools, create a new plugin and register it here
	// No need to modify this function anymore
	loadDefaultPlugins(manager)
	
	// Load plugins into registry
	if err := manager.LoadPlugins(registry); err != nil {
		return nil, err
	}

	return registry, nil
}

// loadDefaultPlugins loads all default tool plugins
// This function can be moved to a separate configuration file
func loadDefaultPlugins(manager *PluginManager) {
	// Import plugins package to access plugin implementations
	// In a real implementation, this could use dynamic loading
	// For now, we'll use compile-time registration
	
	// This is a temporary implementation until we move plugins to their own package
	// The actual implementation would be:
	// plugins.LoadAllPlugins(manager)
	
	// For now, keep the old implementation to maintain compatibility
	tempRegistry := NewRegistry()
	oldLoadTools(tempRegistry)
	
	// Wrap old tools in a compatibility plugin
	manager.RegisterPlugin(&legacyPlugin{registry: tempRegistry})
}

// legacyPlugin wraps the old tool loading for compatibility
type legacyPlugin struct {
	registry *Registry
}

func (p *legacyPlugin) Name() string    { return "legacy" }
func (p *legacyPlugin) Version() string { return "1.0.0" }
func (p *legacyPlugin) Register(registry *Registry) error {
	for _, tool := range p.registry.GetTools() {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// oldLoadTools contains the previous implementation for compatibility
func oldLoadTools(registry *Registry) {
	// Register messaging tools
	messagingTools := []Tool{
		NewListMessagesTool(),
		NewCreateMessageTool(),
		NewDeleteMessageTool(),
		NewEditMessageTool(),
		NewGetMessageDetailsTool(),
		NewListDirectMessagesTool(),
	}

	// Register room tools
	roomTools := []Tool{
		NewListRoomsTool(),
		NewCreateRoomTool(),
		NewGetRoomDetailsTool(),
		NewUpdateRoomTool(),
		NewDeleteRoomTool(),
		NewGetRoomMeetingDetailsTool(),
	}

	// Register people tools
	peopleTools := []Tool{
		NewListPeopleTool(),
		NewCreatePersonTool(),
		NewGetPersonDetailsTool(),
		NewUpdatePersonTool(),
		NewDeletePersonTool(),
		NewGetMyOwnDetailsTool(),
	}

	// Register membership tools
	membershipTools := []Tool{
		NewListMembershipsTool(),
		NewCreateMembershipTool(),
		NewGetMembershipDetailsTool(),
		NewUpdateMembershipTool(),
		NewDeleteMembershipTool(),
	}

	// Register team tools
	teamTools := []Tool{
		NewListTeamsTool(),
		NewCreateTeamTool(),
		NewGetTeamDetailsTool(),
		NewUpdateTeamTool(),
		NewDeleteTeamTool(),
	}

	// Register team membership tools
	teamMembershipTools := []Tool{
		NewListTeamMembershipsTool(),
		NewCreateTeamMembershipTool(),
		NewGetTeamMembershipDetailsTool(),
		NewUpdateTeamMembershipTool(),
		NewDeleteTeamMembershipTool(),
	}

	// Register webhook tools
	webhookTools := []Tool{
		NewListWebhooksTool(),
		NewCreateWebhookTool(),
		NewGetWebhookDetailsTool(),
		NewUpdateWebhookTool(),
		NewDeleteWebhookTool(),
	}

	// Register room tab tools
	roomTabTools := []Tool{
		NewListRoomTabsTool(),
		NewCreateRoomTabTool(),
		NewGetRoomTabDetailsTool(),
		NewUpdateRoomTabTool(),
		NewDeleteRoomTabTool(),
	}

	// Register attachment and event tools
	attachmentTools := []Tool{
		NewCreateAttachmentActionTool(),
		NewGetAttachmentActionDetailsTool(),
	}

	eventTools := []Tool{
		NewListEventsTool(),
		NewGetEventDetailsTool(),
	}

	// Register ECM tools
	ecmTools := []Tool{
		NewCreateECMFolderConfigurationTool(),
		NewGetECMFolderDetailsTool(),
		NewListECMFolderTool(),
		NewUpdateECMLinkedFolderTool(),
		NewUnlinkECMLinkedFolderTool(),
	}

	// Register all tools
	allTools := append(messagingTools, roomTools...)
	allTools = append(allTools, peopleTools...)
	allTools = append(allTools, membershipTools...)
	allTools = append(allTools, teamTools...)
	allTools = append(allTools, teamMembershipTools...)
	allTools = append(allTools, webhookTools...)
	allTools = append(allTools, roomTabTools...)
	allTools = append(allTools, attachmentTools...)
	allTools = append(allTools, eventTools...)
	allTools = append(allTools, ecmTools...)

	for _, tool := range allTools {
		registry.Register(tool)
	}
}