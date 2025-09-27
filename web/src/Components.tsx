import { createSignal, createEffect, createMemo, onCleanup } from "solid-js";
import { A } from "@solidjs/router";
import * as api from "./api";
import type { Command } from "./schema";

function CommandList() {
	const [commandList, setCommandList] = createSignal<Command[]>([]);
	const stream = createMemo(() => {
		const ret = api.createStream<Command>(
			`${api.API_ENDPOINT}/api/v1/command/stream`,
		);
		onCleanup(() => ret.source.close());
		return ret;
	});
	const fetcher = async (prev: Command[]) => {
		if (prev.length === 0)
			return await api.getCommandList("", 50);
		else {
			const commands = await api.getCommandList(prev.at(-1)!.id, 50);
			return [...prev, ...commands]
		}
	}

	// const loadMore = () => fetcher(commandList()).then(setCommandList);

	createEffect(() => {
		(async () => {
			const prevCommands = await fetcher([]);
			setCommandList(prevCommands);

			for await (const command of stream().iterator) {
				setCommandList((prevCommands) => {
					const idx = prevCommands.findIndex((c) => c.id === command.id);
					if (idx < 0) {
						return [command, ...prevCommands];
					} else {
						prevCommands[idx].status = command.status;
						return [...prevCommands];
					}
				});
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

export { CommandList };
