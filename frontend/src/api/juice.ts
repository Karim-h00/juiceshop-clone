import { BASE_URL } from './config'

export const getJuices = async () => {
    const response = await fetch(`${BASE_URL}/api`)
    if (!response.ok) {
        throw new Error('Failed to fetch juices')
    }
    return response.json()
}