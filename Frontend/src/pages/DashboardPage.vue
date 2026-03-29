<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useAccount } from '@wagmi/vue'
import { wagmiConfig } from '@/composables/useWeb3'
import WalletConnect from '@/components/WalletConnect.vue'
import BalanceCard from '@/components/BalanceCard.vue'

const router = useRouter()
const authStore = useAuthStore()
const { address } = useAccount({ config: wagmiConfig })

onMounted(() => {
  if (!authStore.isAuthenticated) {
    router.push('/')
  }
})

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-gray-900">
    <nav class="bg-gray-800 border-b border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold text-white">Web3 Banking</h1>
          </div>
          <div class="flex items-center space-x-4">
            <WalletConnect />
            <button
              @click="handleLogout"
              class="px-3 py-2 text-sm text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="mb-8">
        <h2 class="text-2xl font-bold text-white">Dashboard</h2>
        <p class="text-gray-400 mt-1">Manage your W3B tokens</p>
      </div>

      <div v-if="address" class="space-y-6">
        <BalanceCard :address="address" />
      </div>

      <div v-else class="text-center py-12">
        <p class="text-gray-400">Please connect your wallet to view your balance</p>
      </div>
    </main>
  </div>
</template>
