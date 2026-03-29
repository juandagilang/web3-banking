<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import type { Transaction, Pagination } from '@/types'
import { getTransactions } from '@/api/account'

const props = defineProps<{
  address: string
}>()

const transactions = ref<Transaction[]>([])
const pagination = ref<Pagination>({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
})
const isLoading = ref(false)
const error = ref<string | null>(null)

const fetchTransactions = async () => {
  isLoading.value = true
  error.value = null
  
  const result = await getTransactions(props.address, pagination.value.page, pagination.value.limit)
  
  if (result) {
    transactions.value = result.data
    pagination.value = result.pagination
  } else {
    error.value = 'Failed to load transactions'
  }
  
  isLoading.value = false
}

const goToPage = (page: number) => {
  if (page >= 1 && page <= pagination.value.total_pages) {
    pagination.value.page = page
    fetchTransactions()
  }
}

const formatAmount = (amount: string) => {
  const value = BigInt(amount)
  const formatted = Number(value) / 1e18
  return formatted.toFixed(4)
}

const formatTimestamp = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString()
}

const formatAddress = (address: string) => {
  if (!address) return ''
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}

const getTypeColor = (type: string) => {
  switch (type) {
    case 'deposit':
      return 'text-green-400 bg-green-400/10'
    case 'withdrawal':
      return 'text-red-400 bg-red-400/10'
    case 'transfer':
      return 'text-blue-400 bg-blue-400/10'
    default:
      return 'text-gray-400 bg-gray-400/10'
  }
}

watch(() => props.address, () => {
  if (props.address) {
    fetchTransactions()
  }
})

onMounted(() => {
  if (props.address) {
    fetchTransactions()
  }
})
</script>

<template>
  <div class="bg-gray-800 rounded-xl border border-gray-700">
    <div class="p-6 border-b border-gray-700">
      <h3 class="text-lg font-semibold text-white">Transaction History</h3>
    </div>

    <div v-if="isLoading" class="p-8 text-center">
      <svg class="animate-spin h-8 w-8 text-primary-500 mx-auto" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
      <p class="mt-2 text-gray-400">Loading transactions...</p>
    </div>

    <div v-else-if="error" class="p-8 text-center">
      <p class="text-red-400">{{ error }}</p>
      <button
        @click="fetchTransactions"
        class="mt-4 px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg"
      >
        Retry
      </button>
    </div>

    <div v-else-if="transactions.length === 0" class="p-8 text-center">
      <p class="text-gray-400">No transactions yet</p>
    </div>

    <div v-else class="divide-y divide-gray-700">
      <div
        v-for="tx in transactions"
        :key="tx.id"
        class="p-4 hover:bg-gray-750 transition-colors"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <span
              :class="[
                'px-2 py-1 text-xs font-medium rounded uppercase',
                getTypeColor(tx.type)
              ]"
            >
              {{ tx.type }}
            </span>
            <div>
              <p class="text-white font-medium">
                {{ formatAmount(tx.amount) }} W3B
              </p>
              <p class="text-sm text-gray-400">
                From: {{ formatAddress(tx.from) }}
                <span v-if="tx.to"> → To: {{ formatAddress(tx.to) }}</span>
              </p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm text-gray-400">{{ formatTimestamp(tx.timestamp) }}</p>
            <p class="text-xs text-gray-500 mt-1">Block #{{ tx.block_number }}</p>
          </div>
        </div>
      </div>

      <div v-if="pagination.total_pages > 1" class="p-4 border-t border-gray-700">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-400">
            Page {{ pagination.page }} of {{ pagination.total_pages }}
          </p>
          <div class="flex space-x-2">
            <button
              @click="goToPage(pagination.page - 1)"
              :disabled="pagination.page <= 1"
              class="px-3 py-1 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm rounded"
            >
              Previous
            </button>
            <button
              @click="goToPage(pagination.page + 1)"
              :disabled="pagination.page >= pagination.total_pages"
              class="px-3 py-1 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm rounded"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
