import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true },
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/test',
      name: 'websocket-test',
      component: () => import('../views/WebSocketTest.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // Wait for auth to initialize
  if (authStore.loading) {
    // Watch for loading to complete
    const unwatch = router.app?.$watch(
      () => authStore.loading,
      (loading) => {
        if (!loading) {
          unwatch?.()
          checkAuth()
        }
      }
    )
    return
  }

  checkAuth()

  function checkAuth() {
    const requiresAuth = to.meta.requiresAuth
    const isPublic = to.meta.public

    if (requiresAuth && !authStore.isAuthenticated) {
      // Redirect to login if not authenticated
      next({ name: 'login' })
    } else if (isPublic && authStore.isAuthenticated) {
      // Redirect to home if already authenticated
      next({ name: 'home' })
    } else {
      next()
    }
  }
})

export default router
