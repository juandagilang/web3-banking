import { defineStore } from 'pinia'
import { ref, readonly } from 'vue'
import { useWeb3, useWeb3Modal } from '@/composables/useWeb3'

export const useWalletStore = defineStore('wallet', () => {
  const { address, isConnected, chainId } = useWeb3()
  const { open: openModal } = useWeb3Modal()

  const showDetails = ref(false)

  const toggleDetails = () => {
    showDetails.value = !showDetails.value
  }

  const hideDetails = () => {
    showDetails.value = false
  }

  const formatAddress = (addr: string | null) => {
    if (!addr) return ''
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  return {
    address,
    isConnected,
    chainId,
    showDetails,
    openModal,
    toggleDetails,
    hideDetails,
    formatAddress,
  }
})
