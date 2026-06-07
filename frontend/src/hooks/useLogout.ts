import { useMutation } from "@tanstack/react-query"
import { logout } from "../api/auth"
import useAuthStore from "../store/authStore";


export const useLogout = () => {
    const refreshToken = useAuthStore((s) => s.refreshToken);
    const clearStore = useAuthStore((s) => s.logout);
    return useMutation({
        mutationFn: () => {
            if (!refreshToken) throw new Error("No refresh token");
            return logout(refreshToken);
        },
        onSuccess:() => {
            clearStore()
        }
    })
}