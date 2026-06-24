import { useMutation } from "@tanstack/react-query";
import { changePassword } from "../api/users";

type ChangePassword = {
    password: string
    new_password: string
}

export const useChangePassword = () => {
    return useMutation({
        mutationFn: ({password, new_password}: ChangePassword) => changePassword(password, new_password),
    })
}