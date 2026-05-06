import { apiFetch } from './client'
import type {
  EventCatalogResponse,
  GamingCatalogResponse,
  HomeCatalogResponse,
  LoungeCatalogResponse
} from './types'

const PUBLIC_API_PREFIX = '/api/v1/public'

export function getHomeCatalog() {
  return apiFetch<HomeCatalogResponse>(`${PUBLIC_API_PREFIX}/home`)
}

export function getGamingCatalog() {
  return apiFetch<GamingCatalogResponse>(`${PUBLIC_API_PREFIX}/gaming`)
}

export function getLoungeCatalog() {
  return apiFetch<LoungeCatalogResponse>(`${PUBLIC_API_PREFIX}/lounge`)
}

export function getEventCatalog() {
  return apiFetch<EventCatalogResponse>(`${PUBLIC_API_PREFIX}/event`)
}
