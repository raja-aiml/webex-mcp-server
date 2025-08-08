package advanced_tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type CreateECMFolderConfigurationParams struct {
	RoomId      string `json:"roomId" required:"true"`
	FolderId    string `json:"folderId" required:"true"`
	DisplayName string `json:"displayName,omitempty"`
}

// NewCreateECMFolderConfigurationTool creates an ECM folder configuration
func NewCreateECMFolderConfigurationTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":      StringProperty("A unique identifier for the room."),
		"folderId":    StringProperty("The ECM folder ID."),
		"displayName": StringProperty("A user-friendly name for the ECM folder."),
	}

	return tools.NewCreateTool[CreateECMFolderConfigurationParams](
		"create_an_ecm_folder_configuration",
		"Create an ECM folder configuration",
		"/rooms/linkedFolders",
		properties,
		[]string{"roomId", "folderId"},
	)
}

// NewGetECMFolderDetailsTool gets ECM folder details
func NewGetECMFolderDetailsTool() Tool {
	return tools.NewGetTool(
		"get_ecm_folder_details",
		"Get details for an ECM folder by ID.",
		"/rooms/linkedFolders",
		"folderId",
		"The unique identifier for the ECM folder.",
	)
}
