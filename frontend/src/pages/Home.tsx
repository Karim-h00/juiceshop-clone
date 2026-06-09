import Hero from "../components/Hero"
import { useJuice } from "../hooks/useJuice"
import { Link } from "react-router-dom"

type juiceData = {
    ID: string,
    ImageUrl: string,
    created_at: string,
    Name: string,
    Description: string,
    Price: number
}

function Home() {
    const { data, isLoading, isError, error } = useJuice()

    if (isLoading) {
        console.log("loading")
        return <div>loading</div>
    }
    if (isError) {
        console.log("err", error)
        return <div>Something went wrong</div>
    }

    console.log(data)
    return (
        <>
        <Hero />
        <div className="container mx-auto px-4 py-6">
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {data.map((item: juiceData) => (
                    <div className="w-full max-w-xs rounded-xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-gray-700 dark:bg-gray-800" key={item.ID}>
                        <img
                            src={item.ImageUrl}
                            alt={item.Name}
                            className="h-48 w-full rounded-t-xl object-cover"
                        />

                        <div className="p-4 space-y-2">
                            <h3 className="text-lg font-semibold text-emerald-700 dark:text-emerald-400">{item.Name}</h3>
                            <p className="text-sm text-gray-600 line-clamp-2 dark:text-gray-300">
                                {item.Description}
                            </p>
                            <p className="text-xl font-bold text-gray-900 dark:text-white">${(item.Price / 100).toFixed(2)}</p>

                            <Link
                                to={`/juice/${item.ID}`}
                                className="mt-2 block w-full rounded bg-emerald-600 py-2 text-center text-sm font-medium text-white hover:bg-emerald-700 dark:bg-emerald-500 dark:hover:bg-emerald-600"
                            >
                                View Details
                            </Link>
                        </div>
                    </div>
                ))}
            </div>
        </div>
        </>
    )

}

export default Home