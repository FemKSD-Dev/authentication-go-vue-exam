<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  onRegister: (username: string, password: string, confirmPassword: string) => Promise<{ success: boolean; error?: string }>
}

const props = defineProps<Props>()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)

// Validation errors
const usernameError = ref('')
const passwordError = ref('')
const confirmPasswordError = ref('')
const isRegistering = ref(false)

const emit = defineEmits<{
  switchToLogin: []
}>()

// Validation functions based on backend rules
const validateUsername = () => {
  usernameError.value = ''

  if (!username.value) {
    usernameError.value = 'Username is required'
    return false
  }

  if (username.value.length < 3) {
    usernameError.value = 'Username must be at least 3 characters'
    return false
  }

  if (username.value.length > 255) {
    usernameError.value = 'Username must not exceed 255 characters'
    return false
  }

  // Check alphanumeric (letters and numbers only)
  if (!/^[a-zA-Z0-9]+$/.test(username.value)) {
    usernameError.value = 'Username must contain only letters and numbers'
    return false
  }

  return true
}

const validatePassword = () => {
  passwordError.value = ''

  if (!password.value) {
    passwordError.value = 'Password is required'
    return false
  }

  if (password.value.length < 8) {
    passwordError.value = 'Password must be at least 8 characters'
    return false
  }

  if (password.value.length > 255) {
    passwordError.value = 'Password must not exceed 255 characters'
    return false
  }

  // Check ASCII characters only
  if (!/^[\x00-\x7F]+$/.test(password.value)) {
    passwordError.value = 'Password must contain only ASCII characters'
    return false
  }

  return true
}

const validateConfirmPassword = () => {
  confirmPasswordError.value = ''

  if (!confirmPassword.value) {
    confirmPasswordError.value = 'Please confirm your password'
    return false
  }

  if (password.value !== confirmPassword.value) {
    confirmPasswordError.value = 'Passwords do not match'
    return false
  }

  return true
}

const handleRegister = async () => {
  // Validate all fields
  const isUsernameValid = validateUsername()
  const isPasswordValid = validatePassword()
  const isConfirmPasswordValid = validateConfirmPassword()

  if (!isUsernameValid || !isPasswordValid || !isConfirmPasswordValid) {
    return
  }

  try {
    isRegistering.value = true
    
    const result = await props.onRegister(
      username.value,
      password.value,
      confirmPassword.value
    )

    if (result.success) {
      // Clear form on success
      username.value = ''
      password.value = ''
      confirmPassword.value = ''
      usernameError.value = ''
      passwordError.value = ''
      confirmPasswordError.value = ''
    } else if (result.error) {
      // Show error from server
      usernameError.value = result.error
    }
  } finally {
    isRegistering.value = false
  }
}

const handleSwitchToLogin = (e: Event) => {
  e.preventDefault()
  emit('switchToLogin')
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

const toggleConfirmPasswordVisibility = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}

// Check if field has error
const hasUsernameError = computed(() => usernameError.value !== '')
const hasPasswordError = computed(() => passwordError.value !== '')
const hasConfirmPasswordError = computed(() => confirmPasswordError.value !== '')

// Check if form is valid for submission
const isFormValid = computed(() => {
  // All fields must have values
  if (!username.value || !password.value || !confirmPassword.value) {
    return false
  }

  // Username validation
  if (username.value.length < 3 || username.value.length > 255) {
    return false
  }
  if (!/^[a-zA-Z0-9]+$/.test(username.value)) {
    return false
  }

  // Password validation
  if (password.value.length < 8 || password.value.length > 255) {
    return false
  }
  if (!/^[\x00-\x7F]+$/.test(password.value)) {
    return false
  }

  // Confirm password validation
  if (password.value !== confirmPassword.value) {
    return false
  }

  return true
})
</script>

