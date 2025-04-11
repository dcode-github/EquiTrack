import { Router, RouterProvider } from "react-router-dom";
import router from "./router";
import { WebSocketProvider } from "./WebSocketContext";
function App() {
  return (
    <>
    <WebSocketProvider>
      <RouterProvider router={router}/>
    </WebSocketProvider>
    </>
  );
}

export default App;
