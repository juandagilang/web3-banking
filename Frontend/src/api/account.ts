import apiClient from './client'
import type { ApiResponse, Transaction, Pagination } from '@/types'

interface TransactionsResponse {
  transactions: Transaction[]
}

export const getTransactions = async (
  address: string,
  page: number = 1,
  limit: number = 20
): Promise<{ data: Transaction[]; pagination: Pagination } | null> => {
  try {
    const response = await apiClient.get<ApiResponse<TransactionsResponse & { pagination: Pagination }>>(
      `/api/v1/transactions`,
      {
        params: { address, page, limit },
      }
    )
    if (response.data.success && response.data.data) {
      return {
        data: response.data.data.transactions,
        pagination: response.data.pagination!,
      }
    }
    return null
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
    return null
  }
}

export const getContractInfo = async () => {
  try {
    const response = await apiClient.get('/api/v1/contract/info')
    return response.data
  } catch (error) {
    console.error('Failed to fetch contract info:', error)
    return null
  }
}
