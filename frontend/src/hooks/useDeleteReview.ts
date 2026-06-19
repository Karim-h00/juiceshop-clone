import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteReview } from "../api/reviews";

export const useDeleteReview = (slug: string) => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (reviewID: string)=> deleteReview(reviewID),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["getReviews", slug] })
        }
    })
}