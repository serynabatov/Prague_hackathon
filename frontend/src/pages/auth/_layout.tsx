import { Text } from "@/components/ui/text";
import { Turtle } from "lucide-react";
import { useMemo } from "react";
import { Link, Outlet, useLocation } from "react-router";
import coolkidsTogether from "@/assets/coolkidToghere.svg";

const redirectingMap = {
  signIn: {
    description: "Don't have an account?",
    action: "Sign up",
    href: "/sign-up",
  },
  signUp: {
    description: "You already have an account?",
    action: "Sign in",
    href: "/sign-in",
  },
};

function AuthLayout() {
  const location = useLocation();

  const discriminateRedirectignMap = useMemo(() => {
    switch (location.pathname) {
      case "/sign-in":
        return redirectingMap.signIn;
      case "/sign-up":
        return redirectingMap.signUp;
      default:
        break;
    }
  }, [location.pathname]);
  return (
    <div className="h-screen bg-zinc-50 flex justify-center items-center">
      <div className="lg:max-w-11/12 lg:max-h-11/12 m-auto lg:rounded-md overflow-hidden min-w-0 flex-1 flex flex-row h-full">
        <div className="bg-white h-full w-full lg:w-1/2 shrink-0 flex flex-col px-2 py-4">
          <Link to="/" className="max-w-fit">
            <Turtle size={32} />
          </Link>
          <Outlet />
          {discriminateRedirectignMap && (
            <Text type="p" className="text-center">
              {discriminateRedirectignMap?.description}{" "}
              <Link
                to={discriminateRedirectignMap?.href}
                className="text-blue-500"
              >
                {discriminateRedirectignMap?.action}
              </Link>
            </Text>
          )}
        </div>
        <div
          className="invisible lg:visible relative w-full h-full"
          style={{ backgroundColor: "#F8F5E9" }}
        >
          <div
            className="absolute inset-0 bg-contain bg-center bg-no-repeat"
            style={{
              backgroundImage: `url(${coolkidsTogether})`,
            }}
          />
        </div>
      </div>
    </div>
  );
}

export default AuthLayout;
