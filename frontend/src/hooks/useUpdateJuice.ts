import { useMutation } from "@tanstack/react-query";
import { updateJuice } from "../api/juice";
import { type JuiceUpdateParams } from "../types";

type UpdateJuiceVariables = {
  id: string
  juiceData: JuiceUpdateParams
}

export const useUpdateJuice = () =>{
    return useMutation({
        mutationFn: ({id, juiceData}: UpdateJuiceVariables) => updateJuice(id, juiceData),
    })
}