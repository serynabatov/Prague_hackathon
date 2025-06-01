import { Button } from "@/components/ui/button";
import { Text } from "@/components/ui/text";

function Events() {
  const events = [
    {
      id: 1,
      name: "Tech Conference 2025",
      date: "June 15, 2025",
      location: "Prague, Czech Republic",
      description: "Join the biggest tech conference in Europe!",
    },
    {
      id: 2,
      name: "Music Festival",
      date: "July 10, 2025",
      location: "Berlin, Germany",
      description: "Experience live music like never before!",
    },
    {
      id: 3,
      name: "Startup Meetup",
      date: "August 20, 2025",
      location: "San Francisco, USA",
      description: "Network with top entrepreneurs and investors.",
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-500 to-indigo-600 text-white">
      {/* Hero Section */}
      <header className="text-center py-16 px-4">
        <Text type="h1" className="text-4xl md:text-6xl font-bold mb-4">
          Discover Amazing Events
        </Text>
        <Text type="p" className="text-lg md:text-xl">
          Join events that inspire, connect, and entertain.
        </Text>
      </header>

      {/* Events Section */}
      <section className="px-4 py-8">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {events.map((event) => (
            <div
              key={event.id}
              className="bg-white text-black rounded-lg shadow-lg p-6"
            >
              <Text type="h3" className="text-xl font-semibold mb-2">
                {event.name}
              </Text>
              <Text type="p" className="text-sm text-gray-600 mb-4">
                {event.date} - {event.location}
              </Text>
              <Text type="p" className="mb-4">{event.description}</Text>
              <Button variant="default" className="w-full">
                Join Now
              </Button>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}

export default Events;