import { useRef, useState } from "react"
import type { JuiceData, JuiceUpdateParams } from "../types"
import { useUpdateJuice } from "../hooks/useUpdateJuice"
import { useQueryClient } from "@tanstack/react-query"
import { useUploadJuiceImage } from "../hooks/useUploadJuiceImage"
import ConfirmImageModal from "./ConfirmImageModal"

type JuiceRowProps = {
    juice: JuiceData
    onDelete: (id: string) => void
}

function JuiceTableRow({ juice, onDelete }: JuiceRowProps) {

    const [isEditing, setisEditing] = useState(false)
    const [displayPrice, setDisplayPrice] = useState((juice.price / 100).toFixed(2))
    const [pendingFile, setPendingFile] = useState<File | null>(null)
    const [editForm, setEditForm] = useState<JuiceUpdateParams>({
        name: juice.name,
        description: juice.description,
        price: juice.price,
        stock: juice.stock,
    })

    const fileInputRef = useRef<HTMLInputElement>(null)
    const updateJuice = useUpdateJuice()
    const uploadImage = useUploadJuiceImage()
    const queryClient = useQueryClient()
    const handleSave = () => {
        updateJuice.mutate({ id: juice.id, juiceData: editForm }, {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: ['juices'] })
                setisEditing(false)
            }
        })
    }

    const onChange = (field: keyof JuiceUpdateParams, value: string | number) => {
        setEditForm(prev => ({ ...prev, [field]: value }))
    }

    const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0]
        if (!file) return
        setPendingFile(file)
        e.target.value = ""
    }

    const handleConfirm = () => {
        if (!pendingFile) return
        uploadImage.mutate({ id: juice.id, file: pendingFile }, {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: ['juices'] })
                setPendingFile(null)
            }
        })
    }

    return (
        <>
            {pendingFile && (
                <ConfirmImageModal
                    currentImageUrl={juice.image_url}
                    newFile={pendingFile}
                    onConfirm={handleConfirm}
                    onCancel={() => setPendingFile(null)}
                    isPending={uploadImage.isPending}
                />
            )}
            <tr className="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                <td className="px-4 py-3">
                    <div
                        className="relative h-12 w-12 cursor-pointer group"
                        onClick={() => fileInputRef.current?.click()}
                    >
                        <img
                            src={juice.image_url}
                            alt={juice.name}
                            className="h-12 w-12 rounded-lg object-cover"
                        />

                        {!uploadImage.isPending && (
                            <div className="absolute inset-0 rounded-lg bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                                <svg className="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                        d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                                </svg>
                            </div>
                        )}

                        {uploadImage.isPending && (
                            <div className="absolute inset-0 rounded-lg bg-black/50 flex items-center justify-center">
                                <svg className="h-5 w-5 text-white animate-spin" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                                </svg>
                            </div>
                        )}

                        <input
                            ref={fileInputRef}
                            type="file"
                            accept="image/jpeg,image/png"
                            className="hidden"
                            onChange={handleImageChange}
                        />
                    </div>
                </td>
                <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                    {isEditing ?
                        <input
                            className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                            value={editForm.name}
                            onChange={(e) => onChange('name', e.target.value)} />
                        : juice.name}
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
                    ) : (juice.price / 100).toFixed(2)}
                </td>
                <td className="px-4 py-3 text-gray-500 dark:text-gray-400 max-w-xs truncate">
                    {isEditing ? <input
                        className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                        value={editForm.description}
                        onChange={(e) => onChange('description', e.target.value)} />
                        : juice.description}
                </td>
                <td className="px-4 py-3 text-gray-500 dark:text-gray-400 max-w-xs truncate">
                    {isEditing ? <input
                        className="w-full rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 text-sm"
                        value={editForm.stock}
                        onChange={(e) => onChange('stock', e.target.value == '' ? 0 : parseInt(e.target.value))} />
                        : juice.stock}
                </td>
                <td className="px-4 py-3">
                    <div className="flex gap-2">
                        {isEditing ?
                            (
                                <>
                                    <button className="rounded px-3 py-1 text-xs font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
                                        onClick={() => {
                                            setEditForm({
                                                name: juice.name,
                                                description: juice.description,
                                                price: juice.price,
                                                stock: juice.stock,
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
                                            onDelete(juice.id)
                                        }}>
                                        Delete
                                    </button>
                                </>
                            )}
                    </div>
                </td>
            </tr>
        </>
    )
}
export default JuiceTableRow