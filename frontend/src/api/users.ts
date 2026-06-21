import { useAuthStore } from "../store/authStore"
import type { userData } from "../types"
import { BASE_URL } from "./config"

export const getAllUsers = async ({ search }: { search: string }): Promise<userData[]> => {
    const token = useAuthStore.getState().token
    const params = new URLSearchParams()
    if (search) params.set("q", search)
    const response = await fetch(`${BASE_URL}/api/admin/users?${params}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
        throw new Error('Failed to get users')
    }
    const data = response.json()
    console.log(data)
    return data
}

export const adminUpdateUser = async (id: string, role: "admin" | "user") => {
    const token = useAuthStore.getState().token
    const response = await fetch(`${BASE_URL}/api/admin/users/${id}/role`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ role })
    })
    if (!response.ok) {
        throw new Error('Failed to update user')
    }
}

export const deleteUser = async (id: string) => {
    const token = useAuthStore.getState().token
    const response = await fetch(`${BASE_URL}/api/admin/users/${id}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
        throw new Error('Failed to delete user')
    }
}