import { useState } from "react"
import { useGetAdminOrders } from "../hooks/useGetAdminOrders"
import type { AdminOrder } from "../types"
import { useNavigate } from "react-router-dom"

function AdminOrders() {
  const [page, setPage] = useState(1)
  const { data, isLoading, isError } = useGetAdminOrders(page)

  const navigate = useNavigate();

  if (isLoading) return <p className="text-gray-500">Loading...</p>
  if (isError) return <p className="text-red-500">Failed to load orders.</p>
  if (!data) return null

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Orders</h1>
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
            {data.map((order: AdminOrder) => (
              <tr key={order.ID} className="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              onClick={()=>navigate(`/admin/orders/${order.ID}`)}>
                <td className="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs">{order.ID.slice(0, 8)}...</td>
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
          disabled={data.length < 10}
          className="rounded px-4 py-2 text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:text-gray-300"
        >
          Next
        </button>
      </div>
    </div>
  )
}

export default AdminOrders