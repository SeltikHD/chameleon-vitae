import { useAuthStore } from '~/stores/auth'

/**
 * Middleware to protect routes that require authentication.
 * Redirects to /login if user is not authenticated.
 *
 * CRITICAL: This middleware waits for Firebase to fully initialize
 * before checking authentication state, preventing the "F5 problem"
 * where refreshing the page causes logout.
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // Only run on client-side.
  if (import.meta.server) return

  const authStore = useAuthStore()

  // If auth is not initialized yet, wait for it.
  // This prevents the "F5 problem" where the page loads before Firebase initializes.
  if (!authStore.initialized) {
    await waitForAuthInitialization(authStore)
  }

  // If still loading after initialization, wait a bit more.
  // This handles the case where onAuthStateChanged is still processing.
  if (authStore.loading) {
    await waitForAuthLoading(authStore)
  }

  // Check authentication.
  if (!authStore.isAuthenticated) {
    // Store the intended destination to redirect back after login.
    const redirectTo = to.fullPath === '/login' ? undefined : to.fullPath
    return navigateTo({
      path: '/login',
      query: redirectTo ? { redirect: redirectTo } : undefined
    })
  }
})

/**
 * Wait for auth store to be initialized.
 * Uses a polling approach with exponential backoff.
 */
async function waitForAuthInitialization(
  authStore: ReturnType<typeof useAuthStore>
): Promise<void> {
  const maxWait = 8000 // 8 seconds max
  const startTime = Date.now()

  return new Promise((resolve) => {
    const check = () => {
      if (authStore.initialized) {
        resolve()
        return
      }

      // Check if we've exceeded the timeout
      if (Date.now() - startTime > maxWait) {
        console.warn('[AuthMiddleware] Auth initialization timeout')
        resolve()
        return
      }

      // Check again after a short delay
      setTimeout(check, 50)
    }

    check()
  })
}

/**
 * Wait for auth loading to complete.
 */
async function waitForAuthLoading(authStore: ReturnType<typeof useAuthStore>): Promise<void> {
  const maxWait = 3000 // 3 seconds max
  const startTime = Date.now()

  return new Promise((resolve) => {
    const check = () => {
      if (!authStore.loading) {
        resolve()
        return
      }

      if (Date.now() - startTime > maxWait) {
        resolve()
        return
      }

      setTimeout(check, 50)
    }

    check()
  })
}
