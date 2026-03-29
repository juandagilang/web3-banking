import { ref, readonly } from 'vue'
import { parseEther, formatEther } from 'viem'
import { usePublicClient, useWalletClient } from '@wagmi/vue'
import { useContractReads, useContractWrite, useWaitForTransactionReceipt } from '@wagmi/vue'
import { wagmiConfig } from './useWeb3'

const TOKEN_ADDRESS = import.meta.env.VITE_TOKEN_ADDRESS as `0x${string}`
const BANK_ADDRESS = import.meta.env.VITE_BANK_ADDRESS as `0x${string}`

const abi = [
  {
    name: 'balanceOf',
    type: 'function',
    inputs: [{ name: 'user', type: 'address' }],
    outputs: [{ type: 'uint256' }],
    stateMutability: 'view',
  },
  {
    name: 'deposit',
    type: 'function',
    inputs: [{ name: 'amount', type: 'uint256' }],
    outputs: [],
    stateMutability: 'nonpayable',
  },
  {
    name: 'withdraw',
    type: 'function',
    inputs: [{ name: 'amount', type: 'uint256' }],
    outputs: [],
    stateMutability: 'nonpayable',
  },
  {
    name: 'transfer',
    type: 'function',
    inputs: [
      { name: 'to', type: 'address' },
      { name: 'amount', type: 'uint256' },
    ],
    outputs: [],
    stateMutability: 'nonpayable',
  },
] as const

const tokenAbi = [
  {
    name: 'approve',
    type: 'function',
    inputs: [
      { name: 'spender', type: 'address' },
      { name: 'amount', type: 'uint256' },
    ],
    outputs: [{ type: 'bool' }],
    stateMutability: 'nonpayable',
  },
  {
    name: 'allowance',
    type: 'function',
    inputs: [
      { name: 'owner', type: 'address' },
      { name: 'spender', type: 'address' },
    ],
    outputs: [{ type: 'uint256' }],
    stateMutability: 'view',
  },
  {
    name: 'balanceOf',
    type: 'function',
    inputs: [{ name: 'account', type: 'address' }],
    outputs: [{ type: 'uint256' }],
    stateMutability: 'view',
  },
] as const

export interface TransactionState {
  isPending: boolean
  isConfirming: boolean
  txHash: string | null
  confirmations: number
  error: string | null
  showDetails: boolean
}

