import { Outlet } from 'react-router-dom'
import Navbar from './components/Navbar'
import './index.css'
import { useEffect } from 'react'
import { getMe, refresh } from './api/auth'
import useAuthStore from './store/authStore'

function App() {

  const { setAuth, isLoading } = useAuthStore()
  const clearAuth = useAuthStore((s) => s.clearAuth)

  const decodeToken = (token: string) => {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
    return payload
  }

  useEffect(() => {
    refresh()
      .then(async (token) => {
        const payload = decodeToken(token)
        const me = await getMe(token)
        setAuth(token, {
          username: me.username,
          email: me.email,
          role: payload.role,
        })
      })
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
