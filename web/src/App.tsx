import { HashRouter, Route } from "@solidjs/router";
import { CommandPane, NewCommandPane } from "./Panes";
import { Layout } from "./Layout";

function App() {
	return (
		<HashRouter root={Layout}>
			<Route path="/new" component={NewCommandPane} />
			<Route path="/:id" component={CommandPane} />
		</HashRouter>
	);
}

export default App;
