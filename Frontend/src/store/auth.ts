import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginResponse, NonceResponse } from '@/types'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const walletAddress = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setWalletAddress = (address: string) => {
    walletAddress.value = address
  }

  const clearAuth = () => {
    token.value = null
    walletAddress.value = null
    localStorage.removeItem('token')
  }

  const getNonce = async (address: string): Promise<NonceResponse | null> => {
    try {
      error.value = null
      const response = await axios.post<NonceResponse>(
        `${API_URL}/api/v1/auth/nonce`,
        { address }
      )
      return response.data.data || null
    } catch (err: unknown) {
      const axiosError = err as { response?: { data?: { error?: string } } }
      error.value = axiosError.response?.data?.error || 'Failed to get nonce'
      return null
    }
  }

  const login = async (address: string, signature: string): Promise<boolean> => {
    try {
      isLoading.value = true
      error.value = null
      const response = await axios.post<LoginResponse>(
        `${API_URL}/api/v1/auth/login`,
        { address, signature }
      )
      if (response.data.data) {
        setToken(response.data.data.token)
        setWalletAddress(address)
        return true
      }
      return false
    } catch (err: unknown) {
      const axiosError = err as { response?: { data?: { error?: string } } }
      error.value = axiosError.response?.data?.error || 'Login failed'
      return false
    } finally {
      isLoading.value = false
    }
  }

  const logout = () => {
    clearAuth()
  }

  return {
    token,
    walletAddress,
    isLoading,
    error,
    isAuthenticated,
    setToken,
    setWalletAddress,
    getNonce,
    login,
    logout,
    clearAuth,
  }
})
