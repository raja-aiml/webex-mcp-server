package advanced_tools

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type ListTeamsParams struct {
	Max int `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreateTeamParams defines the parameters for creating a team
type CreateTeamParams struct {
	Name        string `json:"name" required:"true"`
	Description string `json:"description,omitempty"`
}

// UpdateTeamParams defines the parameters for updating a team
type UpdateTeamParams struct {
	TeamId      string `json:"teamId" required:"true"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewListTeamsTool lists teams
func NewListTeamsTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"max": IntegerProperty("Limit the maximum number of teams in the response."),
	}

	return tools.NewListTool[ListTeamsParams](
		"list_teams",
		"List teams to which the authenticated user belongs.",
		"/teams",
		properties,
		[]string{}, // No required fields
	)
}

// NewCreateTeamTool creates a new team
func NewCreateTeamTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"name":        StringProperty("A user-friendly name for the team."),
		"description": StringProperty("The description of the team."),
	}

	return tools.NewCreateTool[CreateTeamParams](
		"create_a_team",
		"Create a new team.",
		"/teams",
		properties,
		[]string{"name"},
	)
}

// NewGetTeamDetailsTool gets team details
func NewGetTeamDetailsTool() Tool {
	return tools.NewGetTool(
		"get_team_details",
		"Get details for a team by ID.",
		"/teams",
		"teamId",
		"The unique identifier for the team.",
	)
}

// NewUpdateTeamTool updates a team
func NewUpdateTeamTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"teamId":      StringProperty("The unique identifier for the team."),
		"name":        StringProperty("A user-friendly name for the team."),
		"description": StringProperty("The description of the team."),
	}

	return tools.NewUpdateTool[UpdateTeamParams](
		"update_a_team",
		"Update details for a team by ID.",
		"/teams",
		"teamId",
		properties,
		[]string{"teamId"},
	)
}

// NewDeleteTeamTool deletes a team
func NewDeleteTeamTool() Tool {
	return tools.NewDeleteTool(
		"delete_a_team",
		"Delete a team.",
		"/teams",
		"teamId",
		"The unique identifier for the team.",
	)
}
