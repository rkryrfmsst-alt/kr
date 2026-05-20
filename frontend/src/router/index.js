import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthorsView from '../views/AuthorsView.vue'
import AboutView from '../views/AboutView.vue'
import ContactsView from '../views/ContactsView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/',         name: 'home',     component: HomeView },
    { path: '/artists',  name: 'artists',  component: AuthorsView },
    { path: '/about',    name: 'about',    component: AboutView },
    { path: '/contacts', name: 'contacts', component: ContactsView },
    { path: '/login',    name: 'login',    component: LoginView },
    { path: '/register', name: 'register', component: RegisterView },
  ],
})

export default router
