import { apiFetch } from './client'
import type {
  BookingResponse,
  CreateBookingPayload,
  GamingAvailabilityResponse,
  PatchBookingPayload
} from './types'

const BOOKING_API_PREFIX = '/api/v1/booking'

export function createBooking(payload: CreateBookingPayload) {
  return apiFetch<BookingResponse>(BOOKING_API_PREFIX, {
    method: 'POST',
    credentials: 'include',
    body: payload
  })
}

/** Creates a lounge/event booking — contact info is auto-filled from the session on the backend. */
export function createSessionBooking(payload: Omit<CreateBookingPayload, 'contact_name' | 'contact_email' | 'contact_phone'>) {
  return apiFetch<BookingResponse>(BOOKING_API_PREFIX, {
    method: 'POST',
    credentials: 'include',
    body: payload
  })
}

export function getMyBookings() {
  return apiFetch<{ bookings: any[] }>(`${BOOKING_API_PREFIX}/my`, {
    credentials: 'include'
  })
}

export function getMyBookingsCategorized() {
  return apiFetch<{ current: any[], archive: any[] }>('/api/v1/bookings/me', {
    credentials: 'include'
  })
}

export function getGamingAvailability(zoneId: number, date: string) {
  return apiFetch<GamingAvailabilityResponse>(`/api/v1/public/gaming/availability?${new URLSearchParams({
    zone_id: String(zoneId),
    date
  }).toString()}`)
}

export function patchBooking(id: number, payload: PatchBookingPayload) {
  return apiFetch<BookingResponse>(`${BOOKING_API_PREFIX}/${id}`, {
    method: 'PATCH',
    credentials: 'include',
    body: payload
  })
}
