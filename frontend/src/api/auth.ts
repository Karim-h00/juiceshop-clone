import { BASE_URL } from './config'
import { type LoginCredentials, type LoginRes, type SigunpCredentials, type SignupRes, type MeRes } from '../types'

export const login = async (credentials: LoginCredentials): Promise<LoginRes> => {
    const response = await fetch(`${BASE_URL}/api/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
        credentials: "include"
    })
    if (!response.ok) {
        throw new Error('Login failed')
    }
    return response.json()
}

export const signup = async (credentials: SigunpCredentials): Promise<SignupRes> => {
    const response = await fetch(`${BASE_URL}/api/users`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
        credentials: 'include',
    })
    if (!response.ok) {
        throw new Error('signup failed')
    }
    return response.json()
}

export const logout = async (): Promise<void> => {
    const response = await fetch(`${BASE_URL}/api/logout`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
    })
    if (!response.ok) {
        throw new Error('logout failed')
    }
}

export const refresh = async (): Promise<string> => {
    const response = await fetch(`${BASE_URL}/api/refresh`, {
        method: 'POST',
        credentials: 'include',
    })
    if (!response.ok) throw new Error('Refresh failed')
    const data = await response.json()
    return data.token
}

export const getMe = async (token: string): Promise<MeRes> => {
    const response = await fetch(`${BASE_URL}/api/me`, {
        headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!response.ok) throw new Error('Not authenticated')
    return response.json()
}

export const checkPassword = async (password: string): Promise<{ pwned: boolean }> => {
    const res = await fetch(`${BASE_URL}/api/check-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password })
    })
    if (!res.ok) throw new Error('Failed to check password')
    return res.json()
}