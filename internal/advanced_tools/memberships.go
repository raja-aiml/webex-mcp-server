package advanced_tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)


type ListMembershipsParams struct {
	RoomId      string `json:"roomId,omitempty" query:"roomId"`
	PersonId    string `json:"personId,omitempty" query:"personId"`
	PersonEmail string `json:"personEmail,omitempty" query:"personEmail"`
	Max         int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreateMembershipParams defines the parameters for creating a membership
type CreateMembershipParams struct {
	RoomId      string `json:"roomId" required:"true"`
	PersonId    string `json:"personId,omitempty"`
	PersonEmail string `json:"personEmail,omitempty"`
	IsModerator bool   `json:"isModerator,omitempty"`
}

// UpdateMembershipParams defines the parameters for updating a membership
type UpdateMembershipParams struct {
	MembershipId string `json:"membershipId" required:"true"`
	IsModerator  bool   `json:"isModerator,omitempty"`
}

// NewListMembershipsTool lists room memberships
func NewListMembershipsTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":      StringProperty("List memberships in a room, by room ID."),
		"personId":    StringProperty("List memberships for a person, by person ID."),
		"personEmail": StringProperty("List memberships for a person, by email address."),
		"max":         IntegerProperty("Limit the maximum number of memberships."),
	}

	return tools.NewListTool[ListMembershipsParams](
		"list_memberships",
		"List room memberships",
		"/memberships",
		properties,
	)
}

// NewCreateMembershipTool creates a new membership
func NewCreateMembershipTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":      StringProperty("The room ID."),
		"personId":    StringProperty("The person ID."),
		"personEmail": StringProperty("The email address of the person."),
		"isModerator": BooleanProperty("Whether the person is a room moderator."),
	}

	return tools.NewCreateTool[CreateMembershipParams](
		"create_a_membership",
		"Add someone to a room by Person ID or email address.",
		"/memberships",
		properties,
		[]string{"roomId"},
	)
}

// NewGetMembershipDetailsTool gets membership details
func NewGetMembershipDetailsTool() Tool {
	return tools.NewGetTool(
		"get_membership_details",
		"Get details for a membership by ID.",
		"/memberships",
		"membershipId",
		"The unique identifier for the membership.",
	)
}

// NewUpdateMembershipTool updates a membership
func NewUpdateMembershipTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the membership."),
		"isModerator":  BooleanProperty("Whether the person is a room moderator."),
	}

	return tools.NewUpdateTool[UpdateMembershipParams](
		"update_a_membership",
		"Update properties for a membership by ID.",
		"/memberships",
		"membershipId",
		properties,
		[]string{"membershipId"},
	)
}

// NewDeleteMembershipTool deletes a membership
func NewDeleteMembershipTool() Tool {
	return tools.NewDeleteTool(
		"delete_a_membership",
		"Delete a membership by ID.",
		"/memberships",
		"membershipId",
		"The unique identifier for the membership.",
	)
}
