<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo">
        <h1>Ergometer.Live</h1>
        <p class="subtitle">Track your rowing workouts in real-time</p>
      </div>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <div class="login-options">
        <button
          @click="handleGoogleSignIn"
          :disabled="loading"
          class="btn-google"
        >
          <svg class="google-icon" viewBox="0 0 24 24">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          {{ loading ? 'Signing in...' : 'Sign in with Google' }}
        </button>

        <div class="divider">
          <span>or</span>
        </div>

        <button
          @click="handleLocalMode"
          :disabled="loading"
          class="btn-local"
        >
          Continue in Local Mode
        </button>

        <p class="local-mode-info">
          Local mode stores all data in your browser. No account required, but data won't sync across devices.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref<string | null>(null)

const handleGoogleSignIn = async () => {
  try {
    loading.value = true
    error.value = null
    await authStore.signInWithGoogle()
    router.push('/')
  } catch (err: any) {
    error.value = err.message || 'Failed to sign in with Google'
  } finally {
    loading.value = false
  }
}

const handleLocalMode = () => {
  authStore.switchToLocalMode()
  router.push('/')
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
}

.login-card {
  background: white;
  border-radius: 12px;
  padding: 3rem;
  max-width: 450px;
  width: 100%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.logo {
  text-align: center;
  margin-bottom: 2rem;
}

.logo h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  color: #6b7280;
  font-size: 1rem;
  margin: 0;
}

.error-message {
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  padding: 0.75rem;
  margin-bottom: 1.5rem;
  color: #991b1b;
  font-size: 0.875rem;
}

.login-options {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-google {
  background: white;
  color: #1f2937;
  border: 2px solid #e5e7eb;
}

.btn-google:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
}

.google-icon {
  width: 20px;
  height: 20px;
}

.btn-local {
  background: #667eea;
  color: white;
}

.btn-local:hover:not(:disabled) {
  background: #5568d3;
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  color: #9ca3af;
  margin: 0.5rem 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid #e5e7eb;
}

.divider span {
  padding: 0 1rem;
  font-size: 0.875rem;
}

.local-mode-info {
  font-size: 0.75rem;
  color: #6b7280;
  text-align: center;
  margin: -0.5rem 0 0 0;
  line-height: 1.4;
}
</style>
