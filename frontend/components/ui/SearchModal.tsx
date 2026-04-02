'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

export function SearchModal() {
  const [isOpen, setIsOpen] = useState(false);
  const [query, setQuery] = useState('');
  const router = useRouter();

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        setIsOpen((prev) => !prev);
      }
      if (e.key === 'Escape') {
        setIsOpen(false);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  if (!isOpen) return null;

  // Mocked Search Results
  const results = [
    { title: 'Velox Engine Architecture', href: '/docs#architecture', category: 'Documentation' },
    { title: 'API Quick Start', href: '/docs#quick-start', category: 'Documentation' },
    { title: 'System Status', href: '/#features', category: 'Landing Page' },
    { title: 'Create an Account', href: '/signup', category: 'Account' },
  ].filter(res => res.title.toLowerCase().includes(query.toLowerCase()));

  return (
    <div className="fixed inset-0 z-[100] flex items-start justify-center pt-[20vh] px-4">
      {/* Backdrop - Visible Website Background */}
      <div 
        className="absolute inset-0 bg-black/30 backdrop-blur-[2px] transition-opacity"
        onClick={() => setIsOpen(false)}
      ></div>

      {/* Modal Wrapper for Punchy Animation */}
      <div 
        className="relative w-full max-w-xl"
        style={{ animation: 'punchy 0.25s cubic-bezier(0.175, 0.885, 0.32, 1.275)' }}
      >
        {/* Tight Pulsing Edge Glow */}
        <div className="absolute -inset-0.5 bg-primary/20 blur-lg rounded-2xl animate-pulse pointer-events-none z-0"></div>

        {/* Modal Content Box */}
        <div className="relative bg-surface/95 backdrop-blur-xl border border-white/10 rounded-2xl shadow-[0_20px_60px_-15px_rgba(0,0,0,0.8)] overflow-hidden z-10 w-full">
          <div className="flex items-center px-4 py-3 border-b border-white/10">
          <svg className="w-5 h-5 text-white/50 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input 
            type="text" 
            placeholder="Search documentation, features, or commands..." 
            className="w-full bg-transparent border-none text-white focus:outline-none focus:ring-0 placeholder:text-white/30"
            autoFocus
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <div className="flex items-center gap-1 mx-2">
            <kbd className="px-2 py-1 text-[10px] font-mono text-white/50 bg-white/5 rounded border border-white/10">ESC</kbd>
          </div>
        </div>

        <div className="max-h-96 overflow-y-auto p-2">
          {results.length > 0 ? (
            <div className="space-y-1">
              <div className="px-3 py-2 text-xs font-bold text-white/40 uppercase tracking-widest">Results</div>
              {results.map((res, i) => (
                <button
                  key={i}
                  className="w-full text-left px-3 py-3 rounded-xl hover:bg-white/5 transition-colors flex flex-col group focus:bg-white/5 focus:outline-none"
                  onClick={() => {
                    router.push(res.href);
                    setIsOpen(false);
                    setQuery('');
                  }}
                >
                  <span className="text-sm font-semibold text-white/90 group-hover:text-primary transition-colors">{res.title}</span>
                  <span className="text-xs text-white/40 mt-1">{res.category}</span>
                </button>
              ))}
            </div>
          ) : (
            <div className="py-12 text-center text-sm text-white/40">
              No results found for "{query}".
            </div>
          )}
        </div>
      </div>
      </div>
      
      <style dangerouslySetInnerHTML={{__html: `
        @keyframes punchy {
          0% { transform: scale(0.95); opacity: 0; }
          100% { transform: scale(1); opacity: 1; }
        }
      `}} />
    </div>
  );
}
