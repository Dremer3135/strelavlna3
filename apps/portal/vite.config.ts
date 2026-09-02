import fs from 'node:fs';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import monacoEditorPlugin from 'vite-plugin-monaco-editor';

// Polyfill deprecated fs.rmdirSync for compatibility with vite-plugin-monaco-editor on Node 22+
const origRmdirSync = fs.rmdirSync;
(fs as any).rmdirSync = function (path: any, options: any) {
	if (options && typeof options === 'object' && options.recursive) {
		return fs.rmSync(path, options);
	}
	return (origRmdirSync as any).apply(this, arguments);
};

const monaco = (monacoEditorPlugin as any).default ?? monacoEditorPlugin;

export default defineConfig({
	plugins: [
		sveltekit(),
		monaco({
			// languageWorkers: ['css', 'html', 'json', 'ts'] // Temporarily commented out for diagnosis
		})
	],
	server: {
		watch: {
			paths: ['../../shared/src/lib']
		}
	},
	optimizeDeps: {
		exclude: ['shared']
	},
	ssr: {
		noExternal: ['shared']
	}
});