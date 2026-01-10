import {
  signInWithPopup,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut as firebaseSignOut,
  onAuthStateChanged,
  GoogleAuthProvider,
  GithubAuthProvider,
  type User,
  type UserCredential
} from 'firebase/auth'
import { getFirebaseAuth, isFirebaseConfigured } from '~/utils/firebase'

/**
 * Authentication state and methods composable.
 */
export function useAuth() {
  const user = useState<User | null>('auth-user', () => null)
  const loading = useState<boolean>('auth-loading', () => true)
  const error = useState<string | null>('auth-error', () => null)
  const isConfigured = useState<boolean>('auth-configured', () => false)

  /**
   * Initialize authentication state listener.
   * Call this in your app.vue or a plugin.
   */
  function initAuth(): void {
    if (!isFirebaseConfigured()) {
      loading.value = false
      isConfigured.value = false
      console.warn('Firebase is not configured. Authentication is disabled.')
      return
    }

    isConfigured.value = true
    const auth = getFirebaseAuth()

    onAuthStateChanged(auth, (firebaseUser) => {
      user.value = firebaseUser
      loading.value = false
    })
  }

  /**
   * Sign in with email and password.
   */
  async function signInWithEmail(email: string, password: string): Promise<UserCredential> {
    error.value = null
    loading.value = true

    try {
      const auth = getFirebaseAuth()
      const result = await signInWithEmailAndPassword(auth, email, password)
      return result
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to sign in'
      error.value = errorMessage
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new account with email and password.
   */
  async function signUpWithEmail(email: string, password: string): Promise<UserCredential> {
    error.value = null
    loading.value = true

    try {
      const auth = getFirebaseAuth()
      const result = await createUserWithEmailAndPassword(auth, email, password)
      return result
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to sign up'
      error.value = errorMessage
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Sign in with Google OAuth.
   */
  async function signInWithGoogle(): Promise<UserCredential> {
    error.value = null
    loading.value = true

    try {
      const auth = getFirebaseAuth()
      const provider = new GoogleAuthProvider()
      const result = await signInWithPopup(auth, provider)
      return result
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to sign in with Google'
      error.value = errorMessage
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Sign in with GitHub OAuth.
   */
  async function signInWithGitHub(): Promise<UserCredential> {
    error.value = null
    loading.value = true

    try {
      const auth = getFirebaseAuth()
      const provider = new GithubAuthProvider()
      const result = await signInWithPopup(auth, provider)
      return result
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to sign in with GitHub'
      error.value = errorMessage
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Sign out the current user.
   */
  async function signOut(): Promise<void> {
    error.value = null

    try {
      const auth = getFirebaseAuth()
      await firebaseSignOut(auth)
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to sign out'
      error.value = errorMessage
      throw e
    }
  }

  /**
   * Get the current user's ID token for API authentication.
   */
  async function getIdToken(): Promise<string | null> {
    if (!user.value) {
      return null
    }

    try {
      return await user.value.getIdToken()
    } catch (e) {
      console.error('Failed to get ID token:', e)
      return null
    }
  }

  /**
   * Check if user is authenticated.
   */
  const isAuthenticated = computed(() => !!user.value)

  return {
    // State
    user: readonly(user),
    loading: readonly(loading),
    error: readonly(error),
    isAuthenticated,
    isConfigured: readonly(isConfigured),

    // Methods
    initAuth,
    signInWithEmail,
    signUpWithEmail,
    signInWithGoogle,
    signInWithGitHub,
    signOut,
    getIdToken
  }
}
