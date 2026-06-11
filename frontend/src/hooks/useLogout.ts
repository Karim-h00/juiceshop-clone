import { useMutation } from "@tanstack/react-query"
import { logout } from "../api/auth"
import { useAuthStore } from "../store/authStore";


export const useLogout = () => {
  const clearAuth = useAuthStore((s) => s.clearAuth)
  return useMutation({
    mutationFn: () => logout(),
    onSuccess: () => {
      clearAuth()
    },
    onError: (error) => {
      console.log('logout error', error)
      clearAuth()
    }
  })
}