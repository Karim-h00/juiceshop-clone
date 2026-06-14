import { Link } from "react-router-dom"
import { type JuiceData } from "../types"
import type { cartItem } from "../store/cartStore"

type CardProps = {
    item: JuiceData
    addItem: (item: Omit<cartItem, "quantity">) => void
}
function Card({ item, addItem }: CardProps) {

    return (
        <div className="w-full max-w-xs rounded-xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-gray-700 dark:bg-gray-800" key={item.ID}>
            <Link to={`/juices/${item.Name.toLowerCase().replace(/\s+/g, '-')}`} className="block">
                <img
                    src={item.ImageUrl}
                    alt={item.Name}
                    className="h-48 w-full rounded-t-xl object-cover"
                />

                <div className="px-4 pt-4 space-y-1">
                    <h3 className="text-lg font-semibold text-emerald-700 dark:text-emerald-400">{item.Name}</h3>
                    <p className="text-xl font-bold text-gray-900 dark:text-white">${(item.Price / 100).toFixed(2)}</p>
                </div>
            </Link>

            <div className="px-4 pb-4 pt-2 space-y-2">
                <p className="text-sm text-gray-600 line-clamp-2 dark:text-gray-300">
                    {item.Description}
                </p>
                <button
                    onClick={() => addItem({
                        id: item.ID,
                        name: item.Name,
                        price: item.Price,
                        image: item.ImageUrl
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