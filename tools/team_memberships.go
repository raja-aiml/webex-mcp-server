package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Team Membership Tools - CRUD operations for team memberships

// ListTeamMembershipsTool lists team memberships
type ListTeamMembershipsTool struct {
	ToolBase
}

func NewListTeamMembershipsTool() *ListTeamMembershipsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"teamId": StringProperty("List memberships for a team, by team ID."),
		"max":    IntegerProperty("Limit the maximum number of team memberships."),
	}, []string{"teamId"})

	return &ListTeamMembershipsTool{
		ToolBase: NewToolBase("list_team_memberships", "List team memberships", schema),
	}
}

func (t *ListTeamMembershipsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		TeamId string `json:"teamId"`
		Max    int    `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.TeamId == "" {
		return nil, fmt.Errorf("teamId is required")
	}

	queryParams := map[string]string{"teamId": params.TeamId}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/team/memberships", queryParams)
}

func (t *ListTeamMembershipsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// CreateTeamMembershipTool adds someone to a team
type CreateTeamMembershipTool struct {
	ToolBase
}

func NewCreateTeamMembershipTool() *CreateTeamMembershipTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"teamId":      StringProperty("The team ID."),
		"personId":    StringProperty("The person ID."),
		"personEmail": StringProperty("The email address of the person."),
		"isModerator": BooleanProperty("Whether the person is a team moderator."),
	}, []string{"teamId"})

	return &CreateTeamMembershipTool{
		ToolBase: NewToolBase("create_a_team_membership", "Add someone to a team", schema),
	}
}

func (t *CreateTeamMembershipTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params["teamId"] == nil || params["teamId"] == "" {
		return nil, fmt.Errorf("teamId is required")
	}

	// Validate that either personId or personEmail is provided
	if params["personId"] == nil && params["personEmail"] == nil {
		return nil, fmt.Errorf("either personId or personEmail is required")
	}

	return t.client.Post("/team/memberships", params)
}

func (t *CreateTeamMembershipTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// GetTeamMembershipDetailsTool gets team membership details
type GetTeamMembershipDetailsTool struct {
	ToolBase
}

func NewGetTeamMembershipDetailsTool() *GetTeamMembershipDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the team membership."),
	}, []string{"membershipId"})

	return &GetTeamMembershipDetailsTool{
		ToolBase: NewToolBase("get_team_membership_details", "Get team membership details", schema),
	}
}

func (t *GetTeamMembershipDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		MembershipId string `json:"membershipId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.MembershipId == "" {
		return nil, fmt.Errorf("membershipId is required")
	}

	return t.client.Get(fmt.Sprintf("/team/memberships/%s", params.MembershipId), nil)
}

func (t *GetTeamMembershipDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// UpdateTeamMembershipTool updates a team membership
type UpdateTeamMembershipTool struct {
	ToolBase
}

func NewUpdateTeamMembershipTool() *UpdateTeamMembershipTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the team membership."),
		"isModerator":  BooleanProperty("Whether the person is a team moderator."),
	}, []string{"membershipId"})

	return &UpdateTeamMembershipTool{
		ToolBase: NewToolBase("update_a_team_membership", "Update a team membership", schema),
	}
}

func (t *UpdateTeamMembershipTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	membershipId, ok := params["membershipId"].(string)
	if !ok || membershipId == "" {
		return nil, fmt.Errorf("membershipId is required")
	}

	// Remove membershipId from params as it's in the URL
	delete(params, "membershipId")

	// At least one field to update must be specified
	if len(params) == 0 {
		return nil, fmt.Errorf("at least one field to update must be specified")
	}

	return t.client.Put(fmt.Sprintf("/team/memberships/%s", membershipId), params)
}

func (t *UpdateTeamMembershipTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// DeleteTeamMembershipTool removes someone from a team
type DeleteTeamMembershipTool struct {
	ToolBase
}

func NewDeleteTeamMembershipTool() *DeleteTeamMembershipTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the team membership."),
	}, []string{"membershipId"})

	return &DeleteTeamMembershipTool{
		ToolBase: NewToolBase("delete_a_team_membership", "Delete a team membership", schema),
	}
}

func (t *DeleteTeamMembershipTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		MembershipId string `json:"membershipId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.MembershipId == "" {
		return nil, fmt.Errorf("membershipId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/team/memberships/%s", params.MembershipId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true, "message": "Team membership deleted"}, nil
}

func (t *DeleteTeamMembershipTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}