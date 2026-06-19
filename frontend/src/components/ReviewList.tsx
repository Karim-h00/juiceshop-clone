import { useAuthStore } from "../store/authStore"
import { useGetReviews } from "../hooks/useGetReviews"
import { useDeleteReview } from "../hooks/useDeleteReview"

export function ReviewList({slug}: {slug: string}) {
  const { data: reviews, isLoading } = useGetReviews(slug)
  const { mutate: deleteReview } = useDeleteReview(slug)
  const { user } = useAuthStore()

  if (isLoading) return <p className="text-sm text-gray-500 mt-6">Loading reviews...</p>
  if (!reviews?.length) return <p className="text-sm text-gray-500 mt-6">No reviews yet. Be the first!</p>

  return (
    <div className="mt-8">
      <h3 className="text-lg font-semibold text-gray-800 dark:text-white mb-4">Reviews</h3>
      <div className="space-y-4">
        {reviews.map((review: any) => (
          <div key={review.id} className="border border-gray-200 dark:border-gray-700 rounded-xl p-4">
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2">
                <span className="font-medium text-sm text-gray-800 dark:text-white">{review.username}</span>
                <span className="text-amber-400 text-sm">{"★".repeat(review.rating)}{"☆".repeat(5 - review.rating)}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="text-xs text-gray-400">
                  {new Date(review.created_at).toLocaleDateString()}
                </span>
                {(user?.username === review.username || user?.role === "admin") && (
                  <button
                    className="text-xs text-red-400 hover:text-red-600"
                    onClick={()=>deleteReview(review.id)}
                  >
                    Delete
                  </button>
                )}
              </div>
            </div>
            {review.comment && (
              <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">{review.comment}</p>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}