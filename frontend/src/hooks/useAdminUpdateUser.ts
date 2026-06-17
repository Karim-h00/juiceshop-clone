import { useMutation, useQueryClient } from "@tanstack/react-query";
import { adminUpdateUser } from "../api/users";

type updateUserVars = {
    id: string
    role: "admin"|"user"
}
export const useAdminUpdateUser = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({id, role}: updateUserVars) => adminUpdateUser(id, role),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["getUsers"] })
    })
}