import { useAuthStore } from '~/stores/auth'

/**
 * Plugin to initialize authentication on app startup.
 * This runs after the Firebase plugin and sets up the auth state listener.
 */
export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()

  // Initialize auth listener (client-side only).
  authStore.initAuth()
})
