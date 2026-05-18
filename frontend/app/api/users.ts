import { apiFetch } from './client'
import type { BookingRecord } from './types'

export interface UserDoc {
  id: string
  full_name: string
  email: string
  phone: string
  role: 'admin' | 'seller' | 'client'
  is_active: boolean
  created_at: string
}

export function searchAdminUsers(q?: string) {
  const params = q ? `?q=${encodeURIComponent(q)}` : ''
  return apiFetch<{ users: UserDoc[] }>(`/api/v1/admin/users${params}`, {
    credentials: 'include'
  })
}

export function listSellers() {
  return apiFetch<{ sellers: UserDoc[] }>('/api/v1/admin/sellers', {
    credentials: 'include'
  })
}

export function setUserRole(userId: string, role: 'admin' | 'seller' | 'client') {
  return apiFetch<{ user: UserDoc }>(`/api/v1/admin/users/${userId}/role`, {
    method: 'PATCH',
    credentials: 'include',
    body: { role }
  })
}

export function getBookingForSeller(bookingId: number) {
  return apiFetch<{ booking: BookingRecord }>(`/api/v1/seller/booking/${bookingId}`, {
    credentials: 'include'
  })
}

export function listBookingsForSeller() {
  return apiFetch<{ bookings: BookingRecord[] }>('/api/v1/seller/bookings', {
    credentials: 'include'
  })
}
