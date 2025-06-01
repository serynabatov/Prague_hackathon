import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList } from "@/components/ui/navigation-menu";
import { Turtle } from "lucide-react";
import { Link, Outlet } from "react-router";

function EventLayout() {
  return (
    <>
      <div className="h-24 flex justify-between items-center p-4 absolute inset-0">
        <Link to="/events">
          <Turtle size={32} />
        </Link>
        <NavigationMenu>
          <NavigationMenuList>
            <NavigationMenuItem>
              <NavigationMenuLink asChild>
                <Link to="/sign">Profile</Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
      </div>

      <section className="h-screen pt-24">
        <Outlet />
      </section>
    </>
  );
}

export default EventLayout;
