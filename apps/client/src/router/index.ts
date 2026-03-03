import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '../views/AuthView.vue'
import { authService } from '@/services'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'auth',
      component: AuthView,
    },
    {
      path: '/getting-started',
      name: 'getting-started',
      component: () => import('../views/GettingStartedView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    try {
      await authService.me()
      next()
    } catch (error) {
      console.error('Authentication failed:', error)
      next('/')
    }
  } else {
    next()
  }
})

export default router
