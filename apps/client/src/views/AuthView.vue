<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import { authService, type LoginResponse } from '@/services'
import type { AxiosError } from 'axios'

type TabType = 'login' | 'register'

const router = useRouter()
const activeTab = ref<TabType>('login')

const switchTab = (tab: TabType) => {
  activeTab.value = tab
}

const handleLogin = async (username: string, password: string): Promise<{ success: boolean; error?: string }> => {
  try {
    const response: LoginResponse = await authService.login({
      username,
      password,
    })

    if (response.data.access_token) {
      localStorage.setItem('token', response.data.access_token)

      // Redirect to getting-started page
      router.push('/getting-started')
      return { success: true }
    }

    return { success: false, error: 'ไม่ได้รับ token จากเซิร์ฟเวอร์' }
  } catch (error) {
    console.error('Login error:', error)
    if (error instanceof Error && 'response' in error) {
      const axError = error as AxiosError<{ message?: string }>
      return {
        success: false,
        error: axError.response?.data?.message || 'เกิดข้อผิดพลาดในการเข้าสู่ระบบ กรุณาลองใหม่อีกครั้ง'
      }
    }
    return { success: false, error: 'เกิดข้อผิดพลาดในการเข้าสู่ระบบ กรุณาลองใหม่อีกครั้ง' }
  }
}

const handleRegister = async (
  username: string,
  password: string,
  confirmPassword: string
): Promise<{ success: boolean; error?: string }> => {
  try {
    await authService.register({
      username,
      password,
      confirmPassword,
    })

    // After successful registration, switch to login tab
    switchTab('login')

    return { success: true }
  } catch (error) {
    console.error('Registration error:', error)
    if (error instanceof Error && 'response' in error) {
      const axError = error as AxiosError<{ message?: string }>
      return {
        success: false,
        error: axError.response?.data?.message || 'เกิดข้อผิดพลาดในการสมัครสมาชิก กรุณาลองใหม่อีกครั้ง',
      }
    }
    return { success: false, error: 'เกิดข้อผิดพลาดในการสมัครสมาชิก กรุณาลองใหม่อีกครั้ง' }
  }
}
</script>

<template>
  <div class="auth-view">
    <div class="auth-container">
      <div class="auth-card">
        <div class="card-header">
          <h1 class="brand-title">Authentication</h1>
          <p class="brand-subtitle">Welcome back! Please login or register to continue</p>
        </div>

        <div class="tabs-container">
          <button
            :class="['tab', { active: activeTab === 'login' }]"
            @click="switchTab('login')"
          >
            <svg class="tab-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
            </svg>
            <span class="tab-text">ลงชื่อเข้าสู่ระบบ</span>
          </button>
          <button
            :class="['tab', { active: activeTab === 'register' }]"
            @click="switchTab('register')"
          >
            <svg class="tab-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
            </svg>
            <span class="tab-text">สมัครสมาชิก</span>
          </button>
        </div>

        <div class="form-container">
          <Transition name="form-switch" mode="out-in">
            <LoginForm
              v-if="activeTab === 'login'"
              key="login"
              :on-login="handleLogin"
              @switch-to-register="switchTab('register')"
            />
            <RegisterForm
              v-else
              key="register"
              :on-register="handleRegister"
              @switch-to-login="switchTab('login')"
            />
          </Transition>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
  padding: 2rem;
  position: relative;
  overflow: hidden;
}

.auth-view::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.1) 0%, transparent 70%);
  animation: float 20s ease-in-out infinite;
}

.auth-view::after {
  content: '';
  position: absolute;
  bottom: -50%;
  left: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, transparent 70%);
  animation: float 15s ease-in-out infinite reverse;
}

.auth-container {
  width: 100%;
  max-width: 1000px;
  position: relative;
  z-index: 1;
  animation: fadeInUp 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.auth-card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 2rem;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  padding: 3rem 2.5rem 2rem;
  text-align: center;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e2e8f0;
  animation: slideDown 0.5s ease-out 0.15s both;
}

.brand-title {
  font-size: 2.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
  animation: fadeInScale 0.5s ease-out 0.3s both;
}

.brand-subtitle {
  font-size: 1rem;
  color: #64748b;
  margin: 0;
  font-weight: 500;
  animation: fadeIn 0.5s ease-out 0.45s both;
}

.tabs-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  background: #f8fafc;
  padding: 1rem 3rem;
  animation: fadeIn 0.5s ease-out 0.55s both;
}

.tab {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  background: transparent;
  border: 2px solid #e2e8f0;
  border-radius: 0.75rem;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  color: #64748b;
  transition: all 0.2s ease;
}

.tab-icon {
  width: 20px;
  height: 20px;
  transition: transform 0.2s ease;
}

.tab:hover:not(.active) {
  background: #f1f5f9;
  color: #475569;
  border-color: #e2e8f0;
}

.tab:hover .tab-icon {
  transform: scale(1.1);
}

.tab.active {
  background: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  color: white;
  border-color: #6366f1;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.4);
}

.tab.active .tab-icon {
  transform: scale(1.05);
}

.form-container {
  padding: 3rem 3rem 3.5rem;
  background: white;
  position: relative;
  animation: fadeInUp 0.5s ease-out 0.65s both;
}

.form-switch-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.form-switch-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.form-switch-enter-from {
  opacity: 0;
  transform: translateX(40px) scale(0.95);
}

.form-switch-leave-to {
  opacity: 0;
  transform: translateX(-40px) scale(0.95);
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -30px) rotate(120deg);
  }
  66% {
    transform: translate(-20px, 20px) rotate(240deg);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInScale {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@media (max-width: 768px) {
  .auth-view {
    padding: 1rem;
  }

  .card-header {
    padding: 2rem 1.5rem 1.5rem;
  }

  .brand-title {
    font-size: 2rem;
  }

  .brand-subtitle {
    font-size: 0.9rem;
  }

  .form-container {
    padding: 2rem 1.5rem 2.5rem;
  }

  .tab {
    padding: 0.875rem 1rem;
    font-size: 0.95rem;
  }

  .tab-icon {
    width: 18px;
    height: 18px;
  }
}

@media (max-width: 480px) {
  .card-header {
    padding: 1.5rem 1rem 1rem;
  }

  .brand-title {
    font-size: 1.75rem;
  }

  .form-container {
    padding: 1.5rem 1rem 2rem;
  }

  .tab {
    flex-direction: column;
    gap: 0.25rem;
    padding: 0.75rem 0.5rem;
    font-size: 0.875rem;
  }

  .tab-text {
    font-size: 0.875rem;
  }
}
</style>
