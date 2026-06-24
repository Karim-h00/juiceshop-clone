import { useState } from "react"
import { useChangePassword } from "../hooks/useChangePassword"

function ChangePassword() {
    const [passwords, setPasswords] = useState({
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
    })
    const [passwordError, setPasswordError] = useState("")

    const changePw = useChangePassword()

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPasswords(prev => ({ ...prev, [e.target.name]: e.target.value }))
    }

    const handleSubmit = (e: React.SubmitEvent) => {
        e.preventDefault()
        setPasswordError("")
        if (passwords.newPassword !== passwords.confirmPassword) {
            setPasswordError("Passwords don't match")
            return
        }

        changePw.mutate(
            { password: passwords.currentPassword, new_password: passwords.newPassword },
            {
                onSuccess: () => {
                    setPasswords({
                        currentPassword: "",
                        newPassword: "",
                        confirmPassword: ""
                    })
                },
            }
        )
    }

    return (
        <form onSubmit={handleSubmit} className="bg-white border border-gray-200 rounded-xl p-6 flex flex-col gap-4">
            <h2 className="text-base font-medium text-gray-900">Change password</h2>

            <div className="flex flex-col gap-1.5">
                <label className="text-sm text-gray-500" htmlFor="currentPassword">Current password</label>
                <input
                    name="currentPassword"
                    type="password"
                    value={passwords.currentPassword}
                    onChange={handleChange}
                    className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-300"
                />
            </div>

            <div className="flex flex-col gap-1.5">
                <label className="text-sm text-gray-500" htmlFor="newPassword">New password</label>
                <input
                    name="newPassword"
                    type="password"
                    value={passwords.newPassword}
                    onChange={handleChange}
                    className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-300"
                />
            </div>

            <div className="flex flex-col gap-1.5">
                <label className="text-sm text-gray-500" htmlFor="confirmPassword">Confirm new password</label>
                <input
                    name="confirmPassword"
                    type="password"
                    value={passwords.confirmPassword}
                    onChange={handleChange}
                    className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-300"
                />
            </div>

            {passwordError && (
                <p className="text-sm text-red-500">{passwordError}</p>
            )}
            {changePw.isError && (
                <p className="text-sm text-red-500">Incorrect password. Try again.</p>
            )}
            {changePw.isSuccess && (
                <p className="text-sm text-green-600">Password changed.</p>
            )}

            <div className="flex justify-end pt-1">
                <button
                    type="submit"
                    disabled={changePw.isPending}
                    className="bg-gray-900 text-white text-sm px-4 py-2 rounded-lg hover:bg-gray-700 disabled:opacity-50 transition-colors"
                >
                    {changePw.isPending ? "Saving…" : "Update password"}
                </button>
            </div>
        </form>
    )
}

export default ChangePassword