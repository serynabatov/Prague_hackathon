import { RouterProvider } from "react-router";
import { routes } from "./routes";
import { Provider } from "jotai/react";
import { globalStore } from "./features/store";

function App() {
  return (
    <Provider store={globalStore}>
      <RouterProvider router={routes} />
    </Provider>
  );
}

export default App;
