import { useQuery } from "@tanstack/react-query";
import { getOrderHistory } from "../api/order";

export const useGetOrderHistory = () =>{
    return useQuery({
        queryKey: ["getOrderHistory"],
        queryFn: getOrderHistory
    })
}