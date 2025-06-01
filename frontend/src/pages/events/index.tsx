import bro from "@/assets/bro.png";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import { userAtom } from "@/features/auth/store";
import { useAtom } from "jotai/react";
import { Sparkles } from "lucide-react";
import { useState } from "react";
import QRCode from "react-qr-code";

const nfts = [
  { id: 1, emoji: "🎨", name: "Abstract Art", hash: "#8F42E1" },
  { id: 2, emoji: "🌟", name: "Cosmic Star", hash: "#2B9AE6" },
  { id: 3, emoji: "🦁", name: "Golden Lion", hash: "#F4C724" },
  { id: 4, emoji: "🌈", name: "Rainbow Prism", hash: "#E645A7" },
  { id: 5, emoji: "🔮", name: "Crystal Ball", hash: "#9C42E6" },
  { id: 6, emoji: "🎭", name: "Theater Mask", hash: "#42E6B8" },
  { id: 7, emoji: "⚡", name: "Lightning Bolt", hash: "#E6B842" },
  { id: 8, emoji: "🌺", name: "Rare Flower", hash: "#E64275" },
];

const events = [
  {
    id: 1,
    name: "Tech Conference 2024",
    date: "May 15, 2024",
    location: "Prague, Czech Republic",
    description: "Join us for an exciting tech conference!",
    emoji: "🎤",
  },
  {
    id: 2,
    name: "Art Expo",
    date: "April 20, 2024",
    location: "Paris, France",
    description: "Experience stunning artwork from around the world",
    emoji: "🎨",
  },
  {
    id: 3,
    name: "Music Festival",
    date: "March 10, 2024",
    location: "Berlin, Germany",
    description: "A day filled with amazing live performances",
    emoji: "🎵",
  },
  {
    id: 4,
    name: "Business Summit",
    date: "February 28, 2024",
    location: "London, UK",
    description: "Network with industry leaders",
    emoji: "💼",
  },
];

function Events() {
  const [view, setView] = useState<string>("events");
  const [user] = useAtom(userAtom);

  return (
    <>
      <div className="pb-6 w-full flex justify-center items-center mb-12 bg-linear-to-br from-orange-50 to-bg-zinc-50">
        <div className="flex flex-col items-center">
          <Dialog>
            <DialogTrigger>
              <Avatar className="w-28 h-28 mb-2">
                <AvatarImage src={bro} />
                <AvatarFallback>BRO</AvatarFallback>
              </Avatar>
            </DialogTrigger>
            <DialogContent>
              <QRCode
                style={{ width: "100%", height: "auto" }} // Makes it responsive
                value={user?.name ?? crypto.randomUUID()}
              />
            </DialogContent>
          </Dialog>
          <Text type="h2">Alex Johnson</Text>
          <Text type="p">@alexJohnson</Text>
        </div>
      </div>

      <Tabs onValueChange={setView} value={view}>
        <div className="flex px-4 mb-6">
          {view === "events" && (
            <Button variant="outline">
              <Sparkles />
              Create new event
            </Button>
          )}
          <TabsList className="ml-auto">
            <TabsTrigger value="events">events</TabsTrigger>
            <TabsTrigger value="user">profile</TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="events">
          <div className="flex flex-col gap-4 px-4">
            {events.map((event) => (
              <div className="px-4 py-2 border border-solid shadow rounded-md bg-white flex">
                <div className="flex gap-2">
                  <div className="w-16 h-16 font-[24px] p-2 border-2 border-solid rounded-lg flex justify-center items-center">
                    {event.emoji}
                  </div>
                  <div className="flex flex-col my-auto">
                    <p className="scroll-m-20  text-[18px] md:text-sm lg:text-xl font-semibold tracking-tight">
                      {event.name}
                    </p>
                    <Text type="p" className="hidden lg:block">
                      {event.description}
                    </Text>
                  </div>
                </div>
                <div className="ml-auto mt-auto text-[12px] md:text-base shrink-0">
                  {event.date}
                </div>
              </div>
            ))}
          </div>
        </TabsContent>
        <TabsContent value="user">
          <div className="flex flex-col px-4 mb-6">
            <Text type="h3" className="mb-4">
              My NFTs
            </Text>
            <div className="flex gap-2 overflow-hidden w-full">
              {nfts.map((nft) => (
                <Card key={nft.id}>
                  <CardContent>
                    <div className="w-20 h-20 font-[24px] p-2 border-2 border-solid rounded-lg flex justify-center items-center">
                      {nft.emoji}
                    </div>
                  </CardContent>
                  <div className="text-[12px] text-center">{nft.name} </div>
                </Card>
              ))}
            </div>
          </div>

          <div className="flex flex-col px-4 mb-6">
            <Text type="h3" className="mb-4">
              Hosted events
            </Text>
            <div className="flex flex-col gap-2 overflow-hidden w-full">
              {events.map((event) => (
                <div className="flex-1 px-4 py-2 border border-solid shadow rounded-md bg-white flex flex-col">
                  <div className="flex gap-2">
                    <div className="shrink-0 w-16 h-16 font-[24px] p-2 border-2 border-solid rounded-lg flex justify-center items-center">
                      {event.emoji}
                    </div>
                    <div className="flex flex-col">
                      <p className="scroll-m-20  text-[18px] md:text-sm lg:text-xl font-semibold tracking-tight">
                        {event.name}
                      </p>
                    </div>
                  </div>
                  <Text type="p" className="hidden lg:block text-right">
                    {event.date}
                  </Text>
                </div>
              ))}
            </div>
          </div>

          <div className="flex flex-col px-4">
            <Text type="h3" className="mb-4">
              Visiting events
            </Text>
            <div className="flex flex-col gap-2 overflow-hidden w-full">
              {events.map((event) => (
                <div
                  className="px-4 py-2 border border-solid shadow rounded-md bg-white flex"
                  key={event.id}
                >
                  <div className="flex gap-2">
                    <div className="shrink-0 w-16 h-16 font-[24px] p-2 border-2 border-solid rounded-lg flex justify-center items-center">
                      {event.emoji}
                    </div>
                    <div className="flex flex-col my-auto">
                      <p className="scroll-m-20  text-[18px] md:text-sm lg:text-xl font-semibold tracking-tight">
                        {event.name}
                      </p>
                    </div>
                  </div>
                  <div className="ml-auto mt-auto text-[12px] md:text-base shrink-0">
                    {event.date}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </>
  );
}

export default Events;
