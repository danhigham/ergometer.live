<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue'

const router = useRouter()
const authStore = useAuthStore()
const signingOut = ref(false)

const handleSignOut = async () => {
  try {
    signingOut.value = true
    await authStore.signOut()
    router.push('/login')
  } catch (error: any) {
    console.error('Sign out error:', error)
    alert('Failed to sign out: ' + error.message)
  } finally {
    signingOut.value = false
  }
}
</script>

<template>
  <main>
    <div class="home">
      <div class="header">
        <div class="user-info">
          <div v-if="authStore.isOnlineMode && authStore.user" class="user-details">
            <img
              v-if="authStore.user.photoURL"
              :src="authStore.user.photoURL"
              :alt="authStore.user.displayName || 'User'"
              class="user-avatar"
            />
            <div class="user-text">
              <span class="user-name">{{ authStore.user.displayName || authStore.user.email }}</span>
              <span class="mode-badge online">Online Mode</span>
            </div>
          </div>
          <div v-else class="user-details">
            <div class="user-avatar local">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </div>
            <div class="user-text">
              <span class="user-name">Local User</span>
              <span class="mode-badge local">Local Mode</span>
            </div>
          </div>
        </div>

        <button
          @click="handleSignOut"
          :disabled="signingOut"
          class="btn-sign-out"
        >
          {{ signingOut ? 'Signing out...' : 'Sign Out' }}
        </button>
      </div>

      <h1>Ergometer Live</h1>
      <p class="subtitle">Real-time PM5 workout tracking</p>

      <div class="nav-cards">
        <RouterLink to="/test" class="card">
          <h2>WebSocket Test</h2>
          <p>Test the connection to your local PM5 server</p>
        </RouterLink>
      </div>
    </div>
  </main>
</template>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 3rem 2rem;
  text-align: center;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding: 1rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.user-info {
  flex: 1;
}

.user-details {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.user-avatar.local {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
}

.user-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.25rem;
}

.user-name {
  font-weight: 600;
  color: #1f2937;
  font-size: 1rem;
}

.mode-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-weight: 500;
}

.mode-badge.online {
  background: #dbeafe;
  color: #1e40af;
}

.mode-badge.local {
  background: #f3f4f6;
  color: #4b5563;
}

.btn-sign-out {
  padding: 0.5rem 1.5rem;
  border: 2px solid #e5e7eb;
  border-radius: 6px;
  background: white;
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-sign-out:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
  color: #1f2937;
}

.btn-sign-out:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.subtitle {
  font-size: 1.25rem;
  color: #666;
  margin-bottom: 3rem;
}

.nav-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.card {
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  padding: 2rem;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
}

.card:hover {
  border-color: #3b82f6;
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.card h2 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  color: #1f2937;
}

.card p {
  color: #6b7280;
  margin: 0;
}
</style>
