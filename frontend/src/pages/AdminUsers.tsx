
import type { userData } from "../types"
import { useGetAllUsers } from "../hooks/useGetAllUsers"
import { useAdminUpdateUser } from "../hooks/useAdminUpdateUser"
import { useDeleteUser } from "../hooks/useDeleteUser"
import { useEffect, useState } from "react"

function AdminUsers() {
    const [search, setSearch] = useState("")
    const [page, setPage] = useState(1)
    const [debouncedSearch, setDebouncedSearch] = useState("")

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(search)
            setPage(1)
        }, 300)
        return () => clearTimeout(timer)
    }, [search])

    const { data, isLoading, isFetching, isError } = useGetAllUsers({ search: debouncedSearch, page })

    const roleMutation = useAdminUpdateUser()
    const deleteMutation = useDeleteUser()

    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Users</h2>
            <input
                type="text"
                placeholder="Search by username or email..."
                value={search}
                onChange={(e) => {
                    setSearch(e.target.value)
                }}
                className="w-full max-w-sm px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-sm text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {isError && <div className="text-red-500">Failed to load users</div>}

            {isLoading && !data && <div className="text-gray-500">Loading initial users...</div>}

            {!isLoading && (!data || data.length === 0) && <div className="text-gray-500">No users found</div>}

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
                    {isFetching &&
                        <div className="absolute top-2 right-2 text-xs text-blue-500">
                            Updating...
                        </div>
                    }
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
                <div className="flex justify-between items-center">
                    <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                        className="rounded px-4 py-2 text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:text-gray-300"
                    >
                        Previous
                    </button>
                    <span className="text-sm text-gray-500 dark:text-gray-400">Page {page}</span>
                    <button
                        onClick={() => setPage(p => p + 1)}
                        disabled={!data || data.length < 50}
                        className="rounded px-4 py-2 text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:text-gray-300"
                    >
                        Next
                    </button>
                </div>
            </div>
        </div>
    )
}
export default AdminUsers