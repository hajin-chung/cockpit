import { createSignal } from "solid-js";
import { Layout } from "./Layout";
import * as api from "./api";
import { useParams } from "@solidjs/router";
import { LogList } from "./Components";

const NewCommandPage = () => {
	const [command, setCommand] = createSignal("");

	const handleInput = (evt: Event) => {
		const target = evt.target as HTMLTextAreaElement;
		setCommand(target.value);
	};

	const handleRun = () => {
		api
			.createCommand(command())
			.then((c) => console.log(c))
			.catch((e) => console.error(e));
	};

	return (
		<Layout>
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
		</Layout>
	);
};

const CommandPage = () => {
	const params = useParams();

	return (
		<Layout>
			<LogList commandId={() => params.id} />
		</Layout>
	);
};


export { NewCommandPage, CommandPage };
