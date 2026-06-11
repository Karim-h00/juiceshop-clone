import type { cartItem } from "../store/cartStore"

type CartItemProps = {
  item: cartItem
  updateQuantity: (id: string, quantity: number) => void
  removeItem: (id: string) => void
}

function CartItem({item, updateQuantity, removeItem}: CartItemProps){
    return(
        <div key={item.id} className="flex items-center justify-between p-4 border rounded  mb-2">
                            <div className="flex items-center">
                                <img src={item.image} alt={item.name} className="w-20 h-20 object-contain" />
                                <div>
                                    <h3 className="text-semibold dark:text-white">{item.name}</h3>
                                    <p className="dark:text-white">{item.price}</p>
                                </div>
                            </div>

                            <div className="flex items-center gap-4">
                                <button
                                    onClick={() => updateQuantity(item.id, item.quantity - 1)}
                                    className="px-3 py-1 bg-gray-300 rounded cursor-pointer"
                                >
                                    -
                                </button>
                                <span className="dark:text-white">{item.quantity}</span>
                                <button
                                    onClick={() => updateQuantity(item.id, item.quantity + 1)}
                                    className="px-3 py-1 bg-gray-300 rounded cursor-pointer"
                                >
                                    +
                                </button>
                                <button
                                    onClick={() => removeItem(item.id)}
                                    className="px-3 py-1 bg-red-500 text-white rounded ml-4 cursor-pointer"
                                >
                                    Remove
                                </button>
                            </div>
                        </div>
    )
}
export default CartItem