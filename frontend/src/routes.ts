import { createBrowserRouter } from "react-router";
import ErrorBoundary from "./components/common/errorBoundary";
import Home from "./pages/landing";
import LandingLayout from "./pages/landing/_layout";
import UserSign from "./pages/auth";
import AuthLayout from "./pages/auth/_layout";
import Events from "./pages/private";
import RootLayout from "./pages/_layout";
import OneTimePassword from "./pages/auth/oneTimePassword";

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
      {
        Component: OneTimePassword,
        path: "/sign/otp-session",
      },
    ],
  },
  {
    Component: RootLayout,
    children: [
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
    ],
  },
]);

export { routes };
