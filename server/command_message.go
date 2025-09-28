package main

type CommandEventType string

const (
	COMMAND_CREATE CommandEventType = "create"
	COMMAND_UPDATE CommandEventType = "update"
	COMMAND_DELETE CommandEventType = "delete"
)

type CommandEvent struct {
	*Command
	Type CommandEventType `json:"type"`
}

func CommandMessage(command *Command, ty CommandEventType) *CommandEvent {
	return &CommandEvent{
		Command: command,
		Type: ty,
	}
}
