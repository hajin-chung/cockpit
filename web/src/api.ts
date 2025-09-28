import type { Command, Log } from "./schema";

type NewCommand = {
	command: string;
};

type DeleteCommand = {
	command: string;
};

export const API_ENDPOINT = import.meta.env.DEV ? "https://dev1.deps.me" : "";

export async function createCommand(command: string): Promise<Command> {
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

export async function getCommand(id: string): Promise<Command> {
	const res = await fetch(`${API_ENDPOINT}/api/v1/command/${id}`);
	if (!res.ok) {
		throw new Error("Failed to get command");
	}
	return res.json();
}

export async function getCommandList(
	before: string,
	limit: number,
): Promise<Command[]> {
	const res = await fetch(
		`${API_ENDPOINT}/api/v1/command/list?before=${before}&limit=${limit}`,
	);
	if (!res.ok) {
		throw new Error("Failed to get command list");
	}
	return res.json();
}

export async function deleteCommand(id: string) {
	const payload: DeleteCommand = { command: id };
	const res = await fetch(`${API_ENDPOINT}/api/v1/command/${id}`, {
		method: "DELETE",
		body: JSON.stringify(payload),
		headers: { "Content-Type": "application/json" },
	});
	if (!res.ok) {
		throw new Error("failed to delete command");
	}
	return;
}

export async function getLog(
	id: string,
	before: string,
	limit: number,
): Promise<Log[]> {
	const res = await fetch(
		`${API_ENDPOINT}/api/v1/command/${id}/log?before=${before}&limit=${limit}`,
	);
	if (!res.ok) {
		throw new Error("failed to get logs");
	}
	return res.json();
}

export type Stream<T> = {
	iterator: AsyncIterable<T>;
	source: EventSource;
};

export function createStream<T>(sourceUrl: string): Stream<T> {
	const eventSource = new EventSource(sourceUrl);
	const queue: T[] = [];
	let resolve: ((v: IteratorResult<T>) => void) | null = null;

	eventSource.onmessage = (evt) => {
		try {
			const value = JSON.parse(evt.data) as T;
			if (resolve) {
				resolve({ value, done: false });
				resolve = null;
			} else {
				queue.push(value);
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

	const iterator: AsyncIterable<T> = {
		[Symbol.asyncIterator]() {
			return {
				next(): Promise<IteratorResult<T>> {
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
