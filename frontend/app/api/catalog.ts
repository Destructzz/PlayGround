import { apiFetch } from './client'
import type {
  EventCatalogResponse,
  GamingCatalogResponse,
  HomeCatalogResponse,
  LoungeCatalogResponse,
  LoungeAvailabilityResponse,
  ShiftListResponse
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

export function getLoungeAvailability(zoneId: number, date: string) {
  return apiFetch<LoungeAvailabilityResponse>(
    `${PUBLIC_API_PREFIX}/lounge/${zoneId}/availability?date=${date}`
  )
}

export function getEventCatalog() {
  return apiFetch<EventCatalogResponse>(`${PUBLIC_API_PREFIX}/event`)
}

export function getShiftByZoneTagId(zoneTagId: number) {
  return apiFetch<ShiftListResponse>(`/api/v1/shift/zone-tag/${zoneTagId}`)
}
