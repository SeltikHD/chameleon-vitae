import { useAuthStore } from '~/stores/auth'

/**
 * Middleware to protect routes that require authentication.
 * Redirects to /login if user is not authenticated.
 */
export default defineNuxtRouteMiddleware(async () => {
  // Only run on client-side.
  if (import.meta.server) return

  const authStore = useAuthStore()

  // Wait for auth to initialize.
  if (!authStore.initialized) {
    // Wait for initialization with a timeout.
    await new Promise<void>((resolve) => {
      const checkInit = setInterval(() => {
        if (authStore.initialized) {
          clearInterval(checkInit)
          resolve()
        }
      }, 50)

      // Timeout after 5 seconds.
      setTimeout(() => {
        clearInterval(checkInit)
        resolve()
      }, 5000)
    })
  }

  // Check authentication.
  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }
})
