import { useState } from "react"
import type { JuiceData, JuiceUpdateParams } from "../types"
import { useUpdateJuice } from "../hooks/useUpdateJuice"
import { useQueryClient } from "@tanstack/react-query"

type JuiceRowProps = {
    juice: JuiceData
    onDelete: (id: string) => void
}

function JuiceTableRow({ juice, onDelete }: JuiceRowProps) {

    const [isEditing, setisEditing] = useState(false)
    const [displayPrice, setDisplayPrice] = useState((juice.Price / 100).toFixed(2))
    const [editForm, setEditForm] = useState<JuiceUpdateParams>({
        name: juice.Name,
        description: juice.Description,
        price: juice.Price,
        stock: juice.Stock,
    })

    const updateJuice = useUpdateJuice()
    const queryClient = useQueryClient()
    const handleSave = () => {
        updateJuice.mutate({ id: juice.ID, juiceData: editForm }, {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: ['juices'] })
                setisEditing(false)
            }
        })
    }

    const onChange = (field: keyof JuiceUpdateParams, value: string | number) => {
        setEditForm(prev => ({ ...prev, [field]: value }))
    }

    return (
        <tr className="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
            <td className="px-4 py-3">
                <img
                    src={juice.ImageUrl}
                    alt={juice.Name}
                    className="h-12 w-12 rounded-lg object-cover"
                />
            </td>
            <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                {isEditing ?
                    <input
                        className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                        value={editForm.name}
                        onChange={(e) => onChange('name', e.target.value)} />
                    : juice.Name}
            </td>
            <td className="px-4 py-3 text-gray-700 dark:text-gray-300">
                {isEditing ? (
                    <input
                        className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                        value={displayPrice}
                        onChange={(e) => {
                            const raw = e.target.value
                            setDisplayPrice(raw)
                            const parsed = parseFloat(raw)
                            if (!isNaN(parsed)) {
                                onChange('price', Math.round(parsed * 100))
                            }
                        }}
                    />
                ) : (juice.Price / 100).toFixed(2)}
            </td>
            <td className="px-4 py-3 text-gray-500 dark:text-gray-400 max-w-xs truncate">
                {isEditing ? <input
                    className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                    value={editForm.description}
                    onChange={(e) => onChange('description', e.target.value)} />
                    : juice.Description}
            </td>
            <td className="px-4 py-3 text-gray-500 dark:text-gray-400 max-w-xs truncate">
                {isEditing ? <input
                    className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                    value={editForm.stock}
                    onChange={(e) => onChange('stock', e.target.value == '' ? 0 : parseInt(e.target.value))} />
                    : juice.Stock}
            </td>
            <td className="px-4 py-3">
                <div className="flex gap-2">
                    {isEditing ?
                        (
                            <>
                                <button className="rounded px-3 py-1 text-xs font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
                                    onClick={() => {
                                        setEditForm({
                                            name: juice.Name,
                                            description: juice.Description,
                                            price: juice.Price,
                                            stock: juice.Stock,
                                        })
                                        setisEditing(false)
                                    }}>
                                    cancel
                                </button>
                                <button className="rounded px-3 py-1 text-xs font-medium bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400 dark:hover:bg-red-900/40"
                                    onClick={handleSave}>
                                    save
                                </button>
                            </>
                        ) : (
                            <>
                                <button className="rounded px-3 py-1 text-xs font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
                                    onClick={() => {
                                        setisEditing(true)
                                    }}>
                                    Edit
                                </button>
                                <button className="rounded px-3 py-1 text-xs font-medium bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400 dark:hover:bg-red-900/40"
                                    onClick={() => {
                                        onDelete(juice.ID)
                                    }}>
                                    Delete
                                </button>
                            </>
                        )}

                </div>
            </td>
        </tr>
    )
}
export default JuiceTableRow