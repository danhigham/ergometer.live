import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { signInWithPopup, signOut as firebaseSignOut, onAuthStateChanged, type User } from 'firebase/auth'
import { auth, googleProvider } from '@/services/firebase'

export type AppMode = 'online' | 'local'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const mode = ref<AppMode>('online')
  const loading = ref(true)
  const error = ref<string | null>(null)

  // Computed
  const isAuthenticated = computed(() => {
    if (mode.value === 'local') return true
    return user.value !== null
  })

  const isOnlineMode = computed(() => mode.value === 'online')
  const isLocalMode = computed(() => mode.value === 'local')

  // Actions
  const signInWithGoogle = async () => {
    try {
      error.value = null
      const result = await signInWithPopup(auth, googleProvider)
      user.value = result.user
      mode.value = 'online'

      // Store mode preference
      localStorage.setItem('app_mode', 'online')

      return result.user
    } catch (err: any) {
      error.value = err.message
      console.error('Sign in error:', err)
      throw err
    }
  }

  const signOut = async () => {
    try {
      error.value = null
      if (mode.value === 'online' && user.value) {
        await firebaseSignOut(auth)
      }
      user.value = null
      mode.value = 'online'
      localStorage.removeItem('app_mode')
    } catch (err: any) {
      error.value = err.message
      console.error('Sign out error:', err)
      throw err
    }
  }

  const switchToLocalMode = () => {
    mode.value = 'local'
    localStorage.setItem('app_mode', 'local')
  }

  const switchToOnlineMode = () => {
    mode.value = 'online'
    localStorage.setItem('app_mode', 'online')
  }

  // Initialize auth state listener
  const initializeAuth = () => {
    // Check for saved mode preference
    const savedMode = localStorage.getItem('app_mode') as AppMode | null
    if (savedMode === 'local') {
      mode.value = 'local'
      loading.value = false
      return
    }

    // Listen for auth state changes in online mode
    onAuthStateChanged(auth, (firebaseUser) => {
      user.value = firebaseUser
      loading.value = false
    })
  }

  return {
    // State
    user,
    mode,
    loading,
    error,
    // Computed
    isAuthenticated,
    isOnlineMode,
    isLocalMode,
    // Actions
    signInWithGoogle,
    signOut,
    switchToLocalMode,
    switchToOnlineMode,
    initializeAuth,
  }
})
