import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
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
