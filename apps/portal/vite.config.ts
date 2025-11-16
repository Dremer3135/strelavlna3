import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import monacoEditorPlugin from 'vite-plugin-monaco-editor';

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