enum CommandStatus {
	IDLE = "IDLE",
	RUNNING = "RUNNING",
	EXITED = "EXITED",
	ERROR = "ERROR",
}

enum CommandEventType {
	CREATE = "create",
	UPDATE = "update",
	DELETE = "delete",
}

enum LogFD {
	STDOUT = 1,
	STDERR = 2,
	ERROR = -1,
}

type Command = {
	id: string;
	createdAt: string;
	command: string;
	status: CommandStatus;
};

type CommandEvent = Command & {
	type: CommandEventType;
};

type Log = {
	id: string;
	commandId: string;
	createdAt: string;
	content: string;
	fd: LogFD;
};

export type { Log, Command, LogFD, CommandEvent };
export { CommandEventType, CommandStatus };