<template>
  <div class="register-form-container">
    <form @submit.prevent="handleRegister">
      <!-- Username Field -->
      <div class="form-group">
        <label for="register-username" class="form-label">Username</label>
        <div class="input-wrapper">
          <input
            id="register-username"
            v-model="username"
            type="text"
            :class="['form-input', { 'input-error': hasUsernameError }]"
            placeholder="Enter your username"
            @blur="validateUsername"
            @input="usernameError = ''"
          />
          <div v-if="hasUsernameError" class="error-message">
            {{ usernameError }}
          </div>
        </div>
      </div>

      <!-- Password Field -->
      <div class="form-group">
        <label for="register-password" class="form-label">Password</label>
        <div class="input-wrapper">
          <div class="password-input-wrapper">
            <input
              id="register-password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              :class="['form-input', { 'input-error': hasPasswordError }]"
              placeholder="Enter your password"
              @blur="validatePassword"
              @input="passwordError = ''"
            />
            <button
              type="button"
              @click="togglePasswordVisibility"
              class="toggle-password-btn"
              :aria-label="showPassword ? 'Hide password' : 'Show password'"
            >
              <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="icon">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="icon">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
          <div v-if="hasPasswordError" class="error-message">
            {{ passwordError }}
          </div>
        </div>
      </div>

      <!-- Confirm Password Field -->
      <div class="form-group">
        <label for="register-confirm-password" class="form-label">Confirm Password</label>
        <div class="input-wrapper">
          <div class="password-input-wrapper">
            <input
              id="register-confirm-password"
              v-model="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              :class="['form-input', { 'input-error': hasConfirmPasswordError }]"
              placeholder="Confirm your password"
              @blur="validateConfirmPassword"
              @input="confirmPasswordError = ''"
            />
            <button
              type="button"
              @click="toggleConfirmPasswordVisibility"
              class="toggle-password-btn"
              :aria-label="showConfirmPassword ? 'Hide password' : 'Show password'"
            >
              <svg v-if="!showConfirmPassword" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="icon">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="icon">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
          <div v-if="hasConfirmPasswordError" class="error-message">
            {{ confirmPasswordError }}
          </div>
        </div>
      </div>

      <button
        type="submit"
        class="submit-button"
        :disabled="!isFormValid"
        :class="{ 'button-disabled': !isFormValid }"
      >
        สมัครสมาชิก
      </button>

      <div class="footer-link">
        <a href="#" class="link" @click="handleSwitchToLogin">กลับสู่หน้าลงชื่อเข้าสู่ระบบ</a>
      </div>
    </form>
  </div>
</template>

<style scoped>
.register-form-container {
  width: 100%;
}

form {
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
}

.form-group {
  display: flex;
  align-items: flex-start;
  gap: 1.25rem;
}

.form-label {
  min-width: 150px;
  text-align: right;
  font-size: 1rem;
  color: #1e293b;
  font-weight: 600;
  padding-top: 0.875rem;
}

.input-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.password-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
}

.form-input {
  width: 100%;
  padding: 0.875rem 1.25rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.75rem;
  font-size: 1rem;
  transition: all 0.2s ease;
  outline: none;
  background: #f8fafc;
  color: #1e293b;
}

.password-input-wrapper .form-input {
  padding-right: 3rem;
}

.form-input::placeholder {
  color: #94a3b8;
}

.form-input:hover {
  border-color: #cbd5e1;
  background: #ffffff;
}

.form-input:focus {
  border-color: #6366f1;
  background: #ffffff;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1);
}

.form-input.input-error {
  border-color: #ef4444;
  background: #fef2f2;
}

.form-input.input-error:focus {
  border-color: #ef4444;
  box-shadow: 0 0 0 4px rgba(239, 68, 68, 0.1);
}

.error-message {
  font-size: 0.875rem;
  color: #ef4444;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.error-message::before {
  content: '⚠';
  font-size: 1rem;
}

.toggle-password-btn {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  transition: color 0.2s ease;
  border-radius: 0.375rem;
}

.toggle-password-btn:hover {
  color: #475569;
  background: #f1f5f9;
}

.toggle-password-btn:active {
  color: #334155;
}

.toggle-password-btn .icon {
  width: 20px;
  height: 20px;
}

.submit-button {
  width: 100%;
  padding: 1rem 2rem;
  margin-top: 0.5rem;
  background: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  color: white;
  border: none;
  border-radius: 0.75rem;
  font-size: 1.05rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.4);
  letter-spacing: 0.02em;
}

.submit-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.5);
}

.submit-button:active:not(:disabled) {
  transform: translateY(0);
}

.submit-button:disabled,
.submit-button.button-disabled {
  background: linear-gradient(135deg, #cbd5e1 0%, #94a3b8 100%);
  cursor: not-allowed;
  box-shadow: none;
  opacity: 0.6;
}

.submit-button:disabled:hover,
.submit-button.button-disabled:hover {
  transform: none;
  box-shadow: none;
}

.footer-link {
  text-align: center;
  margin-top: -0.5rem;
}

.link {
  color: #6366f1;
  font-size: 0.95rem;
  font-weight: 600;
  text-decoration: none;
  transition: color 0.2s ease;
}

.link:hover {
  color: #4f46e5;
  text-decoration: underline;
}

@media (max-width: 768px) {
  .form-label {
    min-width: 130px;
    font-size: 0.95rem;
  }

  .form-input {
    padding: 0.75rem 1rem;
    font-size: 0.95rem;
  }

  .submit-button {
    padding: 0.875rem 1.5rem;
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .form-group {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .form-label {
    min-width: auto;
    text-align: left;
    padding-top: 0;
  }
}
</style>
