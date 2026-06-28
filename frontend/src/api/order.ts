import { BASE_URL } from './config'
import { useAuthStore } from '../store/authStore'
import { useCartStore } from '../store/cartStore'
import type { OrderDetail } from '../types'

export const checkout = async () => {

    const token = useAuthStore.getState().token
    const items = useCartStore.getState().items

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
    return response.json()
}

export const getOrderHistory = async () => {
    const token = useAuthStore.getState().token

    const response = await fetch(`${BASE_URL}/api/order`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
        throw new Error('Failed to get order history')
    }
    return response.json()
}

export const getOrderByID = async (orderID: string): Promise<OrderDetail> => {
    const token = useAuthStore.getState().token

    const response = await fetch(`${BASE_URL}/api/order/${orderID}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
        throw new Error('Failed to get order')
    }
    return response.json()
}

export const getAdminOrders = async ({ search, page }: { search: string, page: number }) =>  {
    const token = useAuthStore.getState().token
    const url = new URL(`${BASE_URL}/api/admin/orders`)
    if (page) url.searchParams.set("page", page.toString())
    if (search) url.searchParams.set("search", search)
    const response = await fetch(url, {
        headers: {
            "Content-Type": 'application/json',
            "Authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
        throw new Error('Failed to get orders')
    }
    return response.json()
}