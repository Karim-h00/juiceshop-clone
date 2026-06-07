import { Link } from 'react-router-dom'
import useAuthStore from '../store/authStore'
import { useState } from 'react';
import { useLogout } from '../hooks/useLogout';

const Navbar = () => {

    const { token, user } = useAuthStore.getState();
    const [open, setOpen] = useState(false)
    const { mutate: logout } = useLogout();

    return (
         <nav className="flex h-16 items-center justify-between border-b bg-white dark:bg-gray-900 px-4 shadow-sm">

      <Link to="/" className="text-xl font-bold text-emerald-600">
        juiceshop
      </Link>

      <div className="relative">
        {!token ? (
          <div className="flex items-center space-x-3">
            <Link
              to="/login"
              className="rounded px-4 py-2 text-sm font-medium text-emerald-700 hover:bg-emerald-50"
            >
              Sign In
            </Link>
            <Link
              to="/register"
              className="rounded bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700"
            >
              Sign Up
            </Link>
          </div>
        ) : (
          <div className="relative">
            <button
              onClick={() => setOpen(!open)}
              className="flex items-center space-x-2 rounded px-3 py-1 hover:bg-gray-100"
            >
              <span className="text-sm font-medium">{user?.username}</span>
              <div className="h-6 w-6 rounded-full bg-emerald-200" />
            </button>

            {open && (
              <div
                onClick={() => setOpen(false)}
                className="absolute right-0 mt-2 w-40 origin-top-right rounded border bg-white shadow-lg"
              >
                <Link
                  to="/profile"
                  className="block px-4 py-2 text-sm hover:bg-gray-100"
                >
                  Profile
                </Link>
                {user?.role === "admin" && (
                  <Link
                    to="/admin"
                    className="block px-4 py-2 text-sm hover:bg-gray-100"
                  >
                    Admin Panel
                  </Link>
                )}
                <button
                  onClick={()=>logout}
                  className="block w-full px-4 py-2 text-left text-sm hover:bg-gray-100"
                >
                  Logout
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
    )
}

export default Navbar