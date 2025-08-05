package plugins

import "github.com/raja-aiml/webex-mcp-server-go/tools"

// PeoplePlugin provides people-related tools
type PeoplePlugin struct{}

func NewPeoplePlugin() tools.ToolPlugin {
	return &PeoplePlugin{}
}

func (p *PeoplePlugin) Name() string    { return "people" }
func (p *PeoplePlugin) Version() string { return "1.0.0" }

func (p *PeoplePlugin) Register(registry *tools.Registry) error {
	peopleTools := []tools.Tool{
		tools.NewListPeopleTool(),
		tools.NewCreatePersonTool(),
		tools.NewGetPersonDetailsTool(),
		tools.NewUpdatePersonTool(),
		tools.NewDeletePersonTool(),
		tools.NewGetMyOwnDetailsTool(),
	}

	for _, tool := range peopleTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}