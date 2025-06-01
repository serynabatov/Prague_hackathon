import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList } from "@/components/ui/navigation-menu";
import { userRepository } from "@/features/auth/store";
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
                <button onClick={() => userRepository.signOut()}>Sign out</button>
              </NavigationMenuLink>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
      </div>

      <section className="min-h-screen pt-24 max-w-7xl mx-auto">
        <Outlet />
      </section>
    </>
  );
}

export default EventLayout;
