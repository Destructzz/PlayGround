import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  if (!authStore.user && !authStore.isAuthenticated) {
    await authStore.fetchSession()
  }

  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }
})
