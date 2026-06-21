import { Link } from "react-router-dom"
import { type JuiceData } from "../types"
import type { cartItem } from "../store/cartStore"

type CardProps = {
    item: JuiceData
    addItem: (item: Omit<cartItem, "quantity">) => void
}
function Card({ item, addItem }: CardProps) {

    return (
        <div className="w-full max-w-xs rounded-xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-gray-700 dark:bg-gray-800">
            <Link to={`/juices/${item.name.toLowerCase().replace(/\s+/g, '-')}`} className="block">
                <img
                    src={item.image_url}
                    alt={item.name}
                    className="h-48 w-full rounded-t-xl object-cover"
                />

                <div className="px-4 pt-4 space-y-1">
                    <h3 className="text-lg font-semibold text-emerald-700 dark:text-emerald-400">{item.name}</h3>
                    <p className="text-xl font-bold text-gray-900 dark:text-white">${(item.price / 100).toFixed(2)}</p>
                    <div className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
                        <span className="text-yellow-400">★</span>
                        <span>{item.avg_rating > 0 ? item.avg_rating.toFixed(1) : "No ratings yet"}</span>
                        {item.avg_rating > 0 && <span>({item.reviews_count})</span>}
                    </div>
                </div>
            </Link>

            <div className="px-4 pb-4 pt-2 space-y-2">
                <p className="text-sm text-gray-600 line-clamp-2 dark:text-gray-300">
                    {item.description}
                </p>
                <button
                    onClick={() => addItem({
                        id: item.id,
                        name: item.name,
                        price: item.price,
                        image: item.image_url
                    })}
                    className="mt-2 block w-full rounded bg-emerald-600 py-2 text-center text-sm font-medium text-white hover:bg-emerald-700 dark:bg-emerald-500 dark:hover:bg-emerald-600"
                >
                    Add to cart
                </button>
            </div>
        </div>
    )
}

export default Card