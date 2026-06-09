import { create } from 'zustand'

type User = {
    username: string
    email: string
    role: string
}

type AuthStore = {
    user: User | null
    isLoading: boolean
    setUser:(user: User) =>  void
    clearAuth: () => void
}

const useAuthStore = create<AuthStore>((set) => ({
    user: null,
    isLoading: true,
    setUser: (user)=>set({user, isLoading: false}),
    clearAuth: () => set({ user: null, isLoading: false }),
}))

export default useAuthStore