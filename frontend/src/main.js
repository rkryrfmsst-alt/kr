import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { useAuth } from './composables/useAuth.js'
import './assets/main.css'

const app = createApp(App)
app.use(router)

const { fetchMe } = useAuth()
fetchMe().then(() => app.mount('#app'))
