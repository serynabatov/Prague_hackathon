import RotatingText from "@/components/common/rotatingText";
import { seconds } from "@/lib/time";

function Home() {
  return (
    <div
      className="flex items-center justify-center h-full"
      style={{ fontFamily: "Barrio" }}
    >
      <div className="flex gap-2 text-4xl items-center transition-all">
      <RotatingText
        texts={["You was", "We were", "I was"]}
        mainClassName="px-2 sm:px-2 md:px-3 bg-cyan-300 text-black overflow-hidden py-0.5 sm:py-1 md:py-2 justify-center rounded-lg"
        staggerFrom={"last"}
        initial={{ y: "100%" }}
        animate={{ y: 0 }}
        exit={{ y: "-120%" }}
        staggerDuration={0.025}
        splitLevelClassName="overflow-hidden pb-0.5 sm:pb-1 md:pb-1"
        transition={{ type: "spring", damping: 30, stiffness: 400 }}
        rotationInterval={seconds(3)}
      />

      <p>There</p>
      </div>
    </div>
  );
}

export default Home;
