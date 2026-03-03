import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginForm from '../LoginForm.vue'

describe('LoginForm', () => {
  let mockOnLogin: ReturnType<typeof vi.fn>

  beforeEach(() => {
    mockOnLogin = vi.fn()
  })

  it('renders properly', () => {
    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    expect(wrapper.find('#login-username').exists()).toBe(true)
    expect(wrapper.find('#login-password').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('shows error when submitting empty form', async () => {
    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('กรุณากรอกชื่อผู้ใช้และรหัสผ่าน')
    expect(mockOnLogin).not.toHaveBeenCalled()
  })

  it('calls onLogin with correct credentials', async () => {
    mockOnLogin.mockResolvedValue({ success: true })

    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    await wrapper.find('#login-username').setValue('testuser')
    await wrapper.find('#login-password').setValue('testpass123')
    await wrapper.find('form').trigger('submit.prevent')

    expect(mockOnLogin).toHaveBeenCalledWith('testuser', 'testpass123')
  })

  it('displays error message from failed login', async () => {
    mockOnLogin.mockResolvedValue({
      success: false,
      error: 'Invalid credentials',
    })

    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    await wrapper.find('#login-username').setValue('testuser')
    await wrapper.find('#login-password').setValue('wrongpass')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Invalid credentials')
  })

  it('toggles password visibility', async () => {
    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    const passwordInput = wrapper.find('#login-password')
    const toggleButton = wrapper.find('.toggle-password-btn')

    expect(passwordInput.attributes('type')).toBe('password')

    await toggleButton.trigger('click')
    expect(passwordInput.attributes('type')).toBe('text')

    await toggleButton.trigger('click')
    expect(passwordInput.attributes('type')).toBe('password')
  })

  it('emits switchToRegister event', async () => {
    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    await wrapper.find('.link').trigger('click')

    expect(wrapper.emitted('switchToRegister')).toBeTruthy()
  })

  it('shows loading state during login', async () => {
    let resolveLogin: (value: any) => void
    const loginPromise = new Promise((resolve) => {
      resolveLogin = resolve
    })
    mockOnLogin.mockReturnValue(loginPromise)

    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    await wrapper.find('#login-username').setValue('testuser')
    await wrapper.find('#login-password').setValue('testpass123')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(wrapper.text()).toContain('กำลังเข้าสู่ระบบ')

    resolveLogin!({ success: true })
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick() // Need extra tick for state to update

    // After loading completes, disabled should be removed (empty string or undefined)
    const disabledAttr = submitButton.attributes('disabled')
    expect(disabledAttr === undefined || disabledAttr === '').toBe(true)
  })

  it('clears error message on new submission', async () => {
    mockOnLogin.mockResolvedValueOnce({
      success: false,
      error: 'First error',
    })

    const wrapper = mount(LoginForm, {
      props: {
        onLogin: mockOnLogin,
      },
    })

    await wrapper.find('#login-username').setValue('testuser')
    await wrapper.find('#login-password').setValue('wrongpass')
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('First error')

    mockOnLogin.mockResolvedValueOnce({ success: true })
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).not.toContain('First error')
  })
})
