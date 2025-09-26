import type { CommandInfo, CommandLog } from "./schema";

type NewCommand = {
	command: string;
};

const API_ENDPOINT = import.meta.env.DEV ? "https://dev1.deps.me" : "";

async function createCommand(command: string): Promise<CommandInfo> {
	const payload: NewCommand = { command };
	const res = await fetch(`${API_ENDPOINT}/api/v1/command/new`, {
		method: "POST",
		body: JSON.stringify(payload),
		headers: { "Content-Type": "application/json" },
	});
	if (!res.ok) {
		throw new Error("Failed to post new command");
	}
	return res.json();
}

async function getCommand(id: string): Promise<CommandInfo> {
	const res = await fetch(`${API_ENDPOINT}/api/v1/command/${id}`);
	if (!res.ok) {
		throw new Error("Failed to get command");
	}
	return res.json();
}

async function getCommandList(
	before: string,
	limit: number,
): Promise<CommandInfo[]> {
	const res = await fetch(
		`${API_ENDPOINT}/api/v1/command/list?before=${before}&limit=${limit}`,
	);
	if (!res.ok) {
		throw new Error("Failed to get command list");
	}
	return res.json();
}

async function getLog(
	id: string,
	before: string,
	limit: number,
): Promise<CommandLog[]> {
	const res = await fetch(
		`${API_ENDPOINT}/api/v1/command/${id}/log?before=${before}&limit=${limit}`,
	);
	if (!res.ok) {
		throw new Error("failed to get logs");
	}
	return res.json();
}

type LogStream = {
	iterator: AsyncIterable<CommandLog>,
	source: EventSource,
};

function createLogStream(id: string): LogStream {
	const eventSource = new EventSource(
		`${API_ENDPOINT}/api/v1/command/${id}/log/stream`,
	);

	const queue: CommandLog[] = [];
	let resolve: ((v: IteratorResult<CommandLog>) => void) | null = null;

	eventSource.onmessage = (evt) => {
		try {
			const log = JSON.parse(evt.data) as CommandLog;
			if (resolve) {
				resolve({ value: log, done: false });
				resolve = null;
			} else {
				queue.push(log);
			}
		} catch (e) {
			console.error("failed to parse SSE message: ", e);
		}
	};

	eventSource.onerror = (err) => {
		console.error("EventSource failed:", err);
		eventSource.close();
		if (resolve) {
			// Signal the end of the stream to the iterator
			resolve({ value: undefined, done: true });
			resolve = null;
		}
	};

	const iterator: AsyncIterable<CommandLog> = {
		[Symbol.asyncIterator]() {
			return {
				next(): Promise<IteratorResult<CommandLog>> {
					return new Promise((r) => {
						if (queue.length > 0) {
							// If there's a queued value, return it immediately.
							r({ value: queue.shift()!, done: false });
						} else {
							// Otherwise, store the resolver so onmessage can call it later.
							resolve = r;
						}
					});
				},
			};
		},
	};

	return { iterator, source: eventSource };
}

export { createCommand, getCommand, getCommandList, getLog, createLogStream };
