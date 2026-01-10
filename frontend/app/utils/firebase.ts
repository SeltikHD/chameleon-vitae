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

let app: FirebaseApp | null = null
let auth: Auth | null = null

/**
 * Get Firebase configuration from runtime config.
 */
function getFirebaseConfig(): FirebaseConfig {
  const config = useRuntimeConfig()

  return {
    apiKey: config.public.firebaseApiKey as string,
    authDomain: config.public.firebaseAuthDomain as string,
    projectId: config.public.firebaseProjectId as string,
    storageBucket: config.public.firebaseStorageBucket as string,
    messagingSenderId: config.public.firebaseMessagingSenderId as string,
    appId: config.public.firebaseAppId as string
  }
}

/**
 * Initialize or get the Firebase app instance.
 */
export function getFirebaseApp(): FirebaseApp {
  if (app) {
    return app
  }

  const existingApps = getApps()
  if (existingApps.length > 0 && existingApps[0]) {
    app = existingApps[0]
    return app
  }

  const firebaseConfig = getFirebaseConfig()

  // Validate configuration.
  if (!firebaseConfig.apiKey || !firebaseConfig.projectId) {
    throw new Error('Firebase configuration is missing. Check your environment variables.')
  }

  app = initializeApp(firebaseConfig)
  return app
}

/**
 * Get the Firebase Auth instance.
 */
export function getFirebaseAuth(): Auth {
  if (auth) {
    return auth
  }

  const firebaseApp = getFirebaseApp()
  auth = getAuth(firebaseApp)
  return auth
}

/**
 * Check if Firebase is properly configured.
 */
export function isFirebaseConfigured(): boolean {
  try {
    const config = getFirebaseConfig()
    return !!(config.apiKey && config.projectId)
  } catch {
    return false
  }
}
