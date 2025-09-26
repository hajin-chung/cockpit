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

type CommandInfo = {
	id: string;
	createdAt: string;
	command: string;
	status: CommandStatus;
};

type CommandLog = {
	id: string;
	commandId: string;
	createdAt: string;
	content: string;
	fd: LogFD;
};

export type { CommandLog, CommandInfo, CommandStatus, LogFD };
