import { useRequestFetch } from '#app'
import { $fetch } from 'ofetch'

type RequestFetch = ReturnType<typeof useRequestFetch>
type RequestFetchOptions = NonNullable<Parameters<RequestFetch>[1]>
type ApiFetchOptions = Omit<RequestFetchOptions, 'responseType'> & {
  responseType?: 'json'
}

export function apiFetch<T>(url: string, options?: ApiFetchOptions) {
  if (import.meta.server) {
    return useRequestFetch()(url, options) as Promise<T>
  }

  return $fetch<T>(url, options)
}
