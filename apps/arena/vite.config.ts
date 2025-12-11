import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import basicSsl from '@vitejs/plugin-basic-ssl';

export default defineConfig({
	plugins: [sveltekit(), basicSsl()],
	server: {
		watch: {
			paths: ['../../shared/src/lib']
		},
	},
	optimizeDeps: {
		exclude: ['shared']
	},
	ssr: {
		noExternal: ['shared']
	}
});
