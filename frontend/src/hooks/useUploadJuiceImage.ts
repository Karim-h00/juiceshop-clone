import { useMutation } from "@tanstack/react-query";
import { uploadJuiceImage } from "../api/juice";

export const useUploadJuiceImage = () => {
    return useMutation({
        mutationFn: ({id, file}: {id: string; file: File}) =>
            uploadJuiceImage(id, file)
    })
}