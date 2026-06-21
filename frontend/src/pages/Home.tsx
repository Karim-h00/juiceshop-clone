import Hero from "../components/Hero"
import Card from "../components/Card"
import { useJuice } from "../hooks/useJuice"
import { useCartStore } from "../store/cartStore"
import { type JuiceData } from "../types"

function Home() {
    const { data, isLoading, isError } = useJuice()
    const { addItem } = useCartStore()

    if (isLoading) {
        return <div>loading</div>
    }
    if (isError) {
        return <div>Something went wrong</div>
    }
    if (!data){
        return <div>No items</div>
    }

    return (
        <>
            <Hero />
            <div className="container mx-auto px-4 py-6">
                <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {data.map((item: JuiceData) => (
                        <Card
                            key={item.id}
                            item={item}
                            addItem={addItem} />
                    ))}
                </div>
            </div>
        </>
    )

}

export default Home