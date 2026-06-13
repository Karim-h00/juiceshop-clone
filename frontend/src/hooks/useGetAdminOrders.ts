import { useQuery } from "@tanstack/react-query";
import { getAdminOrders } from "../api/order";

export const useGetAdminOrders = (page?: number) =>{
    return useQuery({
        queryKey: ["adminOrders", page],
        queryFn: ()=> getAdminOrders(page)
    })
}