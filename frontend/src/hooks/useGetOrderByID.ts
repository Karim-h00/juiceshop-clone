import { useQuery } from "@tanstack/react-query";
import { getOrderByID } from "../api/order";

export const useGetOrderByID = (orderID: string) =>{
    return useQuery({
        queryKey: ["getOrderByID", orderID],
        queryFn: ()=> getOrderByID(orderID)
    })
}