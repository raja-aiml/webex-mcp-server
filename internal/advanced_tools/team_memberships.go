package advanced_tools

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type ListTeamMembershipsParams struct {
	TeamId string `json:"teamId" required:"true"`
	Max    int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreateTeamMembershipParams defines the parameters for creating a team membership
type CreateTeamMembershipParams struct {
	TeamId      string `json:"teamId" required:"true"`
	PersonId    string `json:"personId,omitempty"`
	PersonEmail string `json:"personEmail,omitempty"`
	IsModerator bool   `json:"isModerator,omitempty"`
}

// UpdateTeamMembershipParams defines the parameters for updating a team membership
type UpdateTeamMembershipParams struct {
	MembershipId string `json:"membershipId" required:"true"`
	IsModerator  bool   `json:"isModerator,omitempty"`
}

// NewListTeamMembershipsTool lists team memberships
func NewListTeamMembershipsTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"teamId": StringProperty("List memberships for a team, by ID."),
		"max":    IntegerProperty("Limit the maximum number of team memberships."),
	}

	return tools.NewListTool[ListTeamMembershipsParams](
		"list_team_memberships",
		"List team memberships for a team.",
		"/team/memberships",
		properties,
		[]string{"teamId"}, // teamId is required
	)
}

// NewCreateTeamMembershipTool creates a new team membership
func NewCreateTeamMembershipTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"teamId":      StringProperty("The team ID."),
		"personId":    StringProperty("The person ID."),
		"personEmail": StringProperty("The email address of the person."),
		"isModerator": BooleanProperty("Whether the person is a team moderator."),
	}

	return tools.NewCreateTool[CreateTeamMembershipParams](
		"create_a_team_membership",
		"Add someone to a team by Person ID or email address.",
		"/team/memberships",
		properties,
		[]string{"teamId"},
	)
}

// NewGetTeamMembershipDetailsTool gets team membership details
func NewGetTeamMembershipDetailsTool() Tool {
	return tools.NewGetTool(
		"get_team_membership_details",
		"Get details for a team membership by ID.",
		"/team/memberships",
		"membershipId",
		"The unique identifier for the team membership.",
	)
}

// NewUpdateTeamMembershipTool updates a team membership
func NewUpdateTeamMembershipTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the team membership."),
		"isModerator":  BooleanProperty("Whether the person is a team moderator."),
	}

	return tools.NewUpdateTool[UpdateTeamMembershipParams](
		"update_a_team_membership",
		"Update a team membership by ID.",
		"/team/memberships",
		"membershipId",
		properties,
		[]string{"membershipId"},
	)
}

// NewDeleteTeamMembershipTool deletes a team membership
func NewDeleteTeamMembershipTool() Tool {
	return tools.NewDeleteTool(
		"delete_a_team_membership",
		"Delete a team membership by ID.",
		"/team/memberships",
		"membershipId",
		"The unique identifier for the team membership.",
	)
}
