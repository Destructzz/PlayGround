export function useMockSubmission() {
  const pending = ref(false)

  async function submit(action: () => void, shouldFail = false) {
    if (pending.value) {
      return
    }

    pending.value = true

    try {
      await new Promise(resolve => setTimeout(resolve, 700))

      if (shouldFail) {
        throw new Error('mock-submission-error')
      }

      action()
    } finally {
      pending.value = false
    }
  }

  return {
    pending,
    submit
  }
}
