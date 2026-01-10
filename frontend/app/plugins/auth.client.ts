import { useAuthStore } from '~/stores/auth'

/**
 * Plugin to initialize authentication on app startup.
 * This runs AFTER the Firebase plugin and sets up the auth state listener.
 *
 * CRITICAL: This plugin depends on 'firebase' plugin via the dependsOn option.
 */
export default defineNuxtPlugin({
  name: 'auth',
  dependsOn: ['firebase'], // Ensure Firebase is initialized first.
  setup() {
    const authStore = useAuthStore()

    // Initialize auth listener (client-side only).
    authStore.initAuth()
  }
})
