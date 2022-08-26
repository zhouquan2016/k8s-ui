import { Outlet } from "react-router-dom";

export default function Index() {
    return (<div>
        hello
        <Outlet />
    </div>);
}