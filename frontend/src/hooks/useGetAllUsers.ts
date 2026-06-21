import { useQuery } from "@tanstack/react-query";
import { getAllUsers } from "../api/users";

export const useGetAllUsers = ({search = "", page = 1}) =>{
    return useQuery({
        queryKey: ["getUsers", search, page],
        queryFn: ()=> getAllUsers({search, page})
    })
}