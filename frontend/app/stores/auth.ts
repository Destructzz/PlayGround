import { useRuntimeConfig } from '#app'
import { defineStore } from 'pinia'
import { ref } from 'vue'

type AuthUser = {
  id: string
  email: string
  avatar_url: string
  name: string
  provider: string
}

type SessionResponse = {
  authenticated?: boolean
  user?: AuthUser | null
}

export const useAuthStore = defineStore('auth', () => {
  const config = useRuntimeConfig()
  const backendUrl = config.public.backendUrl.replace(/\/$/, '')
  const user = ref<AuthUser | null>(null)
  const isAuthenticated = ref(false)
  const isLoading = ref(true)

  async function fetchSession() {
    if (import.meta.server) {
      return
    }

    isLoading.value = true
    try {
      const data = await $fetch<SessionResponse>(`${backendUrl}/api/v1/auth/session`, {
        credentials: 'include'
      })

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
      await $fetch(`${backendUrl}/api/v1/auth/logout`, {
        method: 'POST',
        credentials: 'include'
      })
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
