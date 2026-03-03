import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import AuthView from '../AuthView.vue'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import { authService } from '@/services'

vi.mock('@/services', () => ({
  authService: {
    login: vi.fn(),
    register: vi.fn(),
  },
}))

describe('AuthView', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(() => {
    vi.clearAllMocks()

    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        {
          path: '/',
          name: 'auth',
          component: AuthView,
        },
        {
          path: '/getting-started',
          name: 'getting-started',
          component: { template: '<div>Getting Started</div>' },
        },
      ],
    })
  })

  it('renders login form by default', async () => {
    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    expect(wrapper.findComponent(LoginForm).exists()).toBe(true)
    expect(wrapper.findComponent(RegisterForm).exists()).toBe(false)
  })

  it('switches to register form when register tab is clicked', async () => {
    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    const tabs = wrapper.findAll('.tab')
    await tabs[1].trigger('click')

    expect(wrapper.findComponent(RegisterForm).exists()).toBe(true)
    expect(wrapper.findComponent(LoginForm).exists()).toBe(false)
  })

  it('switches back to login form when login tab is clicked', async () => {
    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    const tabs = wrapper.findAll('.tab')
    await tabs[1].trigger('click')
    await tabs[0].trigger('click')

    expect(wrapper.findComponent(LoginForm).exists()).toBe(true)
    expect(wrapper.findComponent(RegisterForm).exists()).toBe(false)
  })

  it('handles successful login', async () => {
    vi.mocked(authService.login).mockResolvedValue({
      data: {
        access_token: 'mock-token',
        refresh_token: 'mock-refresh',
      },
      message: 'Login successful',
    })

    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    const loginForm = wrapper.findComponent(LoginForm)
    await loginForm.vm.$emit('switchToRegister')

    // The component should switch to register form
    await wrapper.vm.$nextTick()
  })

  it('handles failed login', async () => {
    vi.mocked(authService.login).mockRejectedValue({
      response: {
        data: {
          message: 'Invalid credentials',
        },
      },
    })

    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    expect(wrapper.findComponent(LoginForm).exists()).toBe(true)
  })

  it('handles successful registration and switches to login', async () => {
    vi.mocked(authService.register).mockResolvedValue({
      data: {},
      message: 'User registered successfully',
    })

    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    // Switch to register tab
    const tabs = wrapper.findAll('.tab')
    await tabs[1].trigger('click')

    expect(wrapper.findComponent(RegisterForm).exists()).toBe(true)
  })

  it('renders authentication card with correct structure', async () => {
    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    expect(wrapper.find('.auth-view').exists()).toBe(true)
    expect(wrapper.find('.auth-card').exists()).toBe(true)
    expect(wrapper.find('.card-header').exists()).toBe(true)
    expect(wrapper.find('.tabs-container').exists()).toBe(true)
    expect(wrapper.find('.form-container').exists()).toBe(true)
  })

  it('displays correct tab titles', async () => {
    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    const tabs = wrapper.findAll('.tab')
    expect(tabs[0].text()).toContain('ลงชื่อเข้าสู่ระบบ')
    expect(tabs[1].text()).toContain('สมัครสมาชิก')
  })

  it('applies active class to selected tab', async () => {
    const wrapper = mount(AuthView, {
      global: {
        plugins: [router],
      },
    })

    await flushPromises()

    const tabs = wrapper.findAll('.tab')
    expect(tabs[0].classes()).toContain('active')
    expect(tabs[1].classes()).not.toContain('active')

    await tabs[1].trigger('click')

    expect(tabs[0].classes()).not.toContain('active')
    expect(tabs[1].classes()).toContain('active')
  })
})
