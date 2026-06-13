import { useMutation } from "@tanstack/react-query";
import { addJuice } from "../api/juice";
import { type juiceData } from "../types";

export const useAddJuice = () =>{
    return useMutation({
        mutationFn: (juiceData: juiceData) => addJuice(juiceData),
    })
}