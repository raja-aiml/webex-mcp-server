package plugins

import "github.com/raja-aiml/webex-mcp-server-go/tools"

// RoomsPlugin provides room-related tools
type RoomsPlugin struct{}

func NewRoomsPlugin() tools.ToolPlugin {
	return &RoomsPlugin{}
}

func (p *RoomsPlugin) Name() string    { return "rooms" }
func (p *RoomsPlugin) Version() string { return "1.0.0" }

func (p *RoomsPlugin) Register(registry *tools.Registry) error {
	roomTools := []tools.Tool{
		tools.NewListRoomsTool(),
		tools.NewCreateRoomTool(),
		tools.NewGetRoomDetailsTool(),
		tools.NewUpdateRoomTool(),
		tools.NewDeleteRoomTool(),
		tools.NewGetRoomMeetingDetailsTool(),
	}

	for _, tool := range roomTools {
		if err := registry.Register(tool); err != nil {
			return err
		}
	}
	return nil
}