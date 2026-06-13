import { useGetOrderHistory } from "../hooks/useGetOrderHistory"
import { type Order } from "../types"

function OrderHistory() {

    const { data, isLoading, isError } = useGetOrderHistory()

    if (isLoading) return <div>Loading...</div>
    if (isError) return <div>Something went wrong</div>
    if (!data) return null

    return (
        <div className="container mx-auto px-4 py-10 max-w-3xl">
            <h1 className="text-2xl font-bold text-emerald-700 dark:text-emerald-400 mb-6">Order History</h1>

            {data.length === 0 && (
                <p className="text-gray-500 dark:text-gray-400">You haven't placed any orders yet.</p>
            )}

            <div className="space-y-4">
                {data.map((order: Order) => (
                    <div
                        key={order.order_id}
                        className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-800"
                    >
                        <div className="flex justify-between items-center mb-3">
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                {new Date(order.created_at).toLocaleDateString()}
                            </p>
                            <p className="font-bold text-gray-900 dark:text-white">
                                ${(order.total / 100).toFixed(2)}
                            </p>
                        </div>

                        <ul className="space-y-1">
                            {order.items.map((item, idx) => (
                                <li key={idx} className="flex justify-between text-sm text-gray-700 dark:text-gray-300">
                                    <span>{item.name}</span>
                                    <span className="text-gray-500">x{item.quantity}</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                ))}
            </div>
        </div>
    )
}
export default OrderHistory