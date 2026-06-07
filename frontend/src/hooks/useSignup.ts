import { useMutation } from "@tanstack/react-query";
import { signup } from "../api/auth";
import type { SigunpCredentials } from "../types";

export const useSignup = () => {
    return useMutation({
        mutationFn: (credentials: SigunpCredentials) => signup(credentials),
    })
}