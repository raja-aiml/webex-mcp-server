package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Membership Tools - CRUD operations for room memberships

// ListMembershipsTool lists room memberships
type ListMembershipsTool struct {
	ToolBase
}

func NewListMembershipsTool() *ListMembershipsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":      StringProperty("List memberships in a room, by room ID."),
		"personId":    StringProperty("List memberships for a person, by person ID."),
		"personEmail": StringProperty("List memberships for a person, by email address."),
		"max":         IntegerProperty("Limit the maximum number of memberships."),
	}, []string{})

	return &ListMembershipsTool{
		ToolBase: NewToolBase("list_memberships", "List room memberships", schema),
	}
}

func (t *ListMembershipsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomId      string `json:"roomId,omitempty"`
		PersonId    string `json:"personId,omitempty"`
		PersonEmail string `json:"personEmail,omitempty"`
		Max         int    `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	queryParams := make(map[string]string)
	if params.RoomId != "" {
		queryParams["roomId"] = params.RoomId
	}
	if params.PersonId != "" {
		queryParams["personId"] = params.PersonId
	}
	if params.PersonEmail != "" {
		queryParams["personEmail"] = params.PersonEmail
	}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/memberships", queryParams)
}

func (t *ListMembershipsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// CreateMembershipTool adds someone to a room
type CreateMembershipTool struct {
	ToolBase
}

func NewCreateMembershipTool() *CreateMembershipTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":      StringProperty("The room ID."),
		"personId":    StringProperty("The person ID."),
		"personEmail": StringProperty("The email address of the person."),
		"isModerator": BooleanProperty("Whether the person is a room moderator."),
	}, []string{"roomId"})

	return &CreateMembershipTool{
		ToolBase: NewToolBase("create_a_membership", "Add someone to a room by Person ID or email", schema),
	}
}

func (t *CreateMembershipTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params["roomId"] == nil || params["roomId"] == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	// Validate that either personId or personEmail is provided
	if params["personId"] == nil && params["personEmail"] == nil {
		return nil, fmt.Errorf("either personId or personEmail is required")
	}

	return t.client.Post("/memberships", params)
}

func (t *CreateMembershipTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// GetMembershipDetailsTool gets details for a membership
type GetMembershipDetailsTool struct {
	ToolBase
}

func NewGetMembershipDetailsTool() *GetMembershipDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the membership."),
	}, []string{"membershipId"})

	return &GetMembershipDetailsTool{
		ToolBase: NewToolBase("get_membership_details", "Get details for a membership by ID", schema),
	}
}

func (t *GetMembershipDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		MembershipId string `json:"membershipId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.MembershipId == "" {
		return nil, fmt.Errorf("membershipId is required")
	}

	return t.client.Get(fmt.Sprintf("/memberships/%s", params.MembershipId), nil)
}

func (t *GetMembershipDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// UpdateMembershipTool updates membership properties
type UpdateMembershipTool struct {
	ToolBase
}

func NewUpdateMembershipTool() *UpdateMembershipTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"membershipId":  StringProperty("The unique identifier for the membership."),
		"isModerator":   BooleanProperty("Whether the person is a room moderator."),
		"isRoomHidden":  BooleanProperty("Whether the room is hidden in the Webex app."),
	}, []string{"membershipId"})

	return &UpdateMembershipTool{
		ToolBase: NewToolBase("update_a_membership", "Update properties for a membership", schema),
	}
}

func (t *UpdateMembershipTool) Execute(args json.RawMessage) (interface{}, error) {
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

	return t.client.Put(fmt.Sprintf("/memberships/%s", membershipId), params)
}

func (t *UpdateMembershipTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// DeleteMembershipTool removes someone from a room
type DeleteMembershipTool struct {
	ToolBase
}

func NewDeleteMembershipTool() *DeleteMembershipTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"membershipId": StringProperty("The unique identifier for the membership."),
	}, []string{"membershipId"})

	return &DeleteMembershipTool{
		ToolBase: NewToolBase("delete_a_membership", "Delete a membership by ID", schema),
	}
}

func (t *DeleteMembershipTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		MembershipId string `json:"membershipId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.MembershipId == "" {
		return nil, fmt.Errorf("membershipId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/memberships/%s", params.MembershipId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true, "message": "Membership deleted"}, nil
}

func (t *DeleteMembershipTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}