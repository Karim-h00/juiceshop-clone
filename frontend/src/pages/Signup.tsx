import { useState } from "react"
import { useNavigate } from 'react-router-dom';
import { useSignup } from "../hooks/useSignup"
import { useCheckPassword } from "../hooks/useCheckPassword";

function Signup() {

    const [credentials, setCredentials] = useState({
        email: '',
        username: '',
        password: ''
    })
    const [clientError, setClientError] = useState<string | null>(null)

    const handleSignup = useSignup()
    const checkPassword = useCheckPassword()
    const navigate = useNavigate()

    const handlePasswordBlur = () => {
        if (credentials.password.length < 8) return
        checkPassword.mutate(credentials.password)
    }

    const validate = () => {
        const usernameRe = /^[a-zA-Z0-9_-]{3,32}$/
        if (!usernameRe.test(credentials.username)) {
            return "Username must be 3-32 characters, letters/numbers/_ only"
        }
        if (credentials.password.length < 8 || credentials.password.length > 128) {
            return "Password must be 8-128 characters"
        }
        return null
    }

    const handleSubmit = (e: React.SubmitEvent) => {
        e.preventDefault();
        const error = validate()
        if (error) {
            setClientError(error)
            return
        }
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
                                <label htmlFor="username" className="block text-sm font-medium dark:text-gray-300">
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
                                    onBlur={handlePasswordBlur}
                                    onChange={(e) =>
                                        setCredentials({ ...credentials, password: e.target.value })
                                    }
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    placeholder="*********"
                                    aria-invalid={handleSignup.isError}
                                    aria-describedby="signup-error"
                                />
                                {checkPassword.data?.pwned && (
                                    <p className="text-yellow-500 text-sm mt-1">
                                        ⚠️ This password has appeared in known data breaches. Consider using a different one.
                                    </p>
                                )}
                            </div>

                            <button
                                type="submit"
                                disabled={handleSignup.isPending}
                                className="w-full flex justify-center py-2 px-4 mt-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {handleSignup.isPending ? "signing up…" : "Sign up"}
                            </button>
                            {(clientError || handleSignup.isError) && (
                                <p id="signup-error" className="text-red-500 text-sm mt-2">
                                    {clientError || "Signup failed, please try again"}
                                </p>
                            )}
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

export default Signup