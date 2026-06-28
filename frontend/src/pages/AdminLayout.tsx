import { useEffect } from "react"
import { Link, useLocation, Outlet } from "react-router-dom"
import { useAuthStore } from "../store/authStore"

const navLinks = [
  { to: "/admin/products", label: "Products" },
  { to: "/admin/orders", label: "Orders" },
  { to: "/admin/users", label: "Users" },
  { to: "/admin/audits", label: "Audit Log" },
]

function AdminLayout() {
  const { pathname } = useLocation()

  const initialize = useAuthStore((s)=>s.initialize)
  
  useEffect(()=>{
    initialize()
  },[])

  return (
    <div className="flex min-h-screen bg-gray-50 dark:bg-gray-900">
      <aside className="w-56 shrink-0 border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 flex flex-col">
        <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
          <span className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest">Admin</span>
        </div>
        <nav className="flex-1 px-3 py-4 space-y-1">
          {navLinks.map(({ to, label }) => (
            <Link
              key={to}
              to={to}
              className={`block px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                pathname === to
                  ? "bg-emerald-50 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400"
                  : "text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
              }`}
            >
              {label}
            </Link>
          ))}
        </nav>
        <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
          <Link to="/" className="text-sm text-gray-500 hover:text-emerald-600 dark:text-gray-400">
            ← Back to store
          </Link>
        </div>
      </aside>

      <main className="flex-1 p-8 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  )
}

export default AdminLayout