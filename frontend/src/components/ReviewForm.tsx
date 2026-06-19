import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { useAuthStore } from "../store/authStore"
import { useAddReview } from "../hooks/useAddReview"


function ReviewForm({slug}: {slug: string}) {
  const token = useAuthStore(s => s.token)
  const navigate = useNavigate()
  const [hovered, setHovered] = useState(0)
  const [form, setForm] = useState({ rating: 0, comment: "" })
  const { mutate, isPending, isError } = useAddReview(slug)

  const handleSubmit = () => {
    if (!token) { navigate("/login"); return }
    if (form.rating === 0) return
    mutate({ rating: form.rating, comment: form.comment.trim() || undefined })
    setForm({ rating: 0, comment: "" })
  }

  return (
    <div className="mt-6 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <h3 className="text-lg font-semibold text-gray-800 dark:text-white mb-4">Leave a review</h3>

      <div className="flex gap-1 mb-4">
        {[1, 2, 3, 4, 5].map(star => (
          <button
            key={star}
            onClick={() => setForm(f => ({ ...f, rating: star }))}
            onMouseEnter={() => setHovered(star)}
            onMouseLeave={() => setHovered(0)}
            className="text-2xl transition-transform hover:scale-110"
          >
            <span className={(hovered || form.rating) >= star ? "text-amber-400" : "text-gray-300 dark:text-gray-600"}>
              ★
            </span>
          </button>
        ))}
      </div>

      <textarea
        value={form.comment}
        onChange={e => setForm(f => ({ ...f, comment: e.target.value }))}
        placeholder="Share your thoughts (optional)"
        maxLength={500}
        rows={3}
        className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-transparent px-3 py-2 text-sm text-gray-800 dark:text-gray-200 placeholder-gray-400 resize-none focus:outline-none focus:ring-2 focus:ring-emerald-500"
      />

      {isError && <p className="text-sm text-red-500 mt-1">Failed to submit review. Try again.</p>}

      <button
        onClick={handleSubmit}
        disabled={form.rating === 0 || isPending}
        className="mt-3 rounded-lg bg-emerald-600 px-5 py-2 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-emerald-500 dark:hover:bg-emerald-600"
      >
        {isPending ? "Submitting..." : token ? "Submit review" : "Log in to review"}
      </button>
    </div>
  )
}
export default ReviewForm