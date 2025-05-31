import { createBrowserRouter } from "react-router";
import ErrorBoundary from "./components/common/errorBoundary";
import Home from "./pages";
import LandingLayout from "./pages/_layout";
import UserSign from "./pages/auth";
import AuthLayout from "./pages/auth/_layout";
import Events from "./pages/private";

const routes = createBrowserRouter([
  {
    Component: LandingLayout,
    ErrorBoundary,
    children: [
      {
        path: "/",
        Component: Home,
      },
    ],
  },
  {
    ErrorBoundary,
    Component: AuthLayout,
    children: [
      {
        Component: UserSign,
        path: "/sign",
      },
     
    ],
  },
  {
    // add evant layout
    ErrorBoundary,
    children: [
      {
        Component: Events,
        path: "/events",
      },
     
    ],
  },
]);

export { routes };
