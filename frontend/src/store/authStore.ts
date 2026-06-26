import { create } from 'zustand'
import { getMe, refresh } from '../api/auth'

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
    initialize: () => Promise<void>
}

const decodeToken = (token: string) => {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
    return payload
}

export const useAuthStore = create<AuthStore>((set, get) => ({
    token: null,
    user: null,
    isLoading: true,
    setAuth: (token, user) => set({ token, user, isLoading: false }),
    clearAuth: () => set({ user: null, isLoading: false }),
    initialize: async () => {
        if (get().token) {
            set({ isLoading: false })
            return
        }
        try {
            const token = await refresh()
            const payload = decodeToken(token)
            const me = await getMe(token)
            set({
                token,
                user: { username: me.username, email: me.email, role: payload.role },
                isLoading: false,
            })
        } catch {
            get().clearAuth()
        }
    }
}))