import { useMutation } from "@tanstack/react-query";
import { updateJuice } from "../api/juice";
import { type juiceData } from "../types";

export const useUpdateJuice = () =>{
    return useMutation({
        mutationFn: (juiceData: juiceData) => updateJuice(juiceData),
    })
}