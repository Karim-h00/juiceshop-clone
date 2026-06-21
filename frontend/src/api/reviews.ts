import { useAuthStore } from "../store/authStore"
import { BASE_URL } from "./config"

export const addReview = async (slug: string, rating: number, comment?: string) => {
    const token = useAuthStore.getState().token

    const response = await fetch(`${BASE_URL}/api/juice/${slug}/reviews`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({rating, comment})
        })
        if (!response.ok) {
            const err = await response.json()
            throw new Error(err.error || "failed to add review")
        }
        return response.json()
}

export const getReviews = async (slug: string) => {

     const response = await fetch(`${BASE_URL}/api/juice/${slug}/reviews`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
        if (!response.ok) {
            throw new Error('Failed to get reviews')
        }
        return response.json()
}

export const deleteReview = async(id: string) => {
    const token = useAuthStore.getState().token
    const response = await fetch(`${BASE_URL}/api/review/${id}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
            throw new Error('Failed to delete review')
        }
}