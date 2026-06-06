import { BASE_URL } from './config'
import { type LoginCredentials, type LoginRes, type SigunpCredentials, type SignupRes } from '../types'

export const login = async (credentials: LoginCredentials): Promise<LoginRes> => {
    const response = await fetch(`${BASE_URL}/api/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
    })
    if (!response.ok) {
        throw new Error('Login failed')
    }
    return response.json()
}

export const signup = async (credentials: SigunpCredentials): Promise<SignupRes> => {
    const response = await fetch (`${BASE_URL}/api/signup`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials)
    })
    if (!response.ok){
        throw new Error('signup failed')
    }
    return response.json()
}