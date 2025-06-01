import { RouterProvider } from "react-router";
import { routes } from "./routes";
import { Provider } from "jotai/react";
import { globalStore } from "./features/store";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./api";
import { Toaster } from "sonner";

function App() {
  return (
    <Provider store={globalStore}>
      <QueryClientProvider client={queryClient}>
        <Toaster position="top-right"/>
        <RouterProvider router={routes} />
      </QueryClientProvider>
    </Provider>
  );
}

export default App;
