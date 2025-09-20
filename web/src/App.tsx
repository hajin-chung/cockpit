import { createSignal, type Component } from "solid-js";

type PaneType = "Empty" | "New" | "Log";

function App() {
	const [paneType, setPaneType] = createSignal<PaneType>("Empty");

	return (
		<div class="w-screen h-screen overflow-hidden flex flex-col bg-neutral-900 text-neutral-50 p-8 gap-4">
			<h1 class="text-xl font-bold">Cockpit</h1>
			<div class="flex gap-4 w-full h-full justify-start align-baseline">
				<CommandList setPaneType={setPaneType} />
				{paneType() === "New" && <PaneNew />}
				{paneType() === "Log" && <PaneLog />}
			</div>
		</div>
	);
}

const CommandList: Component<{
	setPaneType: (paneType: PaneType) => void;
}> = ({ setPaneType }) => {
	// fetch commands
	return (
		<div class="w-1/4 bg-neutral-800 shadow-md rounded-md h-full p-2 shrink-0">
			<button
				class="p-2 hover:bg-neutral-700 transition-all w-full rounded-lg"
				onclick={() => setPaneType("New")}
			>
				New
			</button>
		</div>
	);
};

const PaneNew: Component = () => {
	const [command, setCommand] = createSignal("");

	const handleInput = (evt: Event) => {
		const target = evt.target as HTMLTextAreaElement;
		setCommand(target.value);
	};

	const handleRun = () => {
		console.log("hi");
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

const PaneLog: Component = () => {
	return <div>log!</div>;
};

export default App;
