import { type User } from "@/features/auth/store";
import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router";

function RootLayout() {
  const navigate = useNavigate();
  const user = localStorage.getItem("user");
  useEffect(() => {
    if (user && (JSON.parse(user) as User)?.sessionToken) {
      navigate("/events");
    } else {
      navigate("/");
    }
  }, [navigate, user]);

  return <Outlet />;
}

export default RootLayout;
