import { createEffect, createMemo, createSignal, onCleanup } from "solid-js";
import * as api from "./api";
import { useNavigate, useParams } from "@solidjs/router";
import type { Log } from "./schema";

const NewCommandPane = () => {
	const navigate = useNavigate();
	const [command, setCommand] = createSignal("");

	const handleInput = (evt: Event) => {
		const target = evt.target as HTMLTextAreaElement;
		setCommand(target.value);
	};

	const handleRun = () => {
		api
			.createCommand(command())
			.then((c) => navigate(`/${c.id}`))
			.catch((e) => console.error(e));
	};

	return (
		<div class="flex flex-col gap-2 w-full items-start">
			<p class="text-lg font-bold">New Command</p>
			<textarea
				class="w-full h-full bg-neutral-800 rounded-md shadow-md p-2 resize-none"
				textContent={command()}
				onInput={handleInput}
			/>
			<button
				class="p-2 bg-violet-900 hover:bg-violet-800 transition-all rounded-lg self-end"
				onClick={handleRun}
			>
				Run
			</button>
		</div>
	);
};

const CommandPane = () => {
	const params = useParams();
	const [logs, setLogs] = createSignal<Log[]>([]);
	const [showLoadMore, setShowLoadMore] = createSignal(true);
	const stream = createMemo(() => {
		const ret = api.createStream<Log>(
			`${api.API_ENDPOINT}/api/v1/command/${params.id}/log/stream`,
		);
		onCleanup(() => ret.source.close());
		return ret;
	});
	const fetcher = async (prevLogs: Log[]) => {
		const before = prevLogs.length != 0 ? prevLogs.at(-1)!.id : "";
		const logs = await api.getLog(params.id, before, 50);
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
		// <div class="w-full h-full overflow-y-auto">
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
		// </div>
	);
};

export { NewCommandPane, CommandPane };
