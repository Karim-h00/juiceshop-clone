import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteJuice } from "../api/juice";

export const useDeleteJuice = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (juiceID: string) => deleteJuice(juiceID),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["juices"] })
        }
    })
}