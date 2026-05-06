import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchAuthSession, logoutSession } from '~/api/auth'
import type { AuthUser } from '~/api/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(null)
  const isAuthenticated = ref(false)
  const isLoading = ref(true)

  async function fetchSession() {
    if (import.meta.server) {
      return
    }

    isLoading.value = true
    try {
      const data = await fetchAuthSession()

      if (data.authenticated) {
        isAuthenticated.value = true
        user.value = data.user ?? null
      } else {
        isAuthenticated.value = false
        user.value = null
      }
    } catch (e) {
      console.error('Failed to fetch session', e)
      isAuthenticated.value = false
      user.value = null
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      await logoutSession()
      isAuthenticated.value = false
      user.value = null
    } catch (e) {
      console.error('Logout failed', e)
    }
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    fetchSession,
    logout
  }
})
