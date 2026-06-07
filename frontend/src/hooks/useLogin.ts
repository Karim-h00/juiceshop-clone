import { useMutation } from "@tanstack/react-query";
import { login } from "../api/auth";
import { type LoginCredentials } from "../types";

export const useLogin = () => {
    return useMutation({
        mutationFn: (credentials: LoginCredentials) => login(credentials),
    })
}