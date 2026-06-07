import { useJuice } from "../hooks/useJuice"
import { NavLink } from "react-router-dom"

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
            {data.map((item: juiceData) => (
                <div className="bg-neutral-700 rounded-md w-68 h-auto mb-4 p-5" key={item.ID}>
                    <NavLink to={'#'}>
                        <img className="rounded-md" src={`${item.ImageUrl}`} />
                    </NavLink>
                    <NavLink to="">
                        <h4 className="mt-6 mb-2 text-2xl font-semibold tracking-tight text-heading">
                            {item.Name}
                        </h4>
                    </NavLink>
                    <p className="mb-6 text-body">
                        {item.Description}
                    </p>
                    <div className="flex items-center justify-between">
                        <span className="text-3xl font-extrabold text-heading">${(item.Price/100)}</span>
                    </div>
                    <button type="button" className="items-center text-white bg-blue-600 px-3 py-2 rounded-xl">
                        Add to Cart
                    </button>
                </div>
            ))}
        </>
    )

}

export default Home