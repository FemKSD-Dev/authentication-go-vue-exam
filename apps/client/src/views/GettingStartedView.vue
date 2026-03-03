<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import UserDetailCard from '@/components/UserDetailCard.vue'
import { authService } from '@/services'

const router = useRouter()
const showCard = ref(false)
const username = ref('')
const userId = ref('')
const isLoading = ref(true)

const fetchUserData = async () => {
  try {
    const response = await authService.me()
    username.value = response.data.username
    userId.value = response.data.user_id
    
    setTimeout(() => {
      showCard.value = true
    }, 100)
  } catch (error) {
    console.error('Failed to fetch user data:', error)
    router.push('/')
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchUserData()
})
</script>

<template>
  <div class="getting-started-view">
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p class="loading-text">กำลังโหลดข้อมูล...</p>
    </div>
    <div v-else class="content-container">
      <Transition name="card-appear">
        <UserDetailCard v-if="showCard" :username="username" />
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.getting-started-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}

.getting-started-view::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.15) 0%, transparent 70%);
  animation: float 20s ease-in-out infinite;
}

.getting-started-view::after {
  content: '';
  position: absolute;
  bottom: -50%;
  left: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.15) 0%, transparent 70%);
  animation: float 15s ease-in-out infinite reverse;
}

.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1.5rem;
  position: relative;
  z-index: 1;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-text {
  color: rgba(255, 255, 255, 0.8);
  font-size: 1.1rem;
  font-weight: 500;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.content-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

.card-appear-enter-active {
  transition: all 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.card-appear-enter-from {
  opacity: 0;
  transform: translateY(60px) scale(0.9) rotateX(10deg);
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

@media (max-width: 768px) {
  .content-container {
    padding: 1.5rem;
  }
}

@media (max-width: 480px) {
  .content-container {
    padding: 1rem;
  }
}
</style>
