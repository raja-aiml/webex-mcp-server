package advanced_tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type ListPeopleParams struct {
	Email       string `json:"email,omitempty" query:"email"`
	DisplayName string `json:"displayName,omitempty" query:"displayName"`
	Id          string `json:"id,omitempty" query:"id"`
	OrgId       string `json:"orgId,omitempty" query:"orgId"`
	LocationId  string `json:"locationId,omitempty" query:"locationId"`
	Max         int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreatePersonParams defines the parameters for creating a person
type CreatePersonParams struct {
	Emails       []string                 `json:"emails" required:"true"`
	PhoneNumbers []map[string]interface{} `json:"phoneNumbers,omitempty"`
	Extension    string                   `json:"extension,omitempty"`
	LocationId   string                   `json:"locationId,omitempty"`
	DisplayName  string                   `json:"displayName,omitempty"`
	FirstName    string                   `json:"firstName,omitempty"`
	LastName     string                   `json:"lastName,omitempty"`
	Avatar       string                   `json:"avatar,omitempty"`
	OrgId        string                   `json:"orgId,omitempty"`
	Roles        []string                 `json:"roles,omitempty"`
	Licenses     []string                 `json:"licenses,omitempty"`
	Department   string                   `json:"department,omitempty"`
	Manager      string                   `json:"manager,omitempty"`
	ManagerId    string                   `json:"managerId,omitempty"`
	Title        string                   `json:"title,omitempty"`
	Addresses    []map[string]interface{} `json:"addresses,omitempty"`
}

// UpdatePersonParams defines the parameters for updating a person
type UpdatePersonParams struct {
	PersonId     string                   `json:"personId" required:"true"`
	Emails       []string                 `json:"emails,omitempty"`
	PhoneNumbers []map[string]interface{} `json:"phoneNumbers,omitempty"`
	Extension    string                   `json:"extension,omitempty"`
	LocationId   string                   `json:"locationId,omitempty"`
	DisplayName  string                   `json:"displayName,omitempty"`
	FirstName    string                   `json:"firstName,omitempty"`
	LastName     string                   `json:"lastName,omitempty"`
	Avatar       string                   `json:"avatar,omitempty"`
	OrgId        string                   `json:"orgId,omitempty"`
	Roles        []string                 `json:"roles,omitempty"`
	Licenses     []string                 `json:"licenses,omitempty"`
	Department   string                   `json:"department,omitempty"`
	Manager      string                   `json:"manager,omitempty"`
	ManagerId    string                   `json:"managerId,omitempty"`
	Title        string                   `json:"title,omitempty"`
	Addresses    []map[string]interface{} `json:"addresses,omitempty"`
	LoginEnabled bool                     `json:"loginEnabled,omitempty"`
}

// NewListPeopleTool creates a new list people tool
func NewListPeopleTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"email":       StringProperty("List people with this email address. For non-admin requests, require an exact match."),
		"displayName": StringProperty("List people with this display name. For non-admin requests, list people with names starting with this value."),
		"id":          StringProperty("List people with this ID. Accepts comma-separated values for bulk lookups."),
		"orgId":       StringProperty("List people in this organization. Only admin users can set this parameter."),
		"locationId":  StringProperty("List people present in this location."),
		"max":         IntegerProperty("Limit the maximum number of people in the response. Default is 100."),
	}

	return tools.NewListTool[ListPeopleParams](
		"list_people",
		"List people in your organization.",
		"/people",
		properties,
		[]string{}, // No required fields
	)
}

// NewCreatePersonTool creates a new person/user account
func NewCreatePersonTool() Tool {
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

	properties := map[string]*jsonschema.Schema{
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
	}

	return tools.NewCreateTool[CreatePersonParams](
		"create_a_person",
		"Create a new user account for a given organization. Only an admin can create a new user account.",
		"/people",
		properties,
		[]string{"emails"},
	)
}

// NewGetPersonDetailsTool gets details for a specific person
func NewGetPersonDetailsTool() Tool {
	return tools.NewGetTool(
		"get_person_details",
		"Shows details for a person by ID.",
		"/people",
		"personId",
		"A unique identifier for the person.",
	)
}

// NewUpdatePersonTool updates a person's details
func NewUpdatePersonTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"personId":     StringProperty("A unique identifier for the person."),
		"emails":       ArrayProperty("The email addresses of the person.", StringProperty("")),
		"phoneNumbers": ArrayProperty("Phone numbers for the person.", ObjectProperty("Phone number object", map[string]*jsonschema.Schema{})),
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
		"addresses":    ArrayProperty("A person's addresses.", ObjectProperty("Address object", map[string]*jsonschema.Schema{})),
		"loginEnabled": BooleanProperty("Whether the user is allowed to use Webex."),
	}

	return tools.NewUpdateTool[UpdatePersonParams](
		"update_a_person",
		"Update details for a person by ID.",
		"/people",
		"personId",
		properties,
		[]string{"personId"},
	)
}

// NewDeletePersonTool deletes a person
func NewDeletePersonTool() Tool {
	return tools.NewDeleteTool(
		"delete_a_person",
		"Remove a person from the system. Only an admin can remove a person.",
		"/people",
		"personId",
		"A unique identifier for the person.",
	)
}
