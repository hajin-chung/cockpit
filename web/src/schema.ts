enum CommandStatus {
	IDLE = "IDLE",
	RUNNING = "RUNNING",
	EXITED = "EXITED",
	ERROR = "ERROR",
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

type Log = {
	id: string;
	commandId: string;
	createdAt: string;
	content: string;
	fd: LogFD;
};

export type { Log, Command, CommandStatus, LogFD };
