import { Link } from 'react-router-dom'
import useAuthStore from '../store/authStore'

const Navbar = () => {
    const { token, user, logout } = useAuthStore()

    return (
        <nav className="bg-white dark:bg-gray-800 shadow px-4 py-3 flex items-center justify-between">
            <Link to="/" className="text-xl font-bold text-green-500">
                JuiceShop
            </Link>

            {/* desktop menu */}
            <div className="hidden md:flex items-center gap-6">
                <Link to="/">Home</Link>

                {!token && (
                    <>
                        <Link to="/login">Login</Link>
                        <Link to="/register">Register</Link>
                    </>
                )}

                {token && (
                    <>
                        <Link to="/orders">My Orders</Link>
                        <button onClick={logout}>Logout</button>
                    </>
                )}

                {token && user?.role === 'admin' && (
                    <Link to="/admin">Admin</Link>
                )}
            </div>
        </nav>
    )
}

export default Navbar