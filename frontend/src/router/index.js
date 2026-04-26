import {createRouter, createWebHistory} from 'vue-router'
import Home from '@/pages/Home.vue'
import Contact from '@/pages/Contact.vue'
import ConverterView from '@/pages/ConverterView.vue'
import DownloadView from '@/pages/DownloadView.vue'

const routes = [
    {path: '/', name: 'Home', component: Home},
    {path: '/contact', name: 'Contact', component: Contact},
    {
        path: '/converter-images',
        name: 'Converter images',
        component: ConverterView
    },
    {
        path: '/download',
        name: 'Download',
        component: DownloadView
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router