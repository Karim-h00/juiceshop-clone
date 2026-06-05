import { Outlet } from 'react-router-dom'
import Navbar from './components/Navbar'
import './App.css'

function App() {

  return (
    <>
    <div className='min-h-screen flex flex-col'>
      <Navbar />
      <main className='flex-1 bg-white dark:bg-gray-800 container mx-auto px-4 py-6'>
        <Outlet />
      </main>
    </div>
    </>
  )
}

export default App
