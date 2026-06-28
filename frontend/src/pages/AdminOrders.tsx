import { useEffect, useState } from "react"
import { useGetAdminOrders } from "../hooks/useGetAdminOrders"
import type { AdminOrder } from "../types"
import { useNavigate } from "react-router-dom"

function AdminOrders() {
  const [search, setSearch] = useState("")
  const [debouncedSearch, setDebouncedSearch] = useState("")
  const [page, setPage] = useState(1)
  const [copiedId, setCopiedId] = useState<string | null>(null)

  useEffect(() => {
          const timer = setTimeout(() => {
              setDebouncedSearch(search)
              setPage(1)
          }, 300)
          return () => clearTimeout(timer)
      }, [search])

  const { data, isLoading, isError } = useGetAdminOrders({search: debouncedSearch, page})

  const navigate = useNavigate();

  const handleCopy = (id: string) => {
    navigator.clipboard.writeText(id)
    setCopiedId(id)
    setTimeout(() => setCopiedId(null), 2000)
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Orders</h2>
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
              <th className="px-4 py-3">Order ID</th>
              <th className="px-4 py-3">User</th>
              <th className="px-4 py-3">Total</th>
              <th className="px-4 py-3">Date</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-900">
            {data?.map((order: AdminOrder) => (
              <tr key={order.ID} className="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                onClick={() => navigate(`/admin/orders/${order.ID}`)}>
                <td
                  className="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs cursor-pointer hover:text-emerald-500 transition-colors"
                  onClick={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    handleCopy(order.ID)
                  }}
                  title="Click to copy full ID"
                >
                  {copiedId === order.ID ? (
                    <span className="text-emerald-500">Copied!</span>
                  ) : (
                    `${order.ID.slice(0, 8)}...`
                  )}
                </td>
                <td className="px-4 py-3 text-gray-900 dark:text-white">{order.Username}</td>
                <td className="px-4 py-3 text-gray-700 dark:text-gray-300">${(order.Total / 100).toFixed(2)}</td>
                <td className="px-4 py-3 text-gray-500 dark:text-gray-400">{new Date(order.CreatedAt).toLocaleDateString()}</td>
                <td className="px-4 py-3">
                  <button className="rounded px-3 py-1 text-xs font-medium bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400 dark:hover:bg-red-900/40">
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
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
          disabled={!data || data.length < 10}
          className="rounded px-4 py-2 text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:text-gray-300"
        >
          Next
        </button>
      </div>
    </div>
  )
}

export default AdminOrders