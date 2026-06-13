import { useMutation } from "@tanstack/react-query";
import { updateJuice } from "../api/juice";
import { type JuiceData } from "../types";

export const useUpdateJuice = () =>{
    return useMutation({
        mutationFn: (juiceData: JuiceData) => updateJuice(juiceData),
    })
}