import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import api from './services/api';
import { getCSRFService } from './services/csrf';

import './assets/styles/main.scss'

const app = createApp(App);

async function bootstrap() {
    try {
        const csrfService = getCSRFService(api);
        await csrfService.init();

        window.addEventListener('csrf:error', () => {
            console.error('CSRF token validation failed');
        });

        app.use(router);
        app.mount('#app');
    } catch (error) {
        console.error('Failed to initialize application:', error);
    }
}

bootstrap();