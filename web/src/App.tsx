import { Router, Route } from "@solidjs/router";
import { CommandPage, NewCommandPage } from "./Pages";

function App() {
	return (
		<Router>
			<Route path="/" component={NewCommandPage} />
			<Route path="/new" component={NewCommandPage} />
			<Route path="/:id" component={CommandPage} />
		</Router>
	);
}

export default App;
