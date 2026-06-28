import { useAuthStore } from '../store/authStore'
import { BASE_URL } from './config'

export const getAuditLogs = async () => {
    const token = useAuthStore.getState().token
    const response = await fetch(`${BASE_URL}/api/admin/audits`, {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
            "authorization": `Bearer ${token}`
        }
    })
    if (!response.ok) {
        throw new Error('Failed to get audit logs')
    }
    return response.json()
}