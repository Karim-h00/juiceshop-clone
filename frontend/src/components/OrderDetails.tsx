import { useNavigate, useParams } from "react-router-dom"
import { useGetOrderByID } from "../hooks/useGetOrderByID"

function OrderDetails() {
    const { id } = useParams<{ id: string }>()
    const navigate = useNavigate()
    const { data, isLoading, isError } = useGetOrderByID(id!)

    if (isLoading) return <div>loading...</div>
    if (isError) return <div>Error</div>
    if (!data) return <div>no order found</div>

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 px-4 py-10">
            <div className="max-w-2xl mx-auto">
                <button
                    onClick={() => navigate(-1)}
                    className="text-sm text-gray-500 dark:text-gray-400 hover:text-green-700 dark:hover:text-green-500 mb-6 inline-flex items-center gap-1 cursor-pointer"
                >
                    ← Back to orders
                </button>

                <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
                    <div className="flex justify-between items-start mb-6">
                        <div>
                            <h1 className="text-lg font-semibold text-gray-900 dark:text-white">Order Details</h1>
                            <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">#{data.order_id}</p>
                        </div>
                        <div className="text-right">
                            <p className="text-xs text-gray-400 dark:text-gray-500">Placed on</p>
                            <p className="text-sm text-gray-700 dark:text-gray-300">{new Date(data.created_at).toLocaleDateString()}</p>
                        </div>
                    </div>

                    <table className="w-full text-sm mb-6">
                        <thead>
                            <tr className="text-left text-gray-400 dark:text-gray-500 border-b border-gray-100 dark:border-gray-700">
                                <th className="pb-2 font-medium">Item</th>
                                <th className="pb-2 font-medium text-center">Quantity</th>
                            </tr>
                        </thead>
                        <tbody>
                            {data.items.map((item, i) => (
                                <tr key={i} className="border-b border-gray-50 dark:border-gray-700 last:border-0">
                                    <td className="py-3 text-gray-800 dark:text-gray-200">{item.name}</td>
                                    <td className="py-3 text-center text-gray-500 dark:text-gray-400">{item.quantity}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>

                    <div className="flex justify-between items-center pt-4 border-t border-gray-100 dark:border-gray-700">
                        <span className="text-sm font-medium text-gray-500 dark:text-gray-400">Total</span>
                        <span className="text-lg font-semibold text-gray-900 dark:text-white">${(data.total / 100).toFixed(2)}</span>
                    </div>
                </div>
            </div>
        </div>
    )
}
export default OrderDetails