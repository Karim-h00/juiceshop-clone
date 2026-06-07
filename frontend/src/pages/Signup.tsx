import { useState } from "react"
import { useSignup } from "../hooks/useSignup"
import { useNavigate } from 'react-router-dom';

function Signup() {

    const [credentials, setCredentials] = useState({
        email: '',
        username: '',
        password: ''
    })

    const handleSignup = useSignup()
    const navigate = useNavigate()
    
    const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault();
    handleSignup.mutate(credentials, {
      onSuccess: () => {
        navigate('/login')
      }
    })
  }
    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 shadow-2xl">
            <div className="max-w-md w-full space-y-8 p-8 bg-white dark:bg-gray-800 rounded-2xl shadow-2xl">
                <div>
                    <h2 className="mt-6 text-center text-3xl font-bold dark:text-white">
                        Sign up
                    </h2>

                    <form className="space-y-6" onSubmit={handleSubmit}>
                        <div className="space-6-4">
                            <div>
                                <label htmlFor="email" className="block text-sm font-medium dark:text-gray-300">
                                    Email
                                </label>
                                <input
                                    id="email"
                                    name="email"
                                    type="email"
                                    required
                                    value={credentials.email}
                                    onChange={(e) =>
                                        setCredentials({ ...credentials, email: e.target.value })
                                    }
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    placeholder="you@example.com"
                                    aria-invalid={handleSignup.isError}
                                    aria-describedby="signup-error"
                                />
                            </div>

                            <div>
                                <label htmlFor="email" className="block text-sm font-medium dark:text-gray-300">
                                    username
                                </label>
                                <input
                                    id="username"
                                    name="username"
                                    type="text"
                                    required
                                    value={credentials.username}
                                    onChange={(e) =>
                                        setCredentials({ ...credentials, username: e.target.value })
                                    }
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    placeholder="example"
                                    aria-invalid={handleSignup.isError}
                                    aria-describedby="signup-error"
                                />
                            </div>

                            <div>
                                <label htmlFor="password" className="block text-sm font-medium dark:text-gray-300">
                                    password
                                </label>
                                <input
                                    id="password"
                                    name="password"
                                    type="password"
                                    required
                                    value={credentials.password}
                                    onChange={(e) =>
                                        setCredentials({ ...credentials, password: e.target.value })
                                    }
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    placeholder="*********"
                                    aria-invalid={handleSignup.isError}
                                    aria-describedby="signup-error"
                                />
                            </div>

                            <button
                                type="submit"
                                disabled={handleSignup.isPending}
                                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {handleSignup.isPending ? "signing up…" : "Sign up"}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

export default Signup