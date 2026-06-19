import { useMutation, useQueryClient } from "@tanstack/react-query";
import { addReview } from "../api/reviews";

export const useAddReview = (slug: string) => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ rating, comment }: { rating: number; comment?: string }) =>
            addReview(slug, rating, comment),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["getReviews", slug] })
        }
    })
}