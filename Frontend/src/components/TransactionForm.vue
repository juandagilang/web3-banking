<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useContract } from '@/composables/useContract'

const props = defineProps<{
  address: string
}>()

const {
  bankBalance,
  tokenBalance,
  txState,
  allowance,
  isApproveConfirming,
  isApproveConfirmed,
  isDepositConfirming,
  isDepositConfirmed,
  isWithdrawConfirming,
  isWithdrawConfirmed,
  isTransferConfirming,
  isTransferConfirmed,
  resetTxState,
  doApprove,
  doDeposit,
  doWithdraw,
  doTransfer,
  refreshBalances,
  formatBalance,
} = useContract(props.address)

const activeTab = ref<'deposit' | 'withdraw' | 'transfer'>('deposit')
const amount = ref('')
const recipient = ref('')
const userMessage = ref('')
const showUserMessage = ref(false)

const tokenBalanceFormatted = computed(() => formatBalance(tokenBalance.value?.[0]?.result))
const bankBalanceFormatted = computed(() => formatBalance(bankBalance.value?.[0]?.result))
const allowanceFormatted = computed(() => formatBalance(allowance.value?.[0]?.result))

const needsApproval = computed(() => {
  if (!amount.value) return false
  const amountWei = BigInt(amount.value) * BigInt(1e18)
  const currentAllowance = allowance.value?.[0]?.result || BigInt(0)
  return amountWei > currentAllowance
})

const isAnyPending = computed(() => {
  return (
    txState.value.isPending ||
    isApproveConfirming.value ||
    isDepositConfirming.value ||
    isWithdrawConfirming.value ||
    isTransferConfirming.value
  )
})

const isAnyConfirmed = computed(() => {
  return (
    isApproveConfirmed.value ||
    isDepositConfirmed.value ||
    isWithdrawConfirmed.value ||
    isTransferConfirmed.value
  )
})

const currentConfirmations = computed(() => {
  if (isApproveConfirming.value || isApproveConfirmed.value) return 1
  if (isDepositConfirming.value || isDepositConfirmed.value) return 1
  if (isWithdrawConfirming.value || isWithdrawConfirmed.value) return 1
  if (isTransferConfirming.value || isTransferConfirmed.value) return 1
  return 0
})

watch(isAnyConfirmed, (confirmed) => {
  if (confirmed) {
    showUserMessage.value = true
    userMessage.value = 'Transaction confirmed!'
    refreshBalances()
    setTimeout(() => {
      resetTxState()
      amount.value = ''
      recipient.value = ''
    }, 2000)
  }
})

watch(() => txState.value.error, (error) => {
  if (error) {
    showUserMessage.value = true
    userMessage.value = error
  }
})

const handleDeposit = async () => {
  if (!amount.value) return
  userMessage.value = ''
  showUserMessage.value = false

  if (needsApproval.value) {
    await doApprove(amount.value)
    if (!txState.value.error) {
      await doDeposit(amount.value)
    }
  } else {
    await doDeposit(amount.value)
  }
}

const handleWithdraw = async () => {
  if (!amount.value) return
  userMessage.value = ''
  showUserMessage.value = false
  await doWithdraw(amount.value)
}

const handleTransfer = async () => {
  if (!amount.value || !recipient.value) return
  userMessage.value = ''
  showUserMessage.value = false
  await doTransfer(recipient.value, amount.value)
}

const dismissMessage = () => {
  showUserMessage.value = false
  if (!isAnyPending.value) {
    resetTxState()
  }
}
</script>

