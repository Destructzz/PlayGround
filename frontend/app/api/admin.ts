import { apiFetch } from './client'
import type {
  AdminBookingListResponse,
  AdminBookingResponse,
  AdminPlaceListResponse,
  AdminPlaceResponse,
  AdminServiceListResponse,
  AdminZoneListResponse,
  AdminZoneResponse,
  AdminZoneTagListResponse,
  AdminZoneTagResponse,
  PatchBookingPayload,
  CreatePlacePayload,
  CreateShiftPayload,
  CreateZonePayload,
  CreateZoneTagPayload,
  PatchPlacePayload,
  PatchShiftPayload,
  PatchZonePayload,
  PatchZoneTagPayload,
  ShiftListResponse,
  CreateServicePayload,
  PatchServicePayload,
  AdminServiceResponse,
  AdminStatsResponse
} from './types'


const withAuth = {
  credentials: 'include' as const
}

export function getAdminZoneTags() {
  return apiFetch<AdminZoneTagListResponse>('/api/v1/zone-tag', withAuth)
}

export function createAdminZoneTag(payload: CreateZoneTagPayload) {
  return apiFetch<AdminZoneTagResponse>('/api/v1/zone-tag', {
    ...withAuth,
    method: 'POST',
    body: payload
  })
}

export function patchAdminZoneTag(id: number, payload: PatchZoneTagPayload) {
  return apiFetch<AdminZoneTagResponse>(`/api/v1/zone-tag/${id}`, {
    ...withAuth,
    method: 'PATCH',
    body: payload
  })
}

export function deleteAdminZoneTag(id: number) {
  return apiFetch(`/api/v1/zone-tag/${id}`, {
    ...withAuth,
    method: 'DELETE'
  })
}

export function getAdminZones() {
  return apiFetch<AdminZoneListResponse>('/api/v1/zone', withAuth)
}

export function createAdminZone(payload: CreateZonePayload) {
  return apiFetch<AdminZoneResponse>('/api/v1/zone', {
    ...withAuth,
    method: 'POST',
    body: payload
  })
}

export function patchAdminZone(id: number, payload: PatchZonePayload) {
  return apiFetch<AdminZoneResponse>(`/api/v1/zone/${id}`, {
    ...withAuth,
    method: 'PATCH',
    body: payload
  })
}

export function deleteAdminZone(id: number) {
  return apiFetch(`/api/v1/zone/${id}`, {
    ...withAuth,
    method: 'DELETE'
  })
}

export function getAdminPlaces() {
  return apiFetch<AdminPlaceListResponse>('/api/v1/place', withAuth)
}

export function createAdminPlace(payload: CreatePlacePayload) {
  return apiFetch<AdminPlaceResponse>('/api/v1/place', {
    ...withAuth,
    method: 'POST',
    body: payload
  })
}

export function patchAdminPlace(id: number, payload: PatchPlacePayload) {
  return apiFetch<AdminPlaceResponse>(`/api/v1/place/${id}`, {
    ...withAuth,
    method: 'PATCH',
    body: payload
  })
}

export function deleteAdminPlace(id: number) {
  return apiFetch(`/api/v1/place/${id}`, {
    ...withAuth,
    method: 'DELETE'
  })
}

export function getAdminServices() {
  return apiFetch<AdminServiceListResponse>('/api/v1/service', withAuth)
}

export function getAdminBookings() {
  return apiFetch<AdminBookingListResponse>('/api/v1/booking', withAuth)
}

export function patchAdminBooking(id: number, payload: PatchBookingPayload) {
  return apiFetch<AdminBookingResponse>(`/api/v1/booking/${id}`, {
    ...withAuth,
    method: 'PATCH',
    body: payload
  })
}

export function getAdminShifts() {
  return apiFetch<ShiftListResponse>('/api/v1/shift', withAuth)
}

export function createAdminShift(payload: CreateShiftPayload) {
  return apiFetch('/api/v1/shift', {
    ...withAuth,
    method: 'POST',
    body: payload
  })
}

export function patchAdminShift(id: number, payload: PatchShiftPayload) {
  return apiFetch(`/api/v1/shift/${id}`, {
    ...withAuth,
    method: 'PATCH',
    body: payload
  })
}

export function deleteAdminShift(id: number) {
  return apiFetch(`/api/v1/shift/${id}`, {
    ...withAuth,
    method: 'DELETE'
  })
}

export function createAdminService(payload: CreateServicePayload) {
  return apiFetch<AdminServiceResponse>('/api/v1/service', {
    ...withAuth,
    method: 'POST',
    body: payload
  })
}

export function patchAdminService(id: number, payload: PatchServicePayload) {
  return apiFetch<AdminServiceResponse>(`/api/v1/service/${id}`, {
    ...withAuth,
    method: 'PATCH',
    body: payload
  })
}

export function deleteAdminService(id: number) {
  return apiFetch(`/api/v1/service/${id}`, {
    ...withAuth,
    method: 'DELETE'
  })
}

export function getAdminSiteSettings() {
  return apiFetch<{ settings: { id: number; settings_json: number[]; gallery_items_json: any[] } }>('/api/v1/settings', withAuth)
}

export function updateAdminSiteSettings(payload: { settings_json: number[]; gallery_items_json?: any[] }) {
  return apiFetch<{ settings: { id: number; settings_json: number[]; gallery_items_json: any[] } }>('/api/v1/settings', {
    ...withAuth,
    method: 'POST',
    body: payload
  })
}

export function getAdminStats() {
  return apiFetch<AdminStatsResponse>('/api/v1/admin/stats', withAuth)
}

