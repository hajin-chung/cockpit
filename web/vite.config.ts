import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import solid from "vite-plugin-solid";
import { viteSingleFile } from "vite-plugin-singlefile";
import { visualizer } from "rollup-plugin-visualizer";

export default defineConfig({
	build: {
		outDir: "../server/build",
		emptyOutDir: true,
		cssCodeSplit: false,
		// rollupOptions: {
		// 	output: {
		// 		treeshake
		// 	}
		// }
	},
	server: {
		port: 4001,
		host: "0.0.0.0",
		allowedHosts: true,
	},
	plugins: [solid(), tailwindcss(), viteSingleFile(), visualizer()],
	clearScreen: false,
});
