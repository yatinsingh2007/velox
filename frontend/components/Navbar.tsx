'use client';

import Link from "next/link";
import { Button } from "./ui/Button";

export default function Navbar() {
  return (
    <nav className="w-full border-b border-white/5 bg-background/80 backdrop-blur-xl sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link href="/" className="flex flex-shrink-0 items-center gap-2 group">
              <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center font-bold text-primary tracking-tighter border border-primary/20 shadow-[0_0_15px_rgba(255,90,0,0.1)] group-hover:shadow-[0_0_20px_rgba(255,90,0,0.3)] transition-all cursor-pointer">
                <svg className="w-5 h-5 text-primary transform -skew-x-12 group-hover:translate-x-0.5 group-hover:scale-110 transition-transform duration-300" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M4.5 4L10.5 19h3.5L8 4H4.5z" />
                  <path d="M12.5 4L18.5 19h3.5L16 4h-3.5z" className="opacity-40" />
                </svg>
              </div>
              <span className="font-extrabold text-xl tracking-wide text-white ml-2 drop-shadow-md">
                Velox
              </span>
            </Link>
            <div className="hidden md:ml-10 md:flex md:space-x-8 text-sm font-semibold tracking-wide">
              <a href="/#features" className="text-white/60 hover:text-white transition-colors">
                Features
              </a>
              <Link href="/docs" className="text-white/60 hover:text-white transition-colors">
                Documentation
              </Link>
            </div>
          </div>
          <div className="flex items-center gap-4 text-sm">
            <Button href="/login" variant="ghost" className="hidden sm:inline-flex" size="sm">
              Sign in
            </Button>
            <Button href="/signup" variant="primary" size="sm">
              Get Started
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}
