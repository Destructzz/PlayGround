import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useCatalogStore = defineStore('catalog', () => {
  const homeData = ref<any>(null)
  const gamingData = ref<any>(null)
  const loungeData = ref<any>(null)
  const eventData = ref<any>(null)
  const isLoading = ref(false)

  async function fetchHome() {
    isLoading.value = true
    try {
      const { data } = await $fetch<any>('/api/v1/public/home')
      homeData.value = data
    } catch (e) {
      console.error('Failed to fetch home catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchGaming() {
    isLoading.value = true
    try {
      const { data } = await $fetch<any>('/api/v1/public/gaming')
      gamingData.value = data
    } catch (e) {
      console.error('Failed to fetch gaming catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchLounge() {
    isLoading.value = true
    try {
      const { data } = await $fetch<any>('/api/v1/public/lounge')
      loungeData.value = data
    } catch (e) {
      console.error('Failed to fetch lounge catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchEvent() {
    isLoading.value = true
    try {
      const { data } = await $fetch<any>('/api/v1/public/event')
      eventData.value = data
    } catch (e) {
      console.error('Failed to fetch event catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  return {
    homeData,
    gamingData,
    loungeData,
    eventData,
    isLoading,
    fetchHome,
    fetchGaming,
    fetchLounge,
    fetchEvent
  }
})
