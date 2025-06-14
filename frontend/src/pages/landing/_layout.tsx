import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
} from "@/components/ui/navigation-menu";
import { Turtle } from "lucide-react";
import { Link, Outlet } from "react-router";
import coolkids from "@/assets/coolkids.svg";

function LandingLayout() {
  return (
    <>
      <div className="h-24 flex justify-between items-center p-4 absolute inset-0">
        <Link to="/">
          <Turtle size={32} />
        </Link>
        <NavigationMenu>
          <NavigationMenuList>
            <NavigationMenuItem>
              <NavigationMenuLink asChild>
                <Link to="/sign">Sign in</Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
      </div>

      <section
        className="h-screen pt-24 bg-cover bg-center"
        style={{ backgroundImage: `url(${coolkids})` }}
      >
        <Outlet />
      </section>
    </>
  );
}

export default LandingLayout;
