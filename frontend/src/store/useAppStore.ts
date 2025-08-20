import { create } from 'zustand'

type AppState = {
  from: string
  to: string
  amount: number
  setFrom: (v: string) => void
  setTo: (v: string) => void
  setAmount: (v: number) => void
}

export const useAppStore = create<AppState>((set) => ({
  from: 'alice',
  to: 'bob',
  amount: 100,
  setFrom: (v) => set({ from: v }),
  setTo: (v) => set({ to: v }),
  setAmount: (v) => set({ amount: v }),
}))









