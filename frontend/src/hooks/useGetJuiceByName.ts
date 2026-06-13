import { useQuery } from "@tanstack/react-query";
import { getJuiceByName } from "../api/juice";

export const useGetJuiceByName = (juiceName: string) =>{
    return useQuery({
        queryKey: ["getJuiceByName", juiceName],
        queryFn: ()=> getJuiceByName(juiceName)
    })
}