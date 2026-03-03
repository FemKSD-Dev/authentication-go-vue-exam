import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createRouter, createMemoryHistory } from 'vue-router'
import { authService } from '@/services'

vi.mock('@/services', () => ({
  authService: {
    me: vi.fn(),
  },
}))

describe('Router', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  const createTestRouter = () => {
    return createRouter({
      history: createMemoryHistory(),
      routes: [
        {
          path: '/',
          name: 'auth',
          component: { template: '<div>Auth</div>' },
        },
        {
          path: '/getting-started',
          name: 'getting-started',
          component: { template: '<div>Getting Started</div>' },
          meta: { requiresAuth: true },
        },
        {
          path: '/public',
          name: 'public',
          component: { template: '<div>Public</div>' },
        },
      ],
    })
  }

  const setupNavigationGuard = (router: ReturnType<typeof createTestRouter>) => {
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
  }

  it('allows navigation to public routes without authentication', async () => {
    const router = createTestRouter()
    setupNavigationGuard(router)

    await router.push('/')
    expect(router.currentRoute.value.path).toBe('/')

    await router.push('/public')
    expect(router.currentRoute.value.path).toBe('/public')
  })

  it('allows navigation to protected routes when authenticated', async () => {
    vi.mocked(authService.me).mockResolvedValue({
      data: {
        user_id: '123',
        username: 'testuser',
      },
      message: 'Success',
    })

    const router = createTestRouter()
    setupNavigationGuard(router)

    await router.push('/getting-started')

    expect(authService.me).toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/getting-started')
  })

  it('redirects to home when accessing protected route without authentication', async () => {
    vi.mocked(authService.me).mockRejectedValue(new Error('Unauthorized'))

    const router = createTestRouter()
    setupNavigationGuard(router)

    await router.push('/getting-started')

    expect(authService.me).toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/')
  })

  it('has correct route configuration', () => {
    const router = createTestRouter()

    const routes = router.getRoutes()

    expect(routes).toHaveLength(3)
    expect(routes.find((r) => r.path === '/')).toBeDefined()
    expect(routes.find((r) => r.path === '/getting-started')).toBeDefined()
    expect(routes.find((r) => r.path === '/public')).toBeDefined()
  })

  it('marks getting-started route as requiring authentication', () => {
    const router = createTestRouter()

    const gettingStartedRoute = router.getRoutes().find((r) => r.path === '/getting-started')

    expect(gettingStartedRoute?.meta.requiresAuth).toBe(true)
  })

  it('does not mark auth route as requiring authentication', () => {
    const router = createTestRouter()

    const authRoute = router.getRoutes().find((r) => r.path === '/')

    expect(authRoute?.meta.requiresAuth).toBeUndefined()
  })

  it('calls authService.me only for protected routes', async () => {
    vi.mocked(authService.me).mockResolvedValue({
      data: {
        user_id: '123',
        username: 'testuser',
      },
      message: 'Success',
    })

    const router = createTestRouter()
    setupNavigationGuard(router)

    await router.push('/')
    expect(authService.me).not.toHaveBeenCalled()

    await router.push('/getting-started')
    expect(authService.me).toHaveBeenCalledTimes(1)

    await router.push('/public')
    expect(authService.me).toHaveBeenCalledTimes(1) // Still 1, not called again
  })

  it('handles navigation guard errors gracefully', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    vi.mocked(authService.me).mockRejectedValue(new Error('Network error'))

    const router = createTestRouter()
    setupNavigationGuard(router)

    await router.push('/getting-started')

    expect(consoleErrorSpy).toHaveBeenCalledWith(
      'Authentication failed:',
      expect.any(Error),
    )
    expect(router.currentRoute.value.path).toBe('/')

    consoleErrorSpy.mockRestore()
  })
})
