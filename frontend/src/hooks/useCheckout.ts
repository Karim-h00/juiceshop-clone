import { useMutation } from "@tanstack/react-query";
import { checkout } from "../api/order";

export const useCheckout = () => {
  return useMutation({
    mutationFn: checkout,
    onSuccess: (data) => {
      console.log("Order placed:", data)
    },
    onError: (error) => {
      console.error("Checkout failed:", error)
    }
  })
}