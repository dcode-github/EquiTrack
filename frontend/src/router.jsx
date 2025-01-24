import {
Navigate,
Route,
createBrowserRouter,
createRoutesFromElements,
} from "react-router-dom";
import Auth from "./pages/Auth";
import StockTracker from "./pages/StockTracker";

const router = createBrowserRouter(
createRoutesFromElements(
    <>
    <Route path="/" element={<Auth />}/>
    <Route path="/home" element={<StockTracker />} />
    </>
),
);

export default router;