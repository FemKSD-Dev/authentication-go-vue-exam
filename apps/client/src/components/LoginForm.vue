<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  onLogin: (username: string, password: string) => Promise<{ success: boolean; error?: string }>
}

const props = defineProps<Props>()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

const emit = defineEmits<{
  switchToRegister: []
}>()

const handleLogin = async () => {
  if (!username.value || !password.value) {
    errorMessage.value = 'กรุณากรอกชื่อผู้ใช้และรหัสผ่าน'
    return
  }

  try {
    isLoading.value = true
    errorMessage.value = ''

    const result = await props.onLogin(username.value, password.value)

    if (!result.success && result.error) {
      errorMessage.value = result.error
    }
  } finally {
    isLoading.value = false
  }
}

const handleSwitchToRegister = (e: Event) => {
  e.preventDefault()
  emit('switchToRegister')
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}
</script>

<template>
  <div class="login-form-container">
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="login-username" class="form-label">Username</label>
        <input
          id="login-username"
          v-model="username"
          type="text"
          class="form-input"
          placeholder="Enter your username"
          required
        />
      </div>

      <div class="form-group">
        <label for="login-password" class="form-label">Password</label>
        <div class="password-input-wrapper">
          <input
            id="login-password"
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            class="form-input"
            placeholder="Enter your password"
            required
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
      </div>

      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>

      <button type="submit" class="submit-button" :disabled="isLoading">
        <span v-if="!isLoading">ลงชื่อเข้าใช้งาน</span>
        <span v-else>กำลังเข้าสู่ระบบ...</span>
      </button>

      <div class="footer-link">
        <a href="#" class="link" @click="handleSwitchToRegister">สมัครสมาชิก</a>
      </div>
    </form>
  </div>
</template>

<style scoped>
.login-form-container {
  width: 100%;
}

form {
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
}

.form-group {
  display: flex;
  align-items: center;
  gap: 1.25rem;
}

.form-label {
  min-width: 120px;
  text-align: right;
  font-size: 1rem;
  color: #1e293b;
  font-weight: 600;
}

.password-input-wrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}

.form-input {
  flex: 1;
  padding: 0.875rem 1.25rem;
  padding-right: 3rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.75rem;
  font-size: 1rem;
  transition: all 0.2s ease;
  outline: none;
  background: #f8fafc;
  color: #1e293b;
  width: 100%;
}

.form-input::placeholder {
  color: #94a3b8;
}

.form-input:hover {
  border-color: #cbd5e1;
  background: #ffffff;
}

.form-input:focus {
  border-color: #10b981;
  background: #ffffff;
  box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.1);
}

.error-message {
  padding: 0.875rem 1.25rem;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.75rem;
  color: #dc2626;
  font-size: 0.95rem;
  font-weight: 500;
  text-align: center;
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

.submit-button:active {
  transform: translateY(0);
}

.submit-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.footer-link {
  text-align: center;
  margin-top: -0.5rem;
}

.link {
  color: #3b82f6;
  font-size: 0.95rem;
  font-weight: 600;
  text-decoration: none;
  transition: color 0.2s ease;
}

.link:hover {
  color: #2563eb;
  text-decoration: underline;
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

@media (max-width: 768px) {
  .form-label {
    min-width: 100px;
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
  }
}
</style>
