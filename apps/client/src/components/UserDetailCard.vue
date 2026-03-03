<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  username?: string
}>()

const displayUsername = ref(props.username || 'xxx')
</script>

<template>
  <div class="user-card">
    <div class="card-header">
      <div class="avatar-container">
        <div class="avatar">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="avatar-icon"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
            />
          </svg>
        </div>
        <div class="status-indicator"></div>
      </div>
    </div>

    <div class="card-body">
      <h1 class="welcome-title">Welcome User</h1>
      <div class="username-display">
        <span class="username-label">:</span>
        <span class="username">{{ displayUsername }}</span>
      </div>
    </div>

    <div class="card-footer">
      <div class="info-badge">
        <svg class="badge-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>Authenticated</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.user-card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 2rem;
  padding: 0;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  max-width: 700px;
  width: 100%;
  margin: 0 auto;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  background: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  padding: 3rem 2rem 2rem;
  display: flex;
  justify-content: center;
  position: relative;
  animation: slideDown 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) 0.2s both;
}

.card-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  opacity: 0.3;
}

.avatar-container {
  position: relative;
}

.avatar {
  width: 120px;
  height: 120px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  position: relative;
  z-index: 1;
  animation: zoomIn 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) 0.4s both;
}

.avatar-icon {
  width: 64px;
  height: 64px;
  color: #10b981;
}

.status-indicator {
  position: absolute;
  bottom: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  background: #10b981;
  border: 4px solid white;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.5);
  z-index: 2;
  animation: popIn 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55) 0.7s both;
}

.card-body {
  padding: 2.5rem 2rem;
  text-align: center;
  background: white;
  animation: fadeInUp 0.6s ease-out 0.5s both;
}

.welcome-title {
  font-size: 2rem;
  color: #1e293b;
  font-weight: 700;
  margin: 0 0 1rem 0;
  letter-spacing: -0.02em;
  animation: fadeInScale 0.5s ease-out 0.6s both;
}

.username-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1rem 2rem;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 1rem;
  border: 2px solid #e2e8f0;
  animation: fadeInScale 0.5s ease-out 0.75s both;
}

.username-label {
  font-size: 1.5rem;
  color: #64748b;
  font-weight: 600;
}

.username {
  font-size: 1.75rem;
  color: #6366f1;
  font-weight: 800;
  letter-spacing: 0.02em;
}

.card-footer {
  padding: 1.5rem 2rem;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: center;
  animation: fadeInUp 0.6s ease-out 0.7s both;
}

.info-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  color: white;
  border-radius: 2rem;
  font-size: 0.95rem;
  font-weight: 600;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.4);
  animation: bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55) 0.8s both;
}

.badge-icon {
  width: 20px;
  height: 20px;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes zoomIn {
  from {
    opacity: 0;
    transform: scale(0);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes popIn {
  0% {
    opacity: 0;
    transform: scale(0);
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
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

@keyframes bounceIn {
  0% {
    opacity: 0;
    transform: scale(0.3) translateY(20px);
  }
  50% {
    transform: scale(1.05) translateY(-5px);
  }
  70% {
    transform: scale(0.95) translateY(0);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@media (max-width: 768px) {
  .card-header {
    padding: 2.5rem 1.5rem 1.5rem;
  }

  .avatar {
    width: 100px;
    height: 100px;
  }

  .avatar-icon {
    width: 52px;
    height: 52px;
  }

  .status-indicator {
    width: 20px;
    height: 20px;
    bottom: 6px;
    right: 6px;
  }

  .card-body {
    padding: 2rem 1.5rem;
  }

  .welcome-title {
    font-size: 1.75rem;
  }

  .username {
    font-size: 1.5rem;
  }

  .username-label {
    font-size: 1.25rem;
  }
}

@media (max-width: 480px) {
  .card-header {
    padding: 2rem 1rem 1rem;
  }

  .avatar {
    width: 80px;
    height: 80px;
  }

  .avatar-icon {
    width: 44px;
    height: 44px;
  }

  .card-body {
    padding: 1.5rem 1rem;
  }

  .welcome-title {
    font-size: 1.5rem;
  }

  .username-display {
    padding: 0.875rem 1.5rem;
  }

  .username {
    font-size: 1.25rem;
  }

  .username-label {
    font-size: 1.125rem;
  }

  .card-footer {
    padding: 1.25rem 1rem;
  }

  .info-badge {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
  }
}
</style>
