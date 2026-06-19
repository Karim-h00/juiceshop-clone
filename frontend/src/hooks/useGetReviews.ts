import { useQuery } from "@tanstack/react-query";
import { getReviews } from "../api/reviews";

export const useGetReviews = (slug: string) => {
    return useQuery({
        queryKey: ["getReviews", slug],
        queryFn: () => getReviews(slug)
    })
}