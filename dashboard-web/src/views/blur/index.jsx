import { useRef, useState } from "react";

export default function Blur() {
    const [name, setName] = useState("");
    return (<div>
        <div><input type="text" onBlur={(e) => { 
            if (name === "0") {
                setName("");
            }
         }} value={name} onChange={(e) => setName(e.target.value) }/></div>
        <div><input type="submit" onClick={(e) => { 
            e.preventDefault(); 
            console.log("name=", name) 
        } } /></div>
    </div>)
}