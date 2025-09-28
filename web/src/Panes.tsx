import { createEffect, createMemo, createSignal } from "solid-js";
import * as api from "./api";
import { useNavigate, useParams } from "@solidjs/router";
import { commandStore } from "./store";
import { LogList } from "./Components";
import { CommandStatus } from "./schema";

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
	const navigate = useNavigate();
	const params = useParams();
	const command = createMemo(() =>
		commandStore.find((c) => c.id === params.id),
	);

	const handleDelete = () => {
		api
			.deleteCommand(params.id)
			.then(() => navigate("/"))
			.catch((e) => console.error(e));
	};

	return (
		<div class="w-full h-full">
			<div>
				<LogList id={params.id} />
				{(command()?.status === CommandStatus.EXITED ||
					command()?.status === CommandStatus.ERROR) && (
					<button
						class="p-2 bg-violet-900 hover:bg-violet-800 transition-all rounded-lg self-end"
						onClick={handleDelete}
					>
						Delete
					</button>
				)}
			</div>
		</div>
	);
};

export { NewCommandPane, CommandPane };