export function useContract(address: string) {
  const publicClient = usePublicClient({ config: wagmiConfig })
  const { data: walletClient } = useWalletClient({ config: wagmiConfig })

  const txState = ref<TransactionState>({
    isPending: false,
    isConfirming: false,
    txHash: null,
    confirmations: 0,
    error: null,
    showDetails: false,
  })

  const { data: bankBalance, refetch: refetchBankBalance } = useContractReads({
    config: wagmiConfig,
    contracts: [
      {
        address: BANK_ADDRESS,
        abi,
        functionName: 'balanceOf',
        args: [address as `0x${string}`],
      },
    ],
    query: {
      enabled: !!address,
    },
  })

  const { data: tokenBalance, refetch: refetchTokenBalance } = useContractReads({
    config: wagmiConfig,
    contracts: [
      {
        address: TOKEN_ADDRESS,
        abi: tokenAbi,
        functionName: 'balanceOf',
        args: [address as `0x${string}`],
      },
    ],
    query: {
      enabled: !!address,
    },
  })

  const { data: allowance, refetch: refetchAllowance } = useContractReads({
    config: wagmiConfig,
    contracts: [
      {
        address: TOKEN_ADDRESS,
        abi: tokenAbi,
        functionName: 'allowance',
        args: [address as `0x${string}`, BANK_ADDRESS],
      },
    ],
    query: {
      enabled: !!address,
    },
  })

  const { writeAsync: approve, data: approveHash } = useContractWrite({
    config: wagmiConfig,
    mutation: {
      onSuccess: () => {
        txState.value.txHash = approveHash.value ?? null
        txState.value.isPending = true
      },
      onError: (error) => {
        txState.value.error = error.message
        txState.value.isPending = false
      },
    },
  })

  const { writeAsync: deposit, data: depositHash } = useContractWrite({
    config: wagmiConfig,
    mutation: {
      onSuccess: () => {
        txState.value.txHash = depositHash.value ?? null
        txState.value.isPending = true
      },
      onError: (error) => {
        txState.value.error = error.message
        txState.value.isPending = false
      },
    },
  })

  const { writeAsync: withdraw, data: withdrawHash } = useContractWrite({
    config: wagmiConfig,
    mutation: {
      onSuccess: () => {
        txState.value.txHash = withdrawHash.value ?? null
        txState.value.isPending = true
      },
      onError: (error) => {
        txState.value.error = error.message
        txState.value.isPending = false
      },
    },
  })

  const { writeAsync: transfer, data: transferHash } = useContractWrite({
    config: wagmiConfig,
    mutation: {
      onSuccess: () => {
        txState.value.txHash = transferHash.value ?? null
        txState.value.isPending = true
      },
      onError: (error) => {
        txState.value.error = error.message
        txState.value.isPending = false
      },
    },
  })

  const { isLoading: isApproveConfirming, isSuccess: isApproveConfirmed } = useWaitForTransactionReceipt({
    hash: approveHash,
  })

  const { isLoading: isDepositConfirming, isSuccess: isDepositConfirmed } = useWaitForTransactionReceipt({
    hash: depositHash,
  })

  const { isLoading: isWithdrawConfirming, isSuccess: isWithdrawConfirmed } = useWaitForTransactionReceipt({
    hash: withdrawHash,
  })

  const { isLoading: isTransferConfirming, isSuccess: isTransferConfirmed } = useWaitForTransactionReceipt({
    hash: transferHash,
  })

  const resetTxState = () => {
    txState.value = {
      isPending: false,
      isConfirming: false,
      txHash: null,
      confirmations: 0,
      error: null,
      showDetails: false,
    }
  }

  const doApprove = async (amount: string) => {
    try {
      resetTxState()
      txState.value.isPending = true
      await approve({
        address: TOKEN_ADDRESS,
        abi: tokenAbi,
        functionName: 'approve',
        args: [BANK_ADDRESS, parseEther(amount)],
      })
    } catch (error: unknown) {
      txState.value.error = error instanceof Error ? error.message : 'Approval failed'
      txState.value.isPending = false
    }
  }

  const doDeposit = async (amount: string) => {
    try {
      resetTxState()
      txState.value.isPending = true
      await deposit({
        address: BANK_ADDRESS,
        abi,
        functionName: 'deposit',
        args: [parseEther(amount)],
      })
    } catch (error: unknown) {
      txState.value.error = error instanceof Error ? error.message : 'Deposit failed'
      txState.value.isPending = false
    }
  }

  const doWithdraw = async (amount: string) => {
    try {
      resetTxState()
      txState.value.isPending = true
      await withdraw({
        address: BANK_ADDRESS,
        abi,
        functionName: 'withdraw',
        args: [parseEther(amount)],
      })
    } catch (error: unknown) {
      txState.value.error = error instanceof Error ? error.message : 'Withdrawal failed'
      txState.value.isPending = false
    }
  }

  const doTransfer = async (to: string, amount: string) => {
    try {
      resetTxState()
      txState.value.isPending = true
      await transfer({
        address: BANK_ADDRESS,
        abi,
        functionName: 'transfer',
        args: [to as `0x${string}`, parseEther(amount)],
      })
    } catch (error: unknown) {
      txState.value.error = error instanceof Error ? error.message : 'Transfer failed'
      txState.value.isPending = false
    }
  }

  const refreshBalances = () => {
    refetchBankBalance()
    refetchTokenBalance()
    refetchAllowance()
  }

  const formatBalance = (balance: bigint | undefined) => {
    if (!balance) return '0'
    return formatEther(balance)
  }

  return {
    bankBalance,
    tokenBalance,
    allowance,
    txState: readonly(txState),
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
  }
}
