import { initializeApp, getApps, type FirebaseApp } from 'firebase/app'
import { getAuth, type Auth } from 'firebase/auth'

/**
 * Firebase configuration from environment variables.
 * These values are safe to expose in client-side code.
 */
interface FirebaseConfig {
  apiKey: string
  authDomain: string
  projectId: string
  storageBucket: string
  messagingSenderId: string
  appId: string
}

let firebaseApp: FirebaseApp | null = null
let firebaseAuth: Auth | null = null

/**
 * Nuxt plugin to initialize Firebase on client-side only.
 * Firebase is configured using runtime config from environment variables.
 */
export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()

  const firebaseConfig: FirebaseConfig = {
    apiKey: config.public.firebaseApiKey as string,
    authDomain: config.public.firebaseAuthDomain as string,
    projectId: config.public.firebaseProjectId as string,
    storageBucket: config.public.firebaseStorageBucket as string,
    messagingSenderId: config.public.firebaseMessagingSenderId as string,
    appId: config.public.firebaseAppId as string
  }

  // Check if Firebase is configured.
  const isConfigured = !!(firebaseConfig.apiKey && firebaseConfig.projectId)

  if (isConfigured) {
    // Initialize Firebase app (singleton pattern).
    const existingApps = getApps()
    if (existingApps.length > 0 && existingApps[0]) {
      firebaseApp = existingApps[0]
    } else {
      firebaseApp = initializeApp(firebaseConfig)
    }

    // Initialize Firebase Auth.
    firebaseAuth = getAuth(firebaseApp)
  } else {
    console.warn(
      '[Firebase] Not configured. Set NUXT_PUBLIC_FIREBASE_* environment variables to enable authentication.'
    )
  }

  return {
    provide: {
      firebase: firebaseApp,
      firebaseAuth,
      isFirebaseConfigured: isConfigured
    }
  }
})
