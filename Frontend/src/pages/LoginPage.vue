<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useWalletStore } from '@/store/wallet'
import { useAuthStore } from '@/store/auth'
import { useWeb3Modal } from '@web3modal/wagmi/vue'
import { useAccount, useSignMessage } from '@wagmi/vue'
import { wagmiConfig } from '@/composables/useWeb3'

const router = useRouter()
const walletStore = useWalletStore()
const authStore = useAuthStore()
const web3Modal = useWeb3Modal()

const { address, isConnected } = useAccount({ config: wagmiConfig })
const { signMessageAsync } = useSignMessage({ config: wagmiConfig })

const isLoading = ref(false)
const error = ref<string | null>(null)

watch(isConnected, async (connected) => {
  if (connected && address.value) {
    await handleLogin()
  }
})

const handleConnect = async () => {
  try {
    error.value = null
    await web3Modal.open()
  } catch (err: unknown) {
    error.value = 'Failed to open wallet modal'
  }
}

const handleLogin = async () => {
  if (!address.value) return

  try {
    isLoading.value = true
    error.value = null

    const nonceData = await authStore.getNonce(address.value)
    if (!nonceData) {
      error.value = authStore.error || 'Failed to get nonce'
      return
    }

    const signature = await signMessageAsync({
      message: nonceData.message,
    })

    const success = await authStore.login(address.value, signature)

    if (success) {
      router.push('/dashboard')
    } else {
      error.value = authStore.error || 'Login failed'
    }
  } catch (err: unknown) {
    const wagmiError = err as { message?: string }
    if (wagmiError.message?.includes('User rejected')) {
      error.value = 'Signature rejected by user'
    } else {
      error.value = 'Failed to sign message'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-900 px-4">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <h1 class="text-4xl font-bold text-white mb-2">Web3 Banking</h1>
        <p class="text-gray-400">Connect your wallet to access your banking dashboard</p>
      </div>

      <div class="bg-gray-800 rounded-xl p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-6">
          <div class="w-20 h-20 bg-primary-600 rounded-full flex items-center justify-center">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
          </div>

          <button
            @click="handleConnect"
            :disabled="isLoading"
            class="w-full py-3 px-4 bg-primary-600 hover:bg-primary-700 disabled:bg-primary-800 text-white font-semibold rounded-lg transition-colors flex items-center justify-center space-x-2"
          >
            <svg v-if="isLoading" class="animate-spin h-5 w-5" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            <span>{{ isLoading ? 'Connecting...' : 'Connect Wallet' }}</span>
          </button>

          <div v-if="error" class="w-full p-4 bg-red-500/10 border border-red-500/50 rounded-lg">
            <p class="text-red-400 text-sm text-center">{{ error }}</p>
          </div>

          <p class="text-gray-500 text-sm text-center">
            By connecting, you agree to sign a message for authentication
          </p>
        </div>
      </div>

      <div class="text-center text-gray-500 text-sm">
        <p>Your wallet will be used for secure, decentralized authentication</p>
      </div>
    </div>
  </div>
</template>
