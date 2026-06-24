import { useMutation } from "@tanstack/react-query";
import { updateUserData } from "../api/users";

type UpdateUserData = {
    username: string
    email: string
}

export const useUpdateUserData = () => {
    return useMutation({
        mutationFn: ({username, email}: UpdateUserData) => updateUserData(username, email),
    })
}