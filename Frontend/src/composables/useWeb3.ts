import { ref, readonly } from 'vue'
import { createWeb3Modal, defaultWagmiConfig } from '@web3modal/wagmi'
import { mainnet, sepolia, hardhat } from 'viem/chains'
import { reconnect, watchAccount, watchNetwork } from '@wagmi/core'
import { configureChains, createConfig, publicProvider } from 'wagmi'
import type { Config } from '@wagmi/core'

const RPC_URL = import.meta.env.VITE_RPC_URL || 'http://127.0.0.1:8545'
const CHAIN_ID = Number(import.meta.env.VITE_CHAIN_ID) || 31337

const chains = [mainnet, sepolia, hardhat]

const { publicClient, webSocketPublicClient } = configureChains(
  chains,
  [publicProvider()]
)

const wagmiConfig: Config = defaultWagmiConfig({
  chains,
  projectId: 'web3bank-demo',
  metadata: {
    name: 'Web3 Banking',
    description: 'Decentralized banking system',
    url: 'http://localhost:5173',
    icons: ['https://avatars.githubusercontent.com/u/37784886'],
  },
  publicClient,
  webSocketPublicClient,
})

export const web3Modal = createWeb3Modal({
  wagmiConfig,
  chains,
  projectId: 'web3bank-demo',
  enableAnalytics: false,
})

reconnect(wagmiConfig)

const address = ref<string | null>(null)
const isConnected = ref(false)
const chainId = ref<number | null>(null)

watchAccount(wagmiConfig, {
  onChange: (account) => {
    address.value = account.address ?? null
    isConnected.value = account.isConnected
  },
})

watchNetwork(wagmiConfig, {
  onChange: (network) => {
    chainId.value = network.chain?.id ?? null
  },
})

export function useWeb3Modal() {
  const open = async () => {
    web3Modal.open()
  }

  const close = () => {
    web3Modal.close()
  }

  return {
    open,
    close,
  }
}

export function useWeb3() {
  return {
    address: readonly(address),
    isConnected: readonly(isConnected),
    chainId: readonly(chainId),
  }
}

export { wagmiConfig }
