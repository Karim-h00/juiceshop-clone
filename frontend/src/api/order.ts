import { BASE_URL } from './config'
import { useAuthStore } from '../store/authStore'
import { useCartStore } from '../store/cartStore'

export const checkout = async () => {

    const token = useAuthStore.getState().token
    const items = useCartStore.getState().items

    console.log("token:", token)
    console.log("items:", items)
    console.log("body:", JSON.stringify(items))

    const response = await fetch(`${BASE_URL}/api/order`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            items: items.map(item => ({
                juice_id: item.id,
                quantity: item.quantity
            }))
        })

    },)
    if (!response.ok) {
        throw new Error('Failed to checkout')
    }
    const data = await response.json
    return data
}