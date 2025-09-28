import { createStore } from "solid-js/store";
import type { Command } from "./schema";

const [commandStore, setCommandStore] = createStore<Command[]>([]);

export { commandStore, setCommandStore };
