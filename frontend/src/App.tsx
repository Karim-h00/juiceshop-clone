import { Outlet } from 'react-router-dom'
import Navbar from './components/Navbar'
import './index.css'
import { useEffect } from 'react'
import { useAuthStore } from './store/authStore'

function App() {

  const initialize = useAuthStore((s) => s.initialize)


  useEffect(() => {
    initialize()
  }, [])

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
