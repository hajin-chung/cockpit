import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import solid from "vite-plugin-solid";

export default defineConfig({
	build: {
		outDir: "../build/static",
		emptyOutDir: true,
	},
	server: {
		port: 4000,
		host: "0.0.0.0",
		allowedHosts: true,
	},
	plugins: [solid(), tailwindcss()],
});
