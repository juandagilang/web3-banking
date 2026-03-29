export interface User {
  walletAddress: string
}

export interface AuthState {
  token: string | null
  isAuthenticated: boolean
  walletAddress: string | null
}

export interface WalletState {
  address: string | null
  isConnected: boolean
  chainId: number | null
}

export interface Transaction {
  id: number
  type: 'deposit' | 'withdrawal' | 'transfer'
  from: string
  to?: string
  amount: string
  block_number: number
  timestamp: number
  tx_hash: string
}

export interface Pagination {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  pagination?: Pagination
}

export interface NonceResponse {
  nonce: string
  message: string
}

export interface LoginResponse {
  token: string
  expires_in: number
}

export interface BalanceResponse {
  address: string
  balance: string
  symbol: string
}

export interface ContractInfo {
  token_address: string
  bank_address: string
  name: string
  symbol: string
  decimals: string
  total_supply: string
}

export interface TransactionState {
  isPending: boolean
  txHash: string | null
  confirmations: number
  error: string | null
  showDetails: boolean
}
