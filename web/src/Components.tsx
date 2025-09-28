import {
	createEffect,
	createMemo,
	createSignal,
	onCleanup,
	type Component,
} from "solid-js";
import { A } from "@solidjs/router";
import * as api from "./api";
import {
	CommandEventType,
	type Command,
	type CommandEvent,
	type Log,
} from "./schema";
import { commandStore, setCommandStore } from "./store";

function CommandList() {
	const stream = createMemo(() => {
		const ret = api.createStream<CommandEvent>(
			`${api.API_ENDPOINT}/api/v1/command/stream`,
		);
		onCleanup(() => ret.source.close());
		return ret;
	});
	const fetcher = async (prev: Command[]) => {
		if (prev.length === 0) return await api.getCommandList("", 50);
		else {
			const commands = await api.getCommandList(prev.at(-1)!.id, 50);
			return [...prev, ...commands];
		}
	};

	// const loadMore = () => fetcher(commandList()).then(setCommandList);

	createEffect(() => {
		(async () => {
			const prevCommands = await fetcher([]);
			setCommandStore(prevCommands);

			for await (const command of stream().iterator) {
				switch (command.type) {
					case CommandEventType.CREATE:
						setCommandStore((prev) => [command, ...prev]);
						break
					case CommandEventType.UPDATE:
						setCommandStore(
							(c) => c.id === command.id,
							"status",
							command.status,
						);
						break
					case CommandEventType.DELETE:
						setCommandStore((prev) => prev.filter((c) => c.id !== command.id));
						break
				}
			}
		})();
	});

	return (
		<div class="w-1/4 bg-neutral-800 shadow-md rounded-md h-full p-2 shrink-0 flex flex-col gap-2">
			<A
				class="p-2 hover:bg-neutral-700 transition-all w-full rounded-lg"
				href="/new"
			>
				New
			</A>
			<div class="w-full h-full overflow-y-auto">
				<div class="w-full flex flex-col gap-2">
					{commandStore.map((c) => (
						<A
							class="p-2 hover:bg-neutral-700 transition-all w-full rounded-lg"
							href={`/${c.id}`}
						>
							{c.command}
						</A>
					))}
				</div>
			</div>
		</div>
	);
}

type LogListProps = {
	id: string;
};

const LogList: Component<LogListProps> = ({ id }) => {
	const [logs, setLogs] = createSignal<Log[]>([]);
	const [showLoadMore, setShowLoadMore] = createSignal(true);

	const stream = createMemo(() => {
		const ret = api.createStream<Log>(
			`${api.API_ENDPOINT}/api/v1/command/${id}/log/stream`,
		);
		onCleanup(() => ret.source.close());
		return ret;
	});
	const fetcher = async (prevLogs: Log[]) => {
		const before = prevLogs.length != 0 ? prevLogs.at(-1)!.id : "";
		const logs = await api.getLog(id, before, 50);
		setShowLoadMore(logs.length === 50);
		return [...prevLogs, ...logs];
	};

	const loadMore = () => fetcher(logs()).then(setLogs);

	createEffect(() => {
		(async () => {
			const prevLogs = await fetcher([]);
			setLogs(prevLogs);
			for await (const log of stream().iterator) {
				setLogs((prevLogs) => [log, ...prevLogs]);
			}
		})();
	});

	return (
		<div class="w-full h-full flex overflow-y-auto gap-2 flex-col-reverse">
			{logs().map((log) => (
				<div>{log.content}</div>
			))}
			{showLoadMore() && (
				<button
					onclick={loadMore}
					class="p-2 hover:bg-neutral-700 transition-all rounded-lg self-center hover:cursor-pointer"
				>
					Load More
				</button>
			)}
		</div>
	);
};

export { CommandList, LogList };
