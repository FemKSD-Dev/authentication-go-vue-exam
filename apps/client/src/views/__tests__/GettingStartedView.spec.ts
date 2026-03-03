import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import GettingStartedView from '../GettingStartedView.vue'
import UserDetailCard from '@/components/UserDetailCard.vue'
import { authService } from '@/services'

vi.mock('@/services', () => ({
  authService: {
    me: vi.fn(),
  },
}))

describe('GettingStartedView', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(() => {
    vi.clearAllMocks()

    router = createRouter({
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
          component: GettingStartedView,
        },
      ],
    })
  })

  it('shows loading state initially', async () => {
    vi.mocked(authService.me).mockImplementation(
      () => new Promise(() => {}), // Never resolves
    )

    const wrapper = mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.find('.loading-container').exists()).toBe(true)
    expect(wrapper.text()).toContain('กำลังโหลดข้อมูล')
    expect(wrapper.findComponent(UserDetailCard).exists()).toBe(false)
  })

  it('displays user card after successful data fetch', async () => {
    vi.mocked(authService.me).mockResolvedValue({
      data: {
        user_id: '123',
        username: 'testuser',
      },
      message: 'Success',
    })

    const wrapper = mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    // Wait for the timeout in the component
    await new Promise((resolve) => setTimeout(resolve, 150))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.loading-container').exists()).toBe(false)
    expect(wrapper.findComponent(UserDetailCard).exists()).toBe(true)
  })

  it('passes username to UserDetailCard', async () => {
    vi.mocked(authService.me).mockResolvedValue({
      data: {
        user_id: '123',
        username: 'testuser',
      },
      message: 'Success',
    })

    const wrapper = mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()
    await new Promise((resolve) => setTimeout(resolve, 150))
    await wrapper.vm.$nextTick()

    const userCard = wrapper.findComponent(UserDetailCard)
    expect(userCard.props('username')).toBe('testuser')
  })

  it('redirects to home on failed data fetch', async () => {
    const pushSpy = vi.spyOn(router, 'push')

    vi.mocked(authService.me).mockRejectedValue(new Error('Unauthorized'))

    mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    expect(pushSpy).toHaveBeenCalledWith('/')
  })

  it('renders with correct structure', async () => {
    vi.mocked(authService.me).mockResolvedValue({
      data: {
        user_id: '123',
        username: 'testuser',
      },
      message: 'Success',
    })

    const wrapper = mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    expect(wrapper.find('.getting-started-view').exists()).toBe(true)
  })

  it('calls authService.me on mount', async () => {
    vi.mocked(authService.me).mockResolvedValue({
      data: {
        user_id: '123',
        username: 'testuser',
      },
      message: 'Success',
    })

    mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    expect(authService.me).toHaveBeenCalledTimes(1)
  })

  it('handles loading spinner animation', async () => {
    vi.mocked(authService.me).mockImplementation(
      () => new Promise(() => {}), // Never resolves
    )

    const wrapper = mount(GettingStartedView, {
      global: {
        plugins: [router],
      },
    })

    const spinner = wrapper.find('.loading-spinner')
    expect(spinner.exists()).toBe(true)
    expect(spinner.classes()).toContain('loading-spinner')
  })
})
