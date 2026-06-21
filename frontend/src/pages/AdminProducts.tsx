import { useQueryClient } from "@tanstack/react-query"
import JuiceTableRow from "../components/JuiceTableRow"
import { useDeleteJuice } from "../hooks/useDeletejuice"
import { useJuice } from "../hooks/useJuice"
import { type JuiceData } from "../types"
import { useState } from "react"
import AddJuiceModal from "../components/AddJuiceModal"

function AdminProducts() {
  const { data, isLoading, isError } = useJuice()
  const [isOpen, setIsOpen] = useState(false)
  const handleDelete = useDeleteJuice()
  const queryClient = useQueryClient()

  const onClose = () =>{
    setIsOpen(false)
  }
  const onDelete = (id: string) =>{
    handleDelete.mutate(id,{
      onSuccess: () => {
        queryClient.invalidateQueries({queryKey:['juices']})
      }
    })
  }

  if (isLoading) return <p className="text-gray-500">Loading...</p>
  if (isError) return <p className="text-red-500">Failed to load products.</p>
  if (!data) return null

  return (
    <>
    {isOpen && <AddJuiceModal onClose={onClose}/>}
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Products</h1>
        <button className="rounded bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700"
        onClick={()=>setIsOpen(true)}>
          + Add Product
        </button>
      </div>

      <div className="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 dark:bg-gray-800 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
            <tr>
              <th className="px-4 py-3">Image</th>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Price</th>
              <th className="px-4 py-3">Description</th>
              <th className="px-4 py-3">Stock</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-900">
            {data.map((juice: JuiceData) => (
              <JuiceTableRow key={juice.id} juice={juice} onDelete={onDelete} />
            ))}
          </tbody>
        </table>
      </div>
    </div>
      </>
  )
}

export default AdminProducts