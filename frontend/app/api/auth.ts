import { apiFetch } from './client'
import type { SessionResponse } from './types'

const AUTH_API_PREFIX = '/api/v1/auth'

export function fetchAuthSession() {
  return apiFetch<SessionResponse>(`${AUTH_API_PREFIX}/session`, {
    credentials: 'include'
  })
}

export function logoutSession() {
  return apiFetch(`${AUTH_API_PREFIX}/logout`, {
    method: 'POST',
    credentials: 'include'
  })
}

export function buildGoogleAuthUrl(backendUrl: string, returnTo: string) {
  const normalizedBackendUrl = backendUrl.replace(/\/$/, '')

  return `${normalizedBackendUrl}${AUTH_API_PREFIX}/google?${new URLSearchParams({
    return_to: returnTo
  }).toString()}`
}

export function updateUserProfile(data: { full_name?: string; phone?: string }) {
	return apiFetch<{ user: any }>('/api/v1/user/me', {
		method: 'PATCH',
		body: JSON.stringify(data),
		credentials: 'include'
	})
}
