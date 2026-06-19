
import type { userData } from "../types"
import { useGetAllUsers } from "../hooks/useGetAllUsers"
import { useAdminUpdateUser } from "../hooks/useAdminUpdateUser"
import { useDeleteUser } from "../hooks/useDeleteUser"

function AdminUsers() {

    const { data, isLoading, isError } = useGetAllUsers()
    const roleMutation = useAdminUpdateUser()
    const deleteMutation = useDeleteUser()

    if (isLoading) return <div className="text-gray-500 dark:text-gray-400">Loading...</div>
    if (isError) return <div className="text-red-500">Failed to load users</div>
    if (!data) return <div>no users found</div>

    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Users</h2>
            <div className="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700">
                <table className="w-full text-sm">
                    <thead className="bg-gray-50 dark:bg-gray-800 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        <tr>
                            <th className="px-4 py-3">User ID</th>
                            <th className="px-4 py-3">Username</th>
                            <th className="px-4 py-3">Email</th>
                            <th className="px-4 py-3">Role</th>
                            <th className="px-4 py-3">Created At</th>
                            <th className="px-4 py-3">Updated At</th>
                            <th className="px-4 py-3">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-900">
                        {data?.map((user: userData) => {
                            const isAdmin = user.role === "admin"
                            return (
                                <tr key={user.id} className="text-gray-700 dark:text-gray-300">
                                    <td className="px-4 py-3 font-mono text-xs text-gray-400 dark:text-gray-500">
                                        {user.id.slice(0, 8) ?? "-"}…
                                    </td>
                                    <td className="px-4 py-3">{user.username}</td>
                                    <td className="px-4 py-3">{user.email}</td>
                                    <td className="px-4 py-3">
                                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${isAdmin
                                            ? "bg-emerald-100 text-emerald-700 dark:bg-emerald-900 dark:text-emerald-300"
                                            : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400"
                                            }`}>
                                            {user.role}
                                        </span>
                                    </td>
                                    <td className="px-4 py-3 text-gray-400 dark:text-gray-500">
                                        {new Date(user.created_at).toLocaleDateString()}
                                    </td>
                                    <td className="px-4 py-3 text-gray-400 dark:text-gray-500">
                                        {new Date(user.updated_at).toLocaleDateString()}
                                    </td>
                                    <td className="px-4 py-3">
                                        <div className="flex gap-2">
                                            <button
                                                onClick={() => roleMutation.mutate({ id: user.id, role: isAdmin ? "user" : "admin" })}
                                                className="rounded-lg px-3 py-1 text-xs font-medium bg-emerald-50 text-emerald-700 hover:bg-emerald-100 dark:bg-emerald-900/30 dark:text-emerald-400 dark:hover:bg-emerald-900/50 transition-colors cursor-pointer"
                                            >
                                                {isAdmin ? "Revoke admin" : "Make admin"}
                                            </button>
                                            <button
                                                onClick={() => {
                                                    if (confirm(`Delete ${user.username}?`)) {
                                                        deleteMutation.mutate(user.id)
                                                    }
                                                }}
                                                className="rounded-lg px-3 py-1 text-xs font-medium bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/30 dark:text-red-400 dark:hover:bg-red-900/50 transition-colors cursor-pointer"
                                            >
                                                Delete
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            )
                        })}
                    </tbody>
                </table>
            </div>
        </div>
    )
}
export default AdminUsers