import { Outlet } from 'react-router-dom'
import Navbar from './components/Navbar'
import './index.css'
import { useEffect } from 'react'
import { getMe } from './api/auth'
import useAuthStore from './store/authStore'

function App() {

  const { setUser, isLoading } = useAuthStore()
  const clearAuth = useAuthStore((s) => s.clearAuth)

  useEffect(() => {
    getMe()
      .then((me) => setUser({ username: me.username, email: me.email, role: me.role }))
      .catch(() => clearAuth())
  }, [])

  if (isLoading) return null

  return (
    <>
      <div className='min-h-screen flex flex-col dark:bg-gray-900'>
        <Navbar />
        <main className='flex-1 bg-white dark:bg-gray-800'>
          <Outlet />
        </main>
      </div>
    </>
  )
}

export default App
