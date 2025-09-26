import {
	onMount,
	createSignal,
	createEffect,
	createMemo,
	type Component,
    onCleanup,
} from "solid-js";
import { A } from "@solidjs/router";
import * as api from "./api";
import type { CommandInfo, CommandLog } from "./schema";

function CommandList() {
	const [commandList, setCommandList] = createSignal<CommandInfo[]>([]);

	onMount(() => {
		api
			.getCommandList("", 50)
			.then((cl) => setCommandList(cl))
			.catch((e) => console.error(e));
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
					{commandList().map((c) => (
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
	commandId: () => string;
};

const LogList: Component<LogListProps> = ({ commandId }) => {
	const [logs, setLogs] = createSignal<CommandLog[]>([]);

	createEffect(() => {
		const id = commandId();
		const stream = createMemo(() => api.createLogStream(id));

		(async () => {
			const prevLogs = await api.getLog(commandId(), "", 50);
			setLogs(prevLogs.reverse());
			for await (const log of stream().iterator) {
				setLogs((prevLogs) => [...prevLogs, log]);
			}
		})();

		onCleanup(() => {
			stream().source.close();
		});
	});

	return (
		<div class="w-full h-full overflow-y-auto">
			<div class="w-full flex flex-col gap-2">
				{logs().map((log) => (
					<div>{log.content}</div>
				))}
			</div>
		</div>
	);
};

export { CommandList, LogList };
