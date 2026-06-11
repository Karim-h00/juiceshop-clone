import CartItem from "../components/CartItem"
import { useCartStore } from "../store/cartStore"
import { useCheckout } from "../hooks/useCheckout"

function Cart() {

    const { items, updateQuantity, removeItem, clearCart } = useCartStore()
    const total = items.reduce((sum, item) => sum + item.price * item.quantity, 0)
    const { mutate: handleCheckout } = useCheckout()

    return (
        <div className="p-4">
            <h2 className="text-2xl">
                Shopping Cart
            </h2>
            {items.length === 0 ? (
                <p>your cart is empty</p>
            ) : (
                <>
                    {items.map((item) => (
                        <CartItem
                            key={item.id}
                            item={item}
                            updateQuantity={updateQuantity}
                            removeItem={removeItem}
                        />
                    ))}

                    {items.length > 0 && (
                        <>
                            <div className="p-4 border-t dark:border-gray-700 space-y-3">
                                <div className="flex justify-between text-sm font-semibold dark:text-white">
                                    <span>Total</span>
                                    <span>${total.toFixed(2)}</span>
                                </div>
                                <button className="w-full bg-green-500 hover:bg-green-600 active:bg-green-700 text-white py-2.5 rounded-lg font-medium transition-colors cursor-pointer"
                                onClick={()=>handleCheckout()}>
                                    Checkout
                                </button>
                                <button
                                    onClick={() => clearCart()}
                                    className="w-full text-sm text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30 py-2 rounded-lg transition-colors cursor-pointer"
                                >
                                    Clear cart
                                </button>
                            </div>
                        </>
                    )}

                </>
            )}
        </div>
    )

}
export default Cart