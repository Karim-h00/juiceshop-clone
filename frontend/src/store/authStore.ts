import { create } from 'zustand'

type User = {
    username: string
    email: string
    role: string
}

type AuthStore = {
    token: string | null
    refreshToken: string | null
    user: User | null
    setAuthStore: (
        token: string,
        refreshToken: string,
        username: string,
        email: string,
        role: string
    ) => void
    logout: () => void
}

const useAuthStore = create<AuthStore>((set) => ({
    token: null,
    refreshToken: null,
    user: null,
    setAuthStore(token, refreshToken, username, email, role){
        set({
            token, refreshToken,
            user: {username, email, role}
        })
    },
    logout: () => set({ token: null, refreshToken: null, user: null }),
}))

export default useAuthStore