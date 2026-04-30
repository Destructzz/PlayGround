import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(async () => {
  if (import.meta.server) {
    return
  }

  const authStore = useAuthStore()
  await authStore.fetchSession()
})
