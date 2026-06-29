import { useMutation } from '@tanstack/react-query'
import { checkPassword } from '../api/auth'

export const useCheckPassword = () => {
    return useMutation({
        mutationFn: checkPassword
    })
}