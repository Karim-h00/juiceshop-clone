import { create } from 'zustand'

type User = {
    id: string
    username: string
    email: string
    role: string
}

type AuthStore = {
    token: string | null
    refreshToken: string | null
    user: User | null
    setToken: (token: string) => void
    setRefreshToken: (token: string) => void
    setUser: (user: User) => void
    logout: () => void
}

const useAuthStore = create<AuthStore>((set) => ({
    token: null,
    refreshToken: null,
    user: null,
    setToken: (token) => set({ token }),
    setRefreshToken: (token) => set({ refreshToken: token }),
    setUser: (user) => set({ user }),
    logout: () => set({ token: null, refreshToken: null, user: null }),
}))

export default useAuthStore