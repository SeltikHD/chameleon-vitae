import { useAuthStore } from '~/stores/auth'
import type { ErrorResponse } from '~/types/api'

/**
 * HTTP methods supported by the API.
 */
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

/**
 * Options for useApiFetch composable.
 */
export type UseApiFetchOptions<_T> = Omit<Parameters<typeof useFetch>[1], 'baseURL'> & {
  /** Whether to skip authentication (for public endpoints). */
  skipAuth?: boolean
  /** Custom error handler. */
  onError?: (error: Error) => void
}

/**
 * Options for imperative apiFetch function.
 */
export interface ApiFetchOptions {
  /** HTTP method. Defaults to GET. */
  method?: HttpMethod
  /** Request body for POST/PUT/PATCH. */
  body?: unknown
  /** Additional headers. */
  headers?: Record<string, string>
  /** Query parameters. */
  query?: Record<string, unknown>
  /** Whether to skip authentication (for public endpoints). */
  skipAuth?: boolean
  /** Custom error handler. */
  onError?: (error: Error) => void
}

/**
 * API error with structured error information.
 */
export class ApiError extends Error {
  code: string
  status: number
  details?: Array<{ field: string; message: string }>

  constructor(
    message: string,
    code: string,
    status: number,
    details?: Array<{ field: string; message: string }>
  ) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
    this.details = details
  }
}

/**
 * Composable for making authenticated API requests.
 *
 * Features:
 * - Automatically attaches Authorization header with Firebase ID Token
 * - Handles token refresh on 401 errors
 * - Provides typed error handling
 * - Redirects to login on authentication failure
 *
 * @example
 * ```ts
 * // GET request
 * const { data, error, pending } = await useApiFetch<UserResponse>('/me')
 *
 * // POST request
 * const { data } = await useApiFetch<ExperienceResponse>('/experiences', {
 *   method: 'POST',
 *   body: { type: 'work', title: 'Engineer', ... }
 * })
 * ```
 */
export function useApiFetch<T>(
  endpoint: string | (() => string),
  options: UseApiFetchOptions<T> = {}
) {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()
  const toast = useToast()

  const { skipAuth = false, onError, ...fetchOptions } = options

  return useFetch<T>(endpoint, {
    baseURL: config.public.apiBase as string,

    // Attach Authorization header.
    async onRequest({ options }) {
      if (skipAuth) return

      // Ensure we have a fresh token.
      let token = authStore.idToken

      if (!token && authStore.firebaseUser) {
        token = await authStore.refreshIdToken()
      }

      if (token) {
        const existingHeaders = options.headers ?? {}
        options.headers = new Headers(existingHeaders as HeadersInit)
        options.headers.set('Authorization', `Bearer ${token}`)
      }
    },

    // Handle response errors.
    async onResponseError({ response }) {
      const status = response.status
      const body = response._data as ErrorResponse | undefined

      // Extract error details.
      const errorCode = body?.error?.code || 'UNKNOWN_ERROR'
      const errorMessage = body?.error?.message || 'An unexpected error occurred'
      const errorDetails = body?.error?.details

      // Handle 401 Unauthorized.
      if (status === 401) {
        // Try to refresh the token.
        const newToken = await authStore.refreshIdToken()

        if (!newToken) {
          // Token refresh failed, redirect to login.
          toast.add({
            title: 'Session Expired',
            description: 'Please sign in again.',
            color: 'error'
          })
          await authStore.signOut()
          await navigateTo('/login')
          return
        }

        // Token refreshed, the caller should retry.
        throw new ApiError('Token refreshed, please retry', 'TOKEN_REFRESHED', 401)
      }

      // Handle 403 Forbidden.
      if (status === 403) {
        toast.add({
          title: 'Access Denied',
          description: errorMessage,
          color: 'error'
        })
      }

      // Handle validation errors.
      if (status === 400 || status === 422) {
        if (errorDetails && errorDetails.length > 0) {
          const fieldErrors = errorDetails.map((d) => `${d.field}: ${d.message}`).join(', ')
          toast.add({
            title: 'Validation Error',
            description: fieldErrors,
            color: 'warning'
          })
        } else {
          toast.add({
            title: 'Invalid Request',
            description: errorMessage,
            color: 'warning'
          })
        }
      }

      // Handle 404 Not Found.
      if (status === 404) {
        toast.add({
          title: 'Not Found',
          description: errorMessage,
          color: 'error'
        })
      }

      // Handle 500+ Server Errors.
      if (status >= 500) {
        toast.add({
          title: 'Server Error',
          description: 'Something went wrong. Please try again later.',
          color: 'error'
        })
      }

      // Create and throw structured error.
      const apiError = new ApiError(errorMessage, errorCode, status, errorDetails)

      if (onError) {
        onError(apiError)
      }

      throw apiError
    },

    ...fetchOptions
  })
}

/**
 * Imperative API fetch for non-reactive use cases.
 *
 * Use this for mutations (POST, PUT, PATCH, DELETE) where you don't need
 * reactive data binding.
 *
 * @example
 * ```ts
 * const user = await apiFetch<UserResponse>('/me', {
 *   method: 'PATCH',
 *   body: { name: 'New Name' }
 * })
 * ```
 */
export async function apiFetch<T>(endpoint: string, options: ApiFetchOptions = {}): Promise<T> {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()
  const toast = useToast()

  const { skipAuth = false, onError, method, body, query, headers: customHeaders } = options

  // Build headers.
  const headers: Record<string, string> = customHeaders
    ? { 'Content-Type': 'application/json', ...customHeaders }
    : { 'Content-Type': 'application/json' }

  if (!skipAuth) {
    let token = authStore.idToken

    if (!token && authStore.firebaseUser) {
      token = await authStore.refreshIdToken()
    }

    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
  }

  try {
    const data = await $fetch<T>(endpoint, {
      baseURL: config.public.apiBase as string,
      method,
      body: body as BodyInit | Record<string, unknown> | null | undefined,
      query,
      headers
    })

    return data
  } catch (error: unknown) {
    // Handle $fetch errors.
    if (error && typeof error === 'object' && 'response' in error) {
      const fetchError = error as { response: Response; data: ErrorResponse }
      const status = fetchError.response?.status || 500
      const body = fetchError.data

      const errorCode = body?.error?.code || 'UNKNOWN_ERROR'
      const errorMessage = body?.error?.message || 'An unexpected error occurred'
      const errorDetails = body?.error?.details

      // Handle 401.
      if (status === 401) {
        const newToken = await authStore.refreshIdToken()

        if (!newToken) {
          toast.add({
            title: 'Session Expired',
            description: 'Please sign in again.',
            color: 'error'
          })
          await authStore.signOut()
          await navigateTo('/login')
        }
      }

      const apiError = new ApiError(errorMessage, errorCode, status, errorDetails)

      if (onError) {
        onError(apiError)
      }

      throw apiError
    }

    // Re-throw unknown errors.
    throw error
  }
}
