package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// ListPeopleTool lists people in the organization
type ListPeopleTool struct {
	ToolBase
}

func NewListPeopleTool() *ListPeopleTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"email":       StringProperty("List people with this email address. For non-admin requests, require an exact match."),
		"displayName": StringProperty("List people with this display name. For non-admin requests, list people with names starting with this value."),
		"id":          StringProperty("List people with this ID. Accepts comma-separated values for bulk lookups."),
		"orgId":       StringProperty("List people in this organization. Only admin users can set this parameter."),
		"locationId":  StringProperty("List people present in this location."),
		"max":         IntegerProperty("Limit the maximum number of people in the response. Default is 100."),
	}, []string{})

	return &ListPeopleTool{
		ToolBase: NewToolBase("list_people", "List people in your organization.", schema),
	}
}

func (t *ListPeopleTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		Email       string `json:"email,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Id          string `json:"id,omitempty"`
		OrgId       string `json:"orgId,omitempty"`
		LocationId  string `json:"locationId,omitempty"`
		Max         int    `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	queryParams := make(map[string]string)
	if params.Email != "" {
		queryParams["email"] = params.Email
	}
	if params.DisplayName != "" {
		queryParams["displayName"] = params.DisplayName
	}
	if params.Id != "" {
		queryParams["id"] = params.Id
	}
	if params.OrgId != "" {
		queryParams["orgId"] = params.OrgId
	}
	if params.LocationId != "" {
		queryParams["locationId"] = params.LocationId
	}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	} else {
		queryParams["max"] = "100"
	}

	return t.client.Get("/people", queryParams)
}

func (t *ListPeopleTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// CreatePersonTool creates a new person/user account
type CreatePersonTool struct {
	ToolBase
}

func NewCreatePersonTool() *CreatePersonTool {
	phoneNumberSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"type":  StringProperty("Phone number type"),
			"value": StringProperty("Phone number value"),
		},
	}

	addressSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"type":          StringProperty("Address type"),
			"country":       StringProperty("Country"),
			"locality":      StringProperty("Locality"),
			"postalCode":    StringProperty("Postal code"),
			"region":        StringProperty("Region"),
			"streetAddress": StringProperty("Street address"),
		},
	}

	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"emails":       ArrayProperty("The email addresses of the person.", StringProperty("")),
		"phoneNumbers": ArrayProperty("Phone numbers for the person.", phoneNumberSchema),
		"extension":    StringProperty("The Webex Calling extension of the person."),
		"locationId":   StringProperty("The ID of the location for this person."),
		"displayName":  StringProperty("The full name of the person."),
		"firstName":    StringProperty("The first name of the person."),
		"lastName":     StringProperty("The last name of the person."),
		"avatar":       StringProperty("The URL to the person's avatar in PNG format."),
		"orgId":        StringProperty("The ID of the organization to which this person belongs."),
		"roles":        ArrayProperty("An array of role strings representing the roles to which this person belongs.", StringProperty("")),
		"licenses":     ArrayProperty("An array of license strings allocated to this person.", StringProperty("")),
		"department":   StringProperty("The business department the user belongs to."),
		"manager":      StringProperty("A manager identifier."),
		"managerId":    StringProperty("The person ID of the manager."),
		"title":        StringProperty("The person's title."),
		"addresses":    ArrayProperty("A person's addresses.", addressSchema),
	}, []string{"emails"})

	return &CreatePersonTool{
		ToolBase: NewToolBase("create_a_person", "Create a new user account for a given organization. Only an admin can create a new user account.", schema),
	}
}

func (t *CreatePersonTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params["emails"] == nil {
		return nil, fmt.Errorf("emails is required")
	}

	return t.client.Post("/people", params)
}

func (t *CreatePersonTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// GetPersonDetailsTool gets details for a specific person
type GetPersonDetailsTool struct {
	ToolBase
}

func NewGetPersonDetailsTool() *GetPersonDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"personId": StringProperty("A unique identifier for the person."),
	}, []string{"personId"})

	return &GetPersonDetailsTool{
		ToolBase: NewToolBase("get_person_details", "Shows details for a person by ID.", schema),
	}
}

func (t *GetPersonDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		PersonId string `json:"personId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.PersonId == "" {
		return nil, fmt.Errorf("personId is required")
	}

	return t.client.Get(fmt.Sprintf("/people/%s", params.PersonId), nil)
}

func (t *GetPersonDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// UpdatePersonTool updates a person's details
type UpdatePersonTool struct {
	ToolBase
}

func NewUpdatePersonTool() *UpdatePersonTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"personId":     StringProperty("A unique identifier for the person."),
		"emails":       ArrayProperty("The email addresses of the person.", StringProperty("")),
		"phoneNumbers": ArrayProperty("Phone numbers for the person.", ObjectProperty("")),
		"extension":    StringProperty("The Webex Calling extension of the person."),
		"locationId":   StringProperty("The ID of the location for this person."),
		"displayName":  StringProperty("The full name of the person."),
		"firstName":    StringProperty("The first name of the person."),
		"lastName":     StringProperty("The last name of the person."),
		"avatar":       StringProperty("The URL to the person's avatar in PNG format."),
		"orgId":        StringProperty("The ID of the organization to which this person belongs."),
		"roles":        ArrayProperty("An array of role strings representing the roles to which this person belongs.", StringProperty("")),
		"licenses":     ArrayProperty("An array of license strings allocated to this person.", StringProperty("")),
		"department":   StringProperty("The business department the user belongs to."),
		"manager":      StringProperty("A manager identifier."),
		"managerId":    StringProperty("The person ID of the manager."),
		"title":        StringProperty("The person's title."),
		"addresses":    ArrayProperty("A person's addresses.", ObjectProperty("")),
		"loginEnabled": BooleanProperty("Whether the user is allowed to use Webex."),
	}, []string{"personId"})

	return &UpdatePersonTool{
		ToolBase: NewToolBase("update_a_person", "Update details for a person by ID.", schema),
	}
}

func (t *UpdatePersonTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	personId, ok := params["personId"].(string)
	if !ok || personId == "" {
		return nil, fmt.Errorf("personId is required")
	}

	delete(params, "personId")
	return t.client.Put(fmt.Sprintf("/people/%s", personId), params)
}

func (t *UpdatePersonTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// DeletePersonTool deletes a person
type DeletePersonTool struct {
	ToolBase
}

func NewDeletePersonTool() *DeletePersonTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"personId": StringProperty("A unique identifier for the person."),
	}, []string{"personId"})

	return &DeletePersonTool{
		ToolBase: NewToolBase("delete_a_person", "Remove a person from the system. Only an admin can remove a person.", schema),
	}
}

func (t *DeletePersonTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		PersonId string `json:"personId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.PersonId == "" {
		return nil, fmt.Errorf("personId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/people/%s", params.PersonId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Person deleted successfully",
	}, nil
}

func (t *DeletePersonTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// GetMyOwnDetailsTool gets the current user's details
type GetMyOwnDetailsTool struct {
	ToolBase
}

func NewGetMyOwnDetailsTool() *GetMyOwnDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{}, []string{})

	return &GetMyOwnDetailsTool{
		ToolBase: NewToolBase("get_my_own_details", "Get the details of the authenticated user.", schema),
	}
}

func (t *GetMyOwnDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	return t.client.Get("/people/me", nil)
}

func (t *GetMyOwnDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}
