import { defineStore } from 'pinia'
import {
  signInWithPopup,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut as firebaseSignOut,
  onAuthStateChanged,
  GoogleAuthProvider,
  GithubAuthProvider,
  type User as FirebaseUser
} from 'firebase/auth'
import type { UserResponse, SyncUserResponse } from '~/types/api'

/**
 * Authentication state managed by Pinia.
 */
interface AuthState {
  /** Firebase user object. */
  firebaseUser: FirebaseUser | null
  /** Backend user profile after sync. */
  user: UserResponse | null
  /** Firebase ID Token for API authentication. */
  idToken: string | null
  /** Loading state during auth operations. */
  loading: boolean
  /** Whether auth listener has been initialized. */
  initialized: boolean
  /** Last auth error message. */
  error: string | null
}

/**
 * Pinia store for authentication state management.
 *
 * CRITICAL: This store handles the synchronization between
 * Firebase Authentication and the Chameleon Vitae backend.
 */
export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    firebaseUser: null,
    user: null,
    idToken: null,
    loading: true,
    initialized: false,
    error: null
  }),

  getters: {
    /** Whether the user is authenticated with both Firebase and Backend. */
    isAuthenticated: (state) => !!state.user && !!state.idToken,

    /** Whether Firebase is configured. */
    isFirebaseConfigured(): boolean {
      const { $isFirebaseConfigured } = useNuxtApp()
      return $isFirebaseConfigured as boolean
    },

    /** User's display name from backend profile. */
    displayName: (state) => state.user?.name || state.user?.email || 'Anonymous User',

    /** User's avatar URL. */
    avatarUrl: (state) => state.user?.picture_url || null,

    /** Backend user ID (UUID). */
    userId: (state) => state.user?.id || null
  },

  actions: {
    /**
     * Initialize the auth state listener.
     * This should be called once on app mount (via a plugin or app.vue).
     */
    initAuth() {
      if (this.initialized) return

      const { $firebaseAuth, $isFirebaseConfigured } = useNuxtApp()

      if (!$isFirebaseConfigured || !$firebaseAuth) {
        console.warn('[AuthStore] Firebase not configured. Auth disabled.')
        this.loading = false
        this.initialized = true
        return
      }

      // Listen for auth state changes.
      onAuthStateChanged($firebaseAuth, async (firebaseUser) => {
        this.firebaseUser = firebaseUser

        if (firebaseUser) {
          try {
            // Get ID Token from Firebase.
            this.idToken = await firebaseUser.getIdToken()

            // Sync with backend.
            await this.syncWithBackend(firebaseUser)
          } catch (error) {
            console.error('[AuthStore] Failed to sync with backend:', error)
            this.error = error instanceof Error ? error.message : 'Failed to sync with backend'
            this.user = null
            this.idToken = null
          }
        } else {
          this.user = null
          this.idToken = null
        }

        this.loading = false
        this.initialized = true
      })
    },

    /**
     * CRITICAL: Sync Firebase user with Backend.
     * This must be called after every successful Firebase login.
     */
    async syncWithBackend(firebaseUser: FirebaseUser) {
      const config = useRuntimeConfig()
      const apiBase = config.public.apiBase as string

      const syncPayload = {
        firebase_uid: firebaseUser.uid,
        email: firebaseUser.email,
        name: firebaseUser.displayName,
        picture: firebaseUser.photoURL
      }

      try {
        const response = await $fetch<SyncUserResponse>(`${apiBase}/auth/sync`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${this.idToken}`
          },
          body: syncPayload
        })

        // Fetch full user profile after sync.
        await this.fetchUserProfile()

        return response
      } catch (error) {
        console.error('[AuthStore] Backend sync failed:', error)
        throw error
      }
    },

    /**
     * Fetch the current user's profile from the backend.
     */
    async fetchUserProfile() {
      const config = useRuntimeConfig()
      const apiBase = config.public.apiBase as string

      if (!this.idToken) {
        throw new Error('No ID token available')
      }

      const user = await $fetch<UserResponse>(`${apiBase}/me`, {
        headers: {
          Authorization: `Bearer ${this.idToken}`
        }
      })

      this.user = user
    },

    /**
     * Sign in with email and password.
     */
    async signInWithEmail(email: string, password: string) {
      const { $firebaseAuth } = useNuxtApp()

      if (!$firebaseAuth) {
        throw new Error('Firebase not configured')
      }

      this.loading = true
      this.error = null

      try {
        await signInWithEmailAndPassword($firebaseAuth, email, password)
        // Auth state listener will handle the rest.
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to sign in'
        this.loading = false
        throw error
      }
    },

    /**
     * Sign up with email and password.
     */
    async signUpWithEmail(email: string, password: string) {
      const { $firebaseAuth } = useNuxtApp()

      if (!$firebaseAuth) {
        throw new Error('Firebase not configured')
      }

      this.loading = true
      this.error = null

      try {
        await createUserWithEmailAndPassword($firebaseAuth, email, password)
        // Auth state listener will handle the rest.
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to sign up'
        this.loading = false
        throw error
      }
    },

    /**
     * Sign in with Google OAuth.
     */
    async signInWithGoogle() {
      const { $firebaseAuth } = useNuxtApp()

      if (!$firebaseAuth) {
        throw new Error('Firebase not configured')
      }

      this.loading = true
      this.error = null

      try {
        const provider = new GoogleAuthProvider()
        await signInWithPopup($firebaseAuth, provider)
        // Auth state listener will handle the rest.
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to sign in with Google'
        this.loading = false
        throw error
      }
    },

    /**
     * Sign in with GitHub OAuth.
     */
    async signInWithGitHub() {
      const { $firebaseAuth } = useNuxtApp()

      if (!$firebaseAuth) {
        throw new Error('Firebase not configured')
      }

      this.loading = true
      this.error = null

      try {
        const provider = new GithubAuthProvider()
        await signInWithPopup($firebaseAuth, provider)
        // Auth state listener will handle the rest.
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to sign in with GitHub'
        this.loading = false
        throw error
      }
    },

    /**
     * Sign out from both Firebase and clear backend session.
     */
    async signOut() {
      const { $firebaseAuth } = useNuxtApp()

      if (!$firebaseAuth) {
        throw new Error('Firebase not configured')
      }

      try {
        await firebaseSignOut($firebaseAuth)
        this.user = null
        this.idToken = null
        this.firebaseUser = null
        this.error = null
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to sign out'
        throw error
      }
    },

    /**
     * Refresh the ID token.
     * Call this before making API requests to ensure token freshness.
     */
    async refreshIdToken(): Promise<string | null> {
      if (!this.firebaseUser) return null

      try {
        this.idToken = await this.firebaseUser.getIdToken(true)
        return this.idToken
      } catch (error) {
        console.error('[AuthStore] Failed to refresh ID token:', error)
        return null
      }
    },

    /**
     * Clear any auth errors.
     */
    clearError() {
      this.error = null
    }
  }
})
