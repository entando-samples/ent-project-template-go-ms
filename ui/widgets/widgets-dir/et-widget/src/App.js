import {getData} from "./integration/Integration";
import {useState} from "react";
import './App.css';

function App() {
    const [payload, setPayload] = useState("")

    async function callTheApi() {
        const responseObj = await getData();
        if (responseObj["response"]) {
            setPayload(responseObj.response.data.payload)
        }else{
            setPayload(responseObj.error.message)
        }
    }

    return (
        <>
            <div className="btn-div">
                <button onClick={callTheApi}>call the api</button>
            </div>
            <div className="btn-div">
                <span>{payload}</span>
            </div>
        </>
    )
}

export default App
