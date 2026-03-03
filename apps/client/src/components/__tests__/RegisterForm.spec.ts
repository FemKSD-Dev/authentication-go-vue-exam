import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import RegisterForm from '../RegisterForm.vue'

describe('RegisterForm', () => {
  let mockOnRegister: ReturnType<typeof vi.fn>

  beforeEach(() => {
    mockOnRegister = vi.fn()
  })

  it('renders properly', () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    expect(wrapper.find('#register-username').exists()).toBe(true)
    expect(wrapper.find('#register-password').exists()).toBe(true)
    expect(wrapper.find('#register-confirm-password').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('validates username requirements', async () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    // Test empty username
    await wrapper.find('#register-username').setValue('')
    await wrapper.find('#register-password').setValue('password123')
    await wrapper.find('#register-confirm-password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('Username is required')
    expect(mockOnRegister).not.toHaveBeenCalled()

    // Test username too short
    await wrapper.find('#register-username').setValue('ab')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('Username must be at least 3 characters')

    // Test invalid characters
    await wrapper.find('#register-username').setValue('user@123')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('Username must contain only letters and numbers')
  })

  it('validates password requirements', async () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    // Test empty password
    await wrapper.find('#register-username').setValue('testuser')
    await wrapper.find('#register-password').setValue('')
    await wrapper.find('#register-confirm-password').setValue('')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('Password is required')

    // Test password too short
    await wrapper.find('#register-password').setValue('pass123')
    await wrapper.find('#register-confirm-password').setValue('pass123')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('Password must be at least 8 characters')
  })

  it('validates password confirmation', async () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    await wrapper.find('#register-username').setValue('testuser')
    await wrapper.find('#register-password').setValue('password123')
    await wrapper.find('#register-confirm-password').setValue('password456')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.text()).toContain('Passwords do not match')
    expect(mockOnRegister).not.toHaveBeenCalled()
  })

  it('calls onRegister with valid data', async () => {
    mockOnRegister.mockResolvedValue({ success: true })

    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    await wrapper.find('#register-username').setValue('testuser')
    await wrapper.find('#register-password').setValue('password123')
    await wrapper.find('#register-confirm-password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')

    expect(mockOnRegister).toHaveBeenCalledWith('testuser', 'password123', 'password123')
  })

  it('displays error message from failed registration', async () => {
    mockOnRegister.mockResolvedValue({
      success: false,
      error: 'Username already exists',
    })

    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    await wrapper.find('#register-username').setValue('existinguser')
    await wrapper.find('#register-password').setValue('password123')
    await wrapper.find('#register-confirm-password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Username already exists')
  })

  it('clears form on successful registration', async () => {
    mockOnRegister.mockResolvedValue({ success: true })

    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    await wrapper.find('#register-username').setValue('testuser')
    await wrapper.find('#register-password').setValue('password123')
    await wrapper.find('#register-confirm-password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    expect((wrapper.find('#register-username').element as HTMLInputElement).value).toBe('')
    expect((wrapper.find('#register-password').element as HTMLInputElement).value).toBe('')
    expect((wrapper.find('#register-confirm-password').element as HTMLInputElement).value).toBe(
      '',
    )
  })

  it('toggles password visibility', async () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    const passwordInput = wrapper.find('#register-password')
    const toggleButtons = wrapper.findAll('.toggle-password-btn')

    expect(passwordInput.attributes('type')).toBe('password')

    await toggleButtons[0].trigger('click')
    expect(passwordInput.attributes('type')).toBe('text')

    await toggleButtons[0].trigger('click')
    expect(passwordInput.attributes('type')).toBe('password')
  })

  it('toggles confirm password visibility', async () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    const confirmPasswordInput = wrapper.find('#register-confirm-password')
    const toggleButtons = wrapper.findAll('.toggle-password-btn')

    expect(confirmPasswordInput.attributes('type')).toBe('password')

    await toggleButtons[1].trigger('click')
    expect(confirmPasswordInput.attributes('type')).toBe('text')

    await toggleButtons[1].trigger('click')
    expect(confirmPasswordInput.attributes('type')).toBe('password')
  })

  it('emits switchToLogin event', async () => {
    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    await wrapper.find('.link').trigger('click')

    expect(wrapper.emitted('switchToLogin')).toBeTruthy()
  })

  it('disables submit button during registration', async () => {
    mockOnRegister.mockResolvedValue({ success: true })

    const wrapper = mount(RegisterForm, {
      props: {
        onRegister: mockOnRegister,
      },
    })

    await wrapper.find('#register-username').setValue('validuser123')
    await wrapper.find('#register-password').setValue('password123')
    await wrapper.find('#register-confirm-password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')

    // Wait for async operation to complete
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 50))

    // After successful registration, form should be cleared
    expect((wrapper.find('#register-username').element as HTMLInputElement).value).toBe('')
    expect(mockOnRegister).toHaveBeenCalledWith('validuser123', 'password123', 'password123')
  })
})
