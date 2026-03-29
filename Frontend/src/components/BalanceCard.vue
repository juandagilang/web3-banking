<script setup lang="ts">
import { computed } from 'vue'
import { useContract } from '@/composables/useContract'

const props = defineProps<{
  address: string
}>()

const {
  bankBalance,
  tokenBalance,
  formatBalance,
  refreshBalances,
} = useContract(props.address)

const bankBalanceFormatted = computed(() => formatBalance(bankBalance.value?.[0]?.result))
const tokenBalanceFormatted = computed(() => formatBalance(tokenBalance.value?.[0]?.result))

const refresh = () => {
  refreshBalances()
}
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="bg-gray-800 rounded-xl p-6 border border-gray-700">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-white">Bank Balance</h3>
        <button
          @click="refresh"
          class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
      <div class="text-3xl font-bold text-white mb-1">
        {{ bankBalanceFormatted }}
      </div>
      <p class="text-gray-400">W3B in Bank</p>
    </div>

    <div class="bg-gray-800 rounded-xl p-6 border border-gray-700">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-white">Wallet Balance</h3>
      </div>
      <div class="text-3xl font-bold text-white mb-1">
        {{ tokenBalanceFormatted }}
      </div>
      <p class="text-gray-400">W3B in Wallet</p>
    </div>
  </div>
</template>
