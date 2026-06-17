import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteUser } from "../api/users";

export const useDeleteUser = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (id: string) => deleteUser(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["getUsers"] })
        }
    })
}