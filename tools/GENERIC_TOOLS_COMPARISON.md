# Generic Tools Comparison

## Overview
This document demonstrates how generic tools reduce boilerplate code while maintaining type safety.

## Before: Traditional Tool Implementation

```go
// ListMessagesTool - 71 lines of code
type ListMessagesTool struct {
    ToolBase
}

func NewListMessagesTool() *ListMessagesTool {
    schema := SimpleSchema(map[string]*jsonschema.Schema{
        "roomId":          StringProperty("The ID of the room to list messages from."),
        "parentId":        StringProperty("The ID of the parent message to filter by."),
        "mentionedPeople": StringProperty("List messages with these people mentioned, by ID."),
        "before":          StringProperty("List messages sent before a specific date and time."),
        "beforeMessage":   StringProperty("List messages sent before a specific message, by ID."),
        "max":             IntegerProperty("Limit the maximum number of messages in the response."),
    }, []string{"roomId"})

    return &ListMessagesTool{
        ToolBase: NewToolBase("list_messages", "List messages in a Webex room.", schema),
    }
}

func (t *ListMessagesTool) Execute(args json.RawMessage) (interface{}, error) {
    var params struct {
        RoomId          string `json:"roomId"`
        ParentId        string `json:"parentId,omitempty"`
        MentionedPeople string `json:"mentionedPeople,omitempty"`
        Before          string `json:"before,omitempty"`
        BeforeMessage   string `json:"beforeMessage,omitempty"`
        Max             int    `json:"max,omitempty"`
    }

    if err := json.Unmarshal(args, &params); err != nil {
        return nil, fmt.Errorf("failed to parse arguments: %w", err)
    }

    if params.RoomId == "" {
        return nil, fmt.Errorf("roomId is required")
    }

    queryParams := map[string]string{"roomId": params.RoomId}
    if params.ParentId != "" {
        queryParams["parentId"] = params.ParentId
    }
    if params.MentionedPeople != "" {
        queryParams["mentionedPeople"] = params.MentionedPeople
    }
    if params.Before != "" {
        queryParams["before"] = params.Before
    }
    if params.BeforeMessage != "" {
        queryParams["beforeMessage"] = params.BeforeMessage
    }
    if params.Max > 0 {
        queryParams["max"] = strconv.Itoa(params.Max)
    }

    return t.client.Get("/messages", queryParams)
}

func (t *ListMessagesTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
    return ExecuteWithMapBase(t, args)
}
```

## After: Generic Tool Implementation

```go
// ListMessagesParams - Define the parameters
type ListMessagesParams struct {
    RoomId          string `json:"roomId" required:"true"`
    ParentId        string `json:"parentId,omitempty"`
    MentionedPeople string `json:"mentionedPeople,omitempty"`
    Before          string `json:"before,omitempty"`
    BeforeMessage   string `json:"beforeMessage,omitempty"`
    Max             int    `json:"max,omitempty"`
}

// NewListMessagesToolGeneric - 15 lines instead of 71
func NewListMessagesToolGeneric() Tool {
    properties := map[string]*jsonschema.Schema{
        "roomId":          StringProperty("The ID of the room to list messages from."),
        "parentId":        StringProperty("The ID of the parent message to filter by."),
        "mentionedPeople": StringProperty("List messages with these people mentioned, by ID."),
        "before":          StringProperty("List messages sent before a specific date and time."),
        "beforeMessage":   StringProperty("List messages sent before a specific message, by ID."),
        "max":             IntegerProperty("Limit the maximum number of messages in the response."),
    }

    return NewListTool[ListMessagesParams](
        "list_messages_generic",
        "List messages in a Webex room (generic implementation).",
        "/messages",
        properties,
    )
}
```

## Benefits

### 1. Code Reduction
- **Traditional**: ~71 lines per tool
- **Generic**: ~15-20 lines per tool
- **Savings**: ~75% reduction in boilerplate

### 2. Automatic Features
- ✅ Parameter validation (required fields)
- ✅ Query parameter building
- ✅ JSON marshaling/unmarshaling
- ✅ Error handling
- ✅ ExecuteWithMap implementation

### 3. Type Safety
- Compile-time type checking with generics
- Struct tags for validation and serialization
- No runtime type assertions needed

### 4. Consistency
- All tools follow the same pattern
- Reduces bugs from copy-paste errors
- Easier to maintain and update

### 5. Factory Functions for Common Patterns

```go
// Simple CRUD operations become one-liners:

// List operation
NewListTool[RoomListParams]("list_rooms", "List rooms", "/rooms", properties)

// Get by ID
NewGetTool("get_room", "Get room details", "/rooms", "roomId")

// Create
NewCreateTool[CreateRoomParams]("create_room", "Create room", "/rooms", properties, required)

// Update
NewUpdateTool[UpdateRoomParams]("update_room", "Update room", "/rooms", "roomId", properties, required)

// Delete
NewDeleteTool("delete_room", "Delete room", "/rooms", "roomId")
```

## Migration Strategy

1. Keep existing tools working
2. Implement new tools using generics
3. Gradually migrate existing tools when updating them
4. Both implementations can coexist

## Example Usage

```go
// In plugin_loader.go
func (p *messagingPlugin) Register(registry *Registry) error {
    tools := []Tool{
        // Mix of traditional and generic tools
        NewListMessagesTool(),           // Traditional
        NewListMessagesToolGeneric(),    // Generic version
        
        // New tools can use generics from the start
        NewListTool[CustomParams]("custom_list", "List custom items", "/custom", props),
    }
    
    for _, tool := range tools {
        if err := registry.Register(tool); err != nil {
            return err
        }
    }
    return nil
}
```

## Conclusion

Generic tools provide:
- **75% less code** to write and maintain
- **Type safety** without runtime assertions
- **Consistency** across all tools
- **Flexibility** for custom validation logic
- **Backward compatibility** with existing tools

The generic approach follows Go best practices and makes the codebase more maintainable while reducing the chance of bugs.