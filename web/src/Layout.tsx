import type { Component, JSX } from "solid-js";
import { CommandList } from "./Components";

type LayoutProps = {
	children?: JSX.Element;
};

const Layout: Component<LayoutProps> = ({ children }) => {
	return (
		<div class="w-screen h-screen flex flex-col bg-neutral-900 text-neutral-50 p-8 gap-4">
			<h1 class="text-xl font-bold">Cockpit</h1>
			<div class="flex gap-4 w-full h-full flex-1 relative justify-start align-baseline overflow-hidden">
				<CommandList />
				{children}
			</div>
		</div>
	);
};

export { Layout };
