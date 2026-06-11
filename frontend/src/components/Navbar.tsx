import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { useState, useEffect, useRef } from 'react';
import { useLogout } from '../hooks/useLogout';
import { useCartStore } from '../store/cartStore';

const Navbar = () => {

  const { user, isLoading } = useAuthStore()
  const { items } = useCartStore()
  const [open, setOpen] = useState(false)
  const { mutate: logout } = useLogout();
  const navigate = useNavigate()
  const location = useLocation()
  const dropdownRef = useRef<HTMLDivElement>(null)
  const isHome = location.pathname === '/'

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  useEffect(()=>{
    const handleClickOutside = (e: MouseEvent) => {
      if(dropdownRef.current && !dropdownRef.current.contains(e.target as Node)){
        setOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  },[])

  if(isLoading) return null

  return (
    <nav className={`flex h-16 items-center justify-between px-6 ${isHome
      ? 'bg-green-700'
      : 'bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700'
      }`}>

      <Link to="/" className={`text-xl font-bold ${isHome
        ? 'text-white'
        : 'text-emerald-600'
        } `}>
        juiceshop
      </Link>

      <div className="relative">
        {!user ? (
          <div className="flex items-center space-x-3">
            <Link
              to="/login"
              className={`rounded px-4 py-2 text-sm font-medium ${isHome
                ? 'text-white border border-white/40 hover:bg-white/10'
                : 'text-emerald-700 hover:bg-emerald-50'
                }`}
            >
              Sign In
            </Link>
            <Link
              to="/register"
              className={`rounded px-4 py-2 text-sm font-medium ${isHome
                ? 'bg-white text-green-700 hover:bg-green-50'
                : 'bg-emerald-600 text-white hover:bg-emerald-700'
                }`}
            >
              sign up
            </Link>
          </div>
        ) : (
          <div className="flex items-center space-x-3">

            <Link to="/cart" className={`relative p-2 rounded ${isHome ? 'text-white hover:bg-white/10' : 'text-gray-600 hover:bg-gray-100'
              }`}>
              <span className="text-2xl">🛒</span>
              <span className="absolute top-1 right-1 flex h-4 w-4 items-center justify-center rounded-full bg-emerald-400 text-xs text-white">
                {items.reduce((sum, i) => sum + i.quantity, 0)}
              </span>
            </Link>

            <div className="relative" ref={dropdownRef}>
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
                    onClick={handleLogout}
                    className="block w-full px-4 py-2 text-left text-sm hover:bg-gray-100"
                  >
                    Logout
                  </button>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  )
}

export default Navbar