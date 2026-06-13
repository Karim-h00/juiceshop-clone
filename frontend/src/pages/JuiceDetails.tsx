import { useParams } from "react-router-dom"
import { useGetJuiceByName } from "../hooks/useGetJuiceByName"
import { useCartStore } from "../store/cartStore"

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
  if(!data) return null

  return (
    <div className="container mx-auto px-4 py-10 max-w-4xl">
      <div className="flex flex-col md:flex-row gap-10">
        <img
          src={data.ImageUrl}
          alt={data.Name}
          className="w-full md:w-1/2 rounded-xl object-cover max-h-96"
        />
        <div className="flex flex-col justify-center space-y-4">
          <h1 className="text-3xl font-bold text-emerald-700 dark:text-emerald-400">{data.Name}</h1>
          <p className="text-gray-600 dark:text-gray-300">{data.Description}</p>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">${(data.Price / 100).toFixed(2)}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {data.Stock > 0 ? `${data.Stock} in stock` : "Out of stock"}
          </p>
          <button
            disabled={Number(data.Stock) === 0}
            onClick={() => addItem({
              id: data.ID,
              name: data.Name,
              price: data.Price,
              image: data.ImageUrl,
            })}
            className="rounded bg-emerald-600 py-3 text-white font-medium hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-emerald-500 dark:hover:bg-emerald-600"
          >
            Add to cart
          </button>
        </div>
      </div>
    </div>
  )
}

export default JuiceDetails