import { useQuery } from "@tanstack/react-query";
import { getAdminOrders } from "../api/order";

export const useGetAdminOrders = ({search = "", page = 1}) =>{
    return useQuery({
        queryKey: ["adminOrders", search, page],
        queryFn: ()=> getAdminOrders({search, page}),
        placeholderData: (prevData) => prevData
    })
}