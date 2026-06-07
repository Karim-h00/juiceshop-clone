import { useState } from "react"
import { useLogin } from "../hooks/useLogin"
import { type LoginCredentials } from "../types";
import useAuthStore from "../store/authStore";
import { useNavigate } from 'react-router-dom';

function Login() {

  const [credentials, setCredentials] = useState<LoginCredentials>({
    email: '',
    password: ''
  })
  const handleLogin = useLogin();
  const navigate = useNavigate()

  const setAuthStore = useAuthStore((s) => s.setAuthStore)

  const decodeRole = (token: string): string => {
    try {
      const [, payload] = token.split('.');
      return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/'))).role;
    } catch {
      throw new Error('Invalid token payload');
    }
  };

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault();
    handleLogin.mutate(credentials, {
      onSuccess: ({ username, email, token, refreshToken }) => {
        const role = decodeRole(token)
        setAuthStore(token, refreshToken, username, email, role)

        const { token: t, refreshToken: rt, user: u } = useAuthStore.getState();
        console.log('AuthStore after update:', { token: t, refreshToken: rt, user: u });
        navigate('/')
      }
    })
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
      <div className="max-w-md w-full space-y-8 p-8 bg-white dark:bg-gray-800 rounded-2xl shadow-2xl">
        <h2 className="text-center text-3xl font-bold dark:text-white">Sign in</h2>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
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
                aria-invalid={handleLogin.isError}
                aria-describedby="login-error"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium dark:text-gray-300">
                Password
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
              />
            </div>
          </div>

          {handleLogin.isError && (
            <p id="login-error" className="text-red-600 text-sm text-center">
              Login failed. Please check your credentials.
            </p>
          )}

          <button
            type="submit"
            disabled={handleLogin.isPending}
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {handleLogin.isPending ? "Logging in…" : "Sign in"}
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login