package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Team Tools - CRUD operations for teams

// ListTeamsTool lists teams
type ListTeamsTool struct {
	ToolBase
}

func NewListTeamsTool() *ListTeamsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"max": IntegerProperty("Limit the maximum number of teams."),
	}, []string{})

	return &ListTeamsTool{
		ToolBase: NewToolBase("list_teams", "List teams", schema),
	}
}

func (t *ListTeamsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		Max int `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	queryParams := make(map[string]string)
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/teams", queryParams)
}

func (t *ListTeamsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// CreateTeamTool creates a new team
type CreateTeamTool struct {
	ToolBase
}

func NewCreateTeamTool() *CreateTeamTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"name": StringProperty("The name of the team."),
	}, []string{"name"})

	return &CreateTeamTool{
		ToolBase: NewToolBase("create_a_team", "Create a team", schema),
	}
}

func (t *CreateTeamTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params["name"] == nil || params["name"] == "" {
		return nil, fmt.Errorf("name is required")
	}

	return t.client.Post("/teams", params)
}

func (t *CreateTeamTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetTeamDetailsTool gets team details
type GetTeamDetailsTool struct {
	ToolBase
}

func NewGetTeamDetailsTool() *GetTeamDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"teamId": StringProperty("The unique identifier for the team."),
	}, []string{"teamId"})

	return &GetTeamDetailsTool{
		ToolBase: NewToolBase("get_team_details", "Get team details", schema),
	}
}

func (t *GetTeamDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		TeamId string `json:"teamId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.TeamId == "" {
		return nil, fmt.Errorf("teamId is required")
	}

	return t.client.Get(fmt.Sprintf("/teams/%s", params.TeamId), nil)
}

func (t *GetTeamDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// UpdateTeamTool updates a team
type UpdateTeamTool struct {
	ToolBase
}

func NewUpdateTeamTool() *UpdateTeamTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"teamId": StringProperty("The unique identifier for the team."),
		"name":   StringProperty("The name of the team."),
	}, []string{"teamId"})

	return &UpdateTeamTool{
		ToolBase: NewToolBase("update_a_team", "Update a team", schema),
	}
}

func (t *UpdateTeamTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	teamId, ok := params["teamId"].(string)
	if !ok || teamId == "" {
		return nil, fmt.Errorf("teamId is required")
	}

	// Remove teamId from params as it's in the URL
	delete(params, "teamId")

	// At least one field to update must be specified
	if len(params) == 0 {
		return nil, fmt.Errorf("at least one field to update must be specified")
	}

	return t.client.Put(fmt.Sprintf("/teams/%s", teamId), params)
}

func (t *UpdateTeamTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// DeleteTeamTool deletes a team
type DeleteTeamTool struct {
	ToolBase
}

func NewDeleteTeamTool() *DeleteTeamTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"teamId": StringProperty("The unique identifier for the team."),
	}, []string{"teamId"})

	return &DeleteTeamTool{
		ToolBase: NewToolBase("delete_a_team", "Delete a team", schema),
	}
}

func (t *DeleteTeamTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		TeamId string `json:"teamId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.TeamId == "" {
		return nil, fmt.Errorf("teamId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/teams/%s", params.TeamId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true, "message": "Team deleted"}, nil
}

func (t *DeleteTeamTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}