import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  getEventCatalog,
  getGamingCatalog,
  getHomeCatalog,
  getLoungeCatalog
} from '~/api/catalog'
import type {
  EventCatalogResponse,
  GamingCatalogResponse,
  HomeCatalogResponse,
  LoungeCatalogResponse
} from '~/api/types'

export const useCatalogStore = defineStore('catalog', () => {
  const homeData = ref<HomeCatalogResponse | null>(null)
  const gamingData = ref<GamingCatalogResponse | null>(null)
  const loungeData = ref<LoungeCatalogResponse | null>(null)
  const eventData = ref<EventCatalogResponse | null>(null)
  const isLoading = ref(false)

  async function fetchHome() {
    isLoading.value = true
    try {
      homeData.value = await getHomeCatalog()
    } catch (e) {
      console.error('Failed to fetch home catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchGaming() {
    isLoading.value = true
    try {
      gamingData.value = await getGamingCatalog()
    } catch (e) {
      console.error('Failed to fetch gaming catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchLounge() {
    isLoading.value = true
    try {
      loungeData.value = await getLoungeCatalog()
    } catch (e) {
      console.error('Failed to fetch lounge catalog', e)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchEvent() {
    isLoading.value = true
    try {
      eventData.value = await getEventCatalog()
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
