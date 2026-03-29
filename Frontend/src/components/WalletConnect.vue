<script setup lang="ts">
import { useWalletStore } from '@/store/wallet'
import { useAuthStore } from '@/store/auth'
import { useWeb3Modal } from '@web3modal/wagmi/vue'
import { useAccount, useSignMessage } from '@wagmi/vue'
import { wagmiConfig } from '@/composables/useWeb3'
import { useRouter } from 'vue-router'

const walletStore = useWalletStore()
const authStore = useAuthStore()
const router = useRouter()
const web3Modal = useWeb3Modal()
const { address, isConnected } = useAccount({ config: wagmiConfig })
const { signMessageAsync } = useSignMessage({ config: wagmiConfig })

const handleClick = async () => {
  if (isConnected.value) {
    authStore.logout()
    router.push('/')
  } else {
    await web3Modal.open()
  }
}

const handleConnectAndLogin = async () => {
  if (!address.value) return

  try {
    const nonceData = await authStore.getNonce(address.value)
    if (!nonceData) return

    const signature = await signMessageAsync({
      message: nonceData.message,
    })

    const success = await authStore.login(address.value, signature)
    if (!success) {
      authStore.logout()
    }
  } catch {
    // User rejected or error
  }
}
</script>

<template>
  <button
    @click="handleClick"
    class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium rounded-lg transition-colors flex items-center space-x-2"
  >
    <svg v-if="isConnected" class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
      <circle cx="12" cy="12" r="8" />
    </svg>
    <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
    </svg>
    <span>{{ isConnected ? walletStore.formatAddress(address) : 'Connect Wallet' }}</span>
  </button>
</template>
