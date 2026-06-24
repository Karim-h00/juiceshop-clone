import { useState } from "react"
import { useUpdateUserData } from "../hooks/useUpdateUserData"
import { useAuthStore } from "../store/authStore"

function UpdateUserForm() {
  const { user } = useAuthStore()

  const [username, setUsername] = useState(user?.username ?? "")
  const [email, setEmail] = useState(user?.email ?? "")
  const [showSuccess, setShowSuccess] = useState(false)

  const updateUser = useUpdateUserData()

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault()
    updateUser.mutate({ username, email }, {
        onSuccess: () => {
            setShowSuccess(true)
            setTimeout(()=>setShowSuccess(false), 3000)
        }
    })
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white border border-gray-200 rounded-xl p-6 flex flex-col gap-4">
      <h2 className="text-base font-medium text-gray-900">Account info</h2>

      <div className="flex flex-col gap-1.5">
        <label className="text-sm text-gray-500" htmlFor="username">Username</label>
        <input
          id="username"
          type="text"
          value={username}
          onChange={e => setUsername(e.target.value)}
          className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-300"
        />
      </div>

      <div className="flex flex-col gap-1.5">
        <label className="text-sm text-gray-500" htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-300"
        />
      </div>

      {updateUser.isError && (
        <p className="text-sm text-red-500">Failed to update profile. Try again.</p>
      )}
      {showSuccess && (
        <p className="text-sm text-green-600">Profile updated.</p>
      )}

      <div className="flex justify-end pt-1">
        <button
          type="submit"
          disabled={updateUser.isPending}
          className="bg-gray-900 text-white text-sm px-4 py-2 rounded-lg hover:bg-gray-700 disabled:opacity-50 transition-colors"
        >
          {updateUser.isPending ? "Saving…" : "Save changes"}
        </button>
      </div>
    </form>
  )
}

export default UpdateUserForm