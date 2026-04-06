export function useApi() {
  const config = useRuntimeConfig()
  const base = config.public.apiBase as string

  async function get<T = any>(path: string): Promise<T> {
    const res = await $fetch<T>(`${base}${path}`)
    return res
  }

  async function post<T = any>(path: string, body: any): Promise<T> {
    return await $fetch<T>(`${base}${path}`, { method: 'POST', body })
  }

  async function patch<T = any>(path: string, body: any): Promise<T> {
    return await $fetch<T>(`${base}${path}`, { method: 'PATCH', body })
  }

  async function del<T = any>(path: string): Promise<T> {
    return await $fetch<T>(`${base}${path}`, { method: 'DELETE' })
  }

  return { get, post, patch, del }
}
