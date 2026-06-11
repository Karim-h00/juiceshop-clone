import { create } from 'zustand'

type User = {
    username: string
    email: string
    role: string
}

type AuthStore = {
    token: string | null
    user: User | null
    isLoading: boolean
    setAuth: (token: string, user: User) => void
    clearAuth: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
    token: null,
    user: null,
    isLoading: true,
    setAuth: (token, user) => set({ token, user, isLoading: false }),
    clearAuth: () => set({ user: null, isLoading: false }),
}))