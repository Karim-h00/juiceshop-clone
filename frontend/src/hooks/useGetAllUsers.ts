import { useQuery } from "@tanstack/react-query";
import { getAllUsers } from "../api/users";

export const useGetAllUsers = ({search = ""}) =>{
    return useQuery({
        queryKey: ["getUsers", search],
        queryFn: ()=> getAllUsers({search})
    })
}