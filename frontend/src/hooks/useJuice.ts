import { useQuery } from "@tanstack/react-query";
import { getJuices } from "../api/juice"; 

export const useJuice = () =>{
    return useQuery({
        queryKey: ["juices"],
        queryFn: getJuices
    })
}