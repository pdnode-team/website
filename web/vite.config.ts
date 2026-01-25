import { paraglideVitePlugin } from '@inlang/paraglide-js';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		paraglideVitePlugin({
			project: './project.inlang',
			outdir: './src/lib/paraglide',
			strategy: ["localStorage", "preferredLanguage", 'cookie', 'url', 'baseLocale'],
		})
	],
	server: {
		proxy: {
			// 凡是访问 /api 的请求，都转交给 PocketBase
			'/api': 'http://127.0.0.1:8090',
			// 凡是访问 PocketBase 管理面板的请求，也转发
			'/_': 'http://127.0.0.1:8090'
		}
	}
});
