import Navbar from "../components/Navbar";
import Hero from "../components/Hero";
import Comparison from "../components/Comparison";
import Features from "../components/Features";
import Footer from "../components/Footer";

export default function Home() {
  return (
    <>
      <Navbar />
      <main className="flex-grow">
        <Hero />
        <Comparison />
        <Features />
      </main>
    </>
  );
}