<template>
  <div class="bg-gray-800 rounded-xl border border-gray-700">
    <div class="border-b border-gray-700">
      <nav class="flex -mb-px">
        <button
          @click="activeTab = 'deposit'"
          :class="[
            'px-6 py-3 text-sm font-medium border-b-2 transition-colors',
            activeTab === 'deposit'
              ? 'border-primary-500 text-primary-500'
              : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-600'
          ]"
        >
          Deposit
        </button>
        <button
          @click="activeTab = 'withdraw'"
          :class="[
            'px-6 py-3 text-sm font-medium border-b-2 transition-colors',
            activeTab === 'withdraw'
              ? 'border-primary-500 text-primary-500'
              : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-600'
          ]"
        >
          Withdraw
        </button>
        <button
          @click="activeTab = 'transfer'"
          :class="[
            'px-6 py-3 text-sm font-medium border-b-2 transition-colors',
            activeTab === 'transfer'
              ? 'border-primary-500 text-primary-500'
              : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-600'
          ]"
        >
          Transfer
        </button>
      </nav>
    </div>

    <div class="p-6">
      <div v-if="activeTab === 'deposit'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Amount (W3B)</label>
          <input
            v-model="amount"
            type="number"
            step="0.0001"
            placeholder="0.0"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-primary-500"
          />
          <p class="mt-1 text-xs text-gray-500">
            Available: {{ tokenBalanceFormatted }} W3B
          </p>
        </div>
        <button
          @click="handleDeposit"
          :disabled="!amount || isAnyPending"
          class="w-full py-3 bg-primary-600 hover:bg-primary-700 disabled:bg-gray-700 disabled:cursor-not-allowed text-white font-semibold rounded-lg transition-colors flex items-center justify-center space-x-2"
        >
          <svg v-if="isAnyPending" class="animate-spin h-5 w-5" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <span>{{ isAnyPending ? 'Processing...' : (needsApproval ? 'Approve & Deposit' : 'Deposit') }}</span>
        </button>
      </div>

      <div v-if="activeTab === 'withdraw'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Amount (W3B)</label>
          <input
            v-model="amount"
            type="number"
            step="0.0001"
            placeholder="0.0"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-primary-500"
          />
          <p class="mt-1 text-xs text-gray-500">
            Available: {{ bankBalanceFormatted }} W3B
          </p>
        </div>
        <button
          @click="handleWithdraw"
          :disabled="!amount || isAnyPending"
          class="w-full py-3 bg-primary-600 hover:bg-primary-700 disabled:bg-gray-700 disabled:cursor-not-allowed text-white font-semibold rounded-lg transition-colors flex items-center justify-center space-x-2"
        >
          <svg v-if="isAnyPending" class="animate-spin h-5 w-5" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <span>{{ isAnyPending ? 'Processing...' : 'Withdraw' }}</span>
        </button>
      </div>

      <div v-if="activeTab === 'transfer'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Recipient Address</label>
          <input
            v-model="recipient"
            type="text"
            placeholder="0x..."
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-primary-500"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Amount (W3B)</label>
          <input
            v-model="amount"
            type="number"
            step="0.0001"
            placeholder="0.0"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-primary-500"
          />
          <p class="mt-1 text-xs text-gray-500">
            Available: {{ bankBalanceFormatted }} W3B
          </p>
        </div>
        <button
          @click="handleTransfer"
          :disabled="!amount || !recipient || isAnyPending"
          class="w-full py-3 bg-primary-600 hover:bg-primary-700 disabled:bg-gray-700 disabled:cursor-not-allowed text-white font-semibold rounded-lg transition-colors flex items-center justify-center space-x-2"
        >
          <svg v-if="isAnyPending" class="animate-spin h-5 w-5" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <span>{{ isAnyPending ? 'Processing...' : 'Transfer' }}</span>
        </button>
      </div>

      <div v-if="txState.txHash" class="mt-4 p-4 bg-gray-900 rounded-lg border border-gray-700">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm text-gray-400">Transaction Hash</span>
          <button
            @click="txState.showDetails = !txState.showDetails"
            class="text-xs text-primary-400 hover:text-primary-300"
          >
            {{ txState.showDetails ? 'Hide Details' : 'Show Details' }}
          </button>
        </div>
        <p class="text-sm text-white font-mono break-all">{{ txState.txHash }}</p>
        
        <div v-if="txState.showDetails && txState.error" class="mt-4 space-y-2">
          <div class="text-xs text-gray-400">Error</div>
          <p class="text-sm text-red-400">{{ txState.error }}</p>
        </div>
      </div>

      <div v-if="showUserMessage" class="mt-4">
        <div
          :class="[
            'p-4 rounded-lg flex items-start justify-between',
            userMessage === 'Transaction confirmed!'
              ? 'bg-green-500/10 border border-green-500/50'
              : 'bg-red-500/10 border border-red-500/50'
          ]"
        >
          <p :class="userMessage === 'Transaction confirmed!' ? 'text-green-400' : 'text-red-400'">
            {{ userMessage }}
          </p>
          <button @click="dismissMessage" class="text-gray-400 hover:text-white ml-2">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
