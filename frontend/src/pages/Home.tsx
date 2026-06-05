import { useJuice } from "../hooks/useJuice"

function Home() {
    const { data, isLoading, isError , error} = useJuice()

    if(isLoading){
        console.log("loading")
        return <div>loading</div>
    }
    if(isError){
        console.log("err", error)
        return <div>Something went wrong</div>
    }

    console.log(data)
    return <div className="text-white">{JSON.stringify(data)}</div>

}

export default Home