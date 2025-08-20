import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'

export const api = axios.create({ baseURL })

export type TransferReq = { from: string; to: string; amount: number }
export type Transaction = { From: string; To: string; Amount: number; Time: number }
export type Block = { Index: number; Timestamp: number; Transactions: Transaction[]; PrevHash: string; Hash: string; Nonce: number }

export const getBlockchain = () => api.get<Block[]>('/blockchain').then(r => r.data)
export const getBalance = (user: string) => api.get<{ user: string; balance: number }>(`/balance/${user}`).then(r => r.data)
export const getPending = () => api.get<Transaction[]>(`/pending`).then(r => r.data)
export const postTransfer = (body: TransferReq) => api.post<{ status: string }>(`/transfer`, body).then(r => r.data)









