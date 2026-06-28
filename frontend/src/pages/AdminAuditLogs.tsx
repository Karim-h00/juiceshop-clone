import { useGetAuditLogs } from "../hooks/useGetAuditLogs"
import { type AuditLogData } from "../types"


function AdminAuditLogs() {

    const { data, isLoading, isError } = useGetAuditLogs()

    if (isLoading && !data) return <div className="text-gray-500">Loading initial logs...</div>
    if (isError) return <div className="text-red-500">Failed to load logs</div>

    return (
        <div className="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700">
            <table className="w-full text-sm">
                <thead className="bg-gray-50 dark:bg-gray-800 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    <tr>
                        <th className="px-4 py-3">User ID</th>
                        <th className="px-4 py-3">Action</th>
                        <th className="px-4 py-3">Target Type</th>
                        <th className="px-4 py-3">Target ID</th>
                        <th className="px-4 py-3">Target Name</th>
                        <th className="px-4 py-3">Created At</th>
                    </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-900">
                    {data?.map((log: AuditLogData) => (
                        <tr key={log.id} className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                            <td className="px-4 py-3 text-gray-700 dark:text-gray-300 font-mono text-xs">{log.user_id}</td>
                            <td className="px-4 py-3">
                                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300">
                                    {log.action}
                                </span>
                            </td>
                            <td className="px-4 py-3 text-gray-600 dark:text-gray-400">{log.target_type}</td>
                            <td className="px-4 py-3 text-gray-500 dark:text-gray-500 font-mono text-xs">{log.target_id}</td>
                            <td className="px-4 py-3 text-gray-700 dark:text-gray-300">{log.target_name}</td>
                            <td className="px-4 py-3 text-gray-500 dark:text-gray-400 whitespace-nowrap">
                                {new Date(log.created_at).toLocaleString()}
                            </td>
                        </tr>

                    ))
                    }
                </tbody>
            </table>
        </div>
    )
}

export default AdminAuditLogs