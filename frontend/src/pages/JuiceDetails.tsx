import { useParams } from "react-router-dom"
import { useGetJuiceByName } from "../hooks/useGetJuiceByName"
import { useCartStore } from "../store/cartStore"
import { ReviewList } from "../components/ReviewList"
import ReviewForm from "../components/ReviewForm"

function JuiceDetails() {
  const { juiceName } = useParams<{ juiceName: string }>()
  const decoded = juiceName!
    .split('-')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
  const { data, isLoading, isError } = useGetJuiceByName(decoded)
  const { addItem } = useCartStore()

  if (isLoading) return <div>Loading...</div>
  if (isError) return <div>Something went wrong</div>
  if (!data) return null

  return (
    <div className="container mx-auto px-4 py-10 max-w-4xl">
      <div className="flex flex-col md:flex-row gap-10">
        <img
          src={data.image_url}
          alt={data.name}
          className="w-full md:w-1/2 rounded-xl object-cover max-h-96"
        />
        <div className="flex flex-col justify-center space-y-4">
          <h1 className="text-3xl font-bold text-emerald-700 dark:text-emerald-400">{data.name}</h1>
          <p className="text-gray-600 dark:text-gray-300">{data.description}</p>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">${(data.price / 100).toFixed(2)}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {data.stock > 0 ? `${data.stock} in stock` : "Out of stock"}
          </p>
          <div className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
            <span className="text-yellow-400">★</span>
            <span>{data.avg_rating > 0 ? data.avg_rating.toFixed(1) : "No ratings yet"}</span>
            {data.avg_rating > 0 && <span>({data.reviews_count})</span>}
          </div>
          <button
            disabled={Number(data.stock) === 0}
            onClick={() => addItem({
              id: data.id,
              name: data.name,
              price: data.price,
              image: data.image_url,
            })}
            className="rounded bg-emerald-600 py-3 text-white font-medium hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-emerald-500 dark:hover:bg-emerald-600"
          >
            Add to cart
          </button>
        </div>
      </div>

      <div className="mt-10">
        <ReviewList slug={juiceName!} />
        <ReviewForm slug={juiceName!} />
      </div>
    </div>
  )
}

export default JuiceDetails