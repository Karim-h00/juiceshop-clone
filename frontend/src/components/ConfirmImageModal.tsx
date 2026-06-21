type ConfirmImageModalProps = {
    currentImageUrl: string
    newFile: File
    onConfirm: () => void
    onCancel: () => void
    isPending: boolean
}

function ConfirmImageModal({ currentImageUrl, newFile, onConfirm, onCancel, isPending }: ConfirmImageModalProps) {
    const previewUrl = URL.createObjectURL(newFile)

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-white dark:bg-gray-900 rounded-lg p-6 w-full max-w-md">
                <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-white">Update Image</h2>
                <div className="flex gap-6 justify-center">
                    <div className="flex flex-col items-center gap-2">
                        <p className="text-xs font-medium text-gray-500 uppercase tracking-wider">Current</p>
                        <img src={currentImageUrl} className="h-28 w-28 rounded-lg object-cover" />
                    </div>
                    <div className="flex flex-col items-center gap-2">
                        <p className="text-xs font-medium text-emerald-600 uppercase tracking-wider">New</p>
                        <img src={previewUrl} className="h-28 w-28 rounded-lg object-cover" />
                    </div>
                </div>
                <div className="flex justify-end gap-2 mt-6">
                    <button
                        className="px-4 py-2 text-sm text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
                        onClick={onCancel}
                        disabled={isPending}>
                        Cancel
                    </button>
                    <button
                        className="px-4 py-2 text-sm bg-emerald-600 text-white rounded-md hover:bg-emerald-700 disabled:opacity-50"
                        onClick={onConfirm}
                        disabled={isPending}>
                        {isPending ? "Uploading..." : "Upload"}
                    </button>
                </div>
            </div>
        </div>
    )
}

export default ConfirmImageModal