import type { RouteLocationNormalizedLoadedGeneric } from 'vue-router'
import type { MockUiState } from '../utils/experienceData'

const validStates: MockUiState[] = ['default', 'loading', 'empty', 'error', 'success']

export function useExperienceMockState(route: RouteLocationNormalizedLoadedGeneric) {
  const uiState = computed<MockUiState>(() => {
    const state = route.query.state
    const normalized = Array.isArray(state) ? state[0] : state

    if (normalized && validStates.includes(normalized as MockUiState)) {
      return normalized as MockUiState
    }

    return 'default'
  })

  const isLoading = computed(() => uiState.value === 'loading')
  const isEmpty = computed(() => uiState.value === 'empty')
  const isError = computed(() => uiState.value === 'error')
  const isSuccessPreview = computed(() => uiState.value === 'success')

  return {
    uiState,
    isLoading,
    isEmpty,
    isError,
    isSuccessPreview
  }
}
