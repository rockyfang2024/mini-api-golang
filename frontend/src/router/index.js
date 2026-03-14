import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import MyPostsView from '../views/MyPostsView.vue'
import NotificationsView from '../views/NotificationsView.vue'
import UserProfileView from '../views/UserProfileView.vue'
import SettingsView from '../views/SettingsView.vue'

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/login', name: 'login', component: LoginView },
  { path: '/register', name: 'register', component: RegisterView },
  { path: '/my-posts', name: 'my-posts', component: MyPostsView, meta: { requiresAuth: true } },
  { path: '/notifications', name: 'notifications', component: NotificationsView, meta: { requiresAuth: true } },
  { path: '/users/:id', name: 'user-profile', component: UserProfileView },
  { path: '/settings', name: 'settings', component: SettingsView, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !localStorage.getItem('token')) {
    return { name: 'login' }
  }
})

export default router
