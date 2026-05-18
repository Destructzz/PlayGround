import { useRequestFetch, useRuntimeConfig, useRequestHeaders } from '#app'
import { $fetch } from 'ofetch'

type RequestFetch = ReturnType<typeof useRequestFetch>
type RequestFetchOptions = NonNullable<Parameters<RequestFetch>[1]>
type ApiFetchOptions = Omit<RequestFetchOptions, 'responseType'> & {
  responseType?: 'json'
}

export function apiFetch<T>(url: string, options?: ApiFetchOptions) {
  if (import.meta.server) {
    const config = useRuntimeConfig()
    let backendUrl = config.public.backendUrl || 'http://localhost:8080'
    if (backendUrl === 'http://localhost') {
      backendUrl = 'http://localhost:8080'
    }
    const headers = useRequestHeaders(['cookie'])
    return $fetch<T>(backendUrl + url, {
      ...options,
      headers: {
        ...headers,
        ...options?.headers
      }
    })
  }

  return $fetch<T>(url, options)
}
