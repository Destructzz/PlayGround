import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server) return

  const authStore = useAuthStore()

  if (!authStore.user && !authStore.isAuthenticated) {
    await authStore.fetchSession()
  }

  if (!authStore.isAuthenticated || authStore.user?.role !== 'admin') {
    return navigateTo('/')
  }
})
