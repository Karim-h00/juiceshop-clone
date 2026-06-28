import { useAuthStore } from '../store/authStore'
import { BASE_URL } from './config'
import { type JuiceData, type JuiceUpdateParams } from '../types'

export const getJuices = async (): Promise<JuiceData[]> => {
    const response = await fetch(`${BASE_URL}/api`)
    if (!response.ok) {
        throw new Error('Failed to fetch juices')
    }
    return response.json()
}

export const getJuiceByName = async (slug: string): Promise<JuiceData> => {
    const response = await fetch(`${BASE_URL}/api/juice/${slug}`)
    if (!response.ok) {
        throw new Error('Failed to fetch juice data')
    }
    return response.json()
}

export const deleteJuice = async (juiceID: string) => {
    const token = useAuthStore.getState().token

    const response = await fetch(`${BASE_URL}/api/admin/juice/${juiceID}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
    })
    if (!response.ok) {
        throw new Error('Failed to delete juice')
    }
    return response.json()
}

export const addJuice = async (juiceData: JuiceUpdateParams) => {
    const token = useAuthStore.getState().token
    console.log(JSON.stringify(juiceData))

    const response = await fetch(`${BASE_URL}/api/admin/juice`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: juiceData.name,
            description: juiceData.description,
            price: juiceData.price,
            stock: juiceData.stock
        })
    })
    if (!response.ok) {
        throw new Error('Failed to update juice data')
    }
}

export const updateJuice = async (id: string, juiceData: JuiceUpdateParams) => {
    const token = useAuthStore.getState().token

    const response = await fetch(`${BASE_URL}/api/admin/juice/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
            name: juiceData.name,
            description: juiceData.description,
            price: juiceData.price,
            stock: juiceData.stock
        })
    })
    if (!response.ok) {
        throw new Error('Failed to update juice data')
    }
}

export const uploadJuiceImage = async (id: string, file: File) => {
    const token = useAuthStore.getState().token

    const formData = new FormData()
    formData.append("juiceImage", file)

    const response = await fetch(`${BASE_URL}/api/admin/juice/${id}/image`, {
        method: "PUT",
        headers: {
            "Authorization": `Bearer ${token}`
        },
        body: formData
    })
    if (!response.ok) {
        throw new Error('Failed to update juice image')
    }

}