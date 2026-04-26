import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
    plugins: [vue()],
    server: {
        host: '0.0.0.0',
        port: 5173,
        watch: {
            usePolling: true
        }
    },
    css: {
        devSourcemap: true,
        preprocessorOptions: {
            scss: {
                additionalData: `
                    @import "src/assets/styles/variables.scss";
                    @import "src/assets/styles/mixins.scss";
                `
            }
        }
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src')
        }
    },
})