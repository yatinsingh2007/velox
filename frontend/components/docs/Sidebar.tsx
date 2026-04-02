'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';

const docLinks = [
  { name: 'Introduction', href: '/docs#introduction', id: 'introduction' },
  { name: 'Architecture', href: '/docs#architecture', id: 'architecture' },
  { name: 'Quick Start', href: '/docs#quick-start', id: 'quick-start' },
];

export function Sidebar() {
  const [activeId, setActiveId] = useState<string>('introduction');

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
          }
        });
      },
      {
        // Trigger heavily towards the top of the viewport
        rootMargin: '-20% 0px -70% 0px'
      }
    );

    // Initial check and observation
    docLinks.forEach((link) => {
      const el = document.getElementById(link.id);
      if (el) observer.observe(el);
    });

    return () => observer.disconnect();
  }, []);

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>, id: string) => {
    e.preventDefault();
    const element = document.getElementById(id);
    if (element) {
      window.scrollTo({
        top: element.offsetTop - 100, // Account for sticky navbar
        behavior: 'smooth'
      });
      // Allow scroll spy to naturally update the state, but we can optimistically set it
      setActiveId(id);
    }
  };

  return (
    <div className="w-64 flex-shrink-0 hidden lg:block sticky top-24 h-[calc(100vh-6rem)] overflow-y-auto border-r border-white/5 bg-background/50 backdrop-blur-3xl px-6 py-4">
      <h3 className="text-xs font-extrabold text-white/40 uppercase tracking-widest mb-6 border-b border-white/5 pb-4">
        Documentation
      </h3>


      <ul className="relative flex flex-col gap-2">
        {/* Animated active background pill for ultra-smooth SaaS look */}
        <div
          className="absolute left-0 w-[calc(100%+16px)] -ml-2 bg-white/5 rounded-lg border border-white/5 transition-transform duration-300 ease-out z-0 h-10 shadow-[0_4px_10px_rgba(0,0,0,0.5)]"
          style={{
            transform: `translateY(${docLinks.findIndex(l => l.id === activeId) * 48}px)`, // 40px height + 8px gap
            opacity: docLinks.findIndex(l => l.id === activeId) === -1 ? 0 : 1
          }}
        ></div>

        {docLinks.map((link) => {
          const isActive = activeId === link.id;

          return (
            <li key={link.name} className="relative z-10 block h-10">
              <Link
                href={link.href}
                onClick={(e) => handleClick(e, link.id)}
                className={`flex items-center w-full h-full px-2 text-sm font-semibold transition-colors duration-200 group ${isActive
                    ? 'text-primary'
                    : 'text-white/50 hover:text-white/80'
                  }`}
              >
                <div className={`w-1.5 h-1.5 rounded-full mr-3 transition-colors duration-300 ${isActive ? 'bg-primary shadow-[0_0_10px_rgba(255,90,0,0.8)]' : 'bg-transparent group-hover:bg-white/20'}`}></div>
                {link.name}
              </Link>
            </li>
          );
        })}
      </ul>

      <div className="mt-12">
        <h3 className="text-xs font-extrabold text-white/40 uppercase tracking-widest mb-6 border-b border-white/5 pb-4">
          Resources
        </h3>
        <ul className="space-y-4 px-3">
          <li>
            <Link href="https://github.com/" className="text-sm font-semibold text-white/50 hover:text-white transition-colors flex items-center gap-2">
              <svg className="w-4 h-4 opacity-50" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
              GitHub Repository
            </Link>
          </li>
        </ul>
      </div>
    </div>
  );
}
