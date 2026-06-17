import { useQuery } from "@tanstack/react-query";
import { getAllUsers } from "../api/users";

export const useGetAllUsers = () =>{
    return useQuery({
        queryKey: ["getUsers"],
        queryFn: ()=> getAllUsers()
    })
}