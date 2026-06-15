import { useState } from "react"
import { useAddJuice } from "../hooks/useAddJuice"
import type { JuiceUpdateParams } from "../types"
import { useQueryClient } from "@tanstack/react-query"

function AddJuiceModal({ onClose }: {onClose: ()=>void}) {

    const [juiceData, setJuiceData] = useState<JuiceUpdateParams>({
        name: '',
        description: '',
        price: 0,
        stock: 0
    }) 

    const handlerAddJuice = useAddJuice()
    const queryClient = useQueryClient()

    const onChange = (field: keyof JuiceUpdateParams, value: string | number) => {
        setJuiceData(prev => ({ ...prev, [field]: value }))
    }

    const handlerSubmit = () =>{
        handlerAddJuice.mutate(juiceData,{
            onSuccess: () =>{
                queryClient.invalidateQueries({queryKey: ['juices']})
                onClose()
            }
        })
    }
    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">

                <h2 className="text-xl font-semibold mb-4">Add Juice</h2>

                <div className="flex flex-col gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                        <input 
                        className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                        onChange={(e)=>onChange('name', e.target.value)} />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea 
                        rows={3} 
                        className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 resize-none"
                        onChange={(e)=>onChange('description', e.target.value)} />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Price (cents)</label>
                        <input type="number" 
                        className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                        onChange={(e)=>onChange('price', parseInt(e.target.value))} />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Stock</label>
                        <input type="number" 
                        className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                        onChange={(e)=>onChange('stock', parseInt(e.target.value))} />
                    </div>
                </div>

                <div className="flex justify-end gap-2 mt-6">
                    <button 
                    className="px-4 py-2 text-sm text-gray-600 hover:text-gray-800"
                    onClick={onClose}>Cancel</button>
                    <button 
                    className="px-4 py-2 text-sm bg-green-600 text-white rounded-md hover:bg-green-700"
                    onClick={handlerSubmit}>Add Juice</button>
                </div>

            </div>
        </div>
    )

}
export default AddJuiceModal