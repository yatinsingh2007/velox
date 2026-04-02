import React from 'react';
import { Button } from './ui/Button';

export default function Hero() {
  return (
    <div className="relative overflow-hidden bg-background pt-32 pb-40">
      
      {/* Background Gradient Orbs */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-5xl h-full z-0 pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/20 rounded-full blur-3xl opacity-50 mix-blend-screen animate-pulse" style={{ animationDuration: '4s' }}></div>
        <div className="absolute top-1/3 right-1/4 w-[500px] h-[500px] bg-secondary/10 rounded-full blur-[100px] opacity-30 mix-blend-screen"></div>
      </div>
      
      <div className="mx-auto max-w-7xl px-6 lg:px-8 relative z-10">
        <div className="lg:grid lg:grid-cols-12 lg:gap-16 items-center">
          
          {/* Left Text Content */}
          <div className="lg:col-span-6 text-left pb-16 lg:pb-0 relative">
            <div className="mb-10 flex justify-start">
              {/* <div className="relative rounded-full px-5 py-2 text-sm leading-6 text-foreground font-semibold tracking-wide border border-white/10 bg-surface/50 backdrop-blur-md shadow-[0_0_20px_rgba(255,90,0,0.15)] flex items-center gap-3">
                <span className="relative flex h-2.5 w-2.5">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-primary"></span>
                </span>
                <span className="text-white/90">Velox Engine is live</span>
                <span className="text-white/30 ml-2">→</span>
              </div> */}
            </div>
            
            <h1 className="text-5xl font-extrabold tracking-tight text-white sm:text-7xl mb-8 leading-[1.1]">
              The Engine For <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-white via-primary to-primary">Code Judging.</span>
            </h1>
            
            <p className="mt-8 text-xl leading-relaxed text-foreground/70 max-w-lg">
              Stop making your users wait. Velox delivers instant execution, exact millisecond precision, and a totally unfair advantage for your competitive coding platform.
            </p>
            
            <div className="mt-12 flex items-center justify-start gap-x-6">
              <Button href="#get-started" variant="primary" size="lg">
                Deploy Velox
              </Button>
              <Button href="/docs" variant="ghost" size="lg" className="group">
                Read Documentation <span className="ml-2 transform group-hover:translate-x-1 transition-transform">→</span>
              </Button>
            </div>
          </div>

          {/* Right Visual Content */}
          <div className="lg:col-span-6 relative">
            
            {/* Sleek SaaS Graphic */}
            <div className="relative bg-surface/80 backdrop-blur-xl border border-white/10 rounded-2xl p-6 lg:p-8 shadow-[0_20px_60px_-15px_rgba(0,0,0,1)] hover:-translate-y-2 hover:shadow-[0_30px_80px_-20px_rgba(255,90,0,0.2)] transition-all duration-500 overflow-hidden">
               {/* Inner glow line */}
               <div className="absolute top-0 inset-x-0 h-[1px] bg-gradient-to-r from-transparent via-primary/50 to-transparent"></div>

              <div className="flex justify-between items-center border-b border-white/5 pb-5 mb-6">
                <div className="flex gap-2">
                  <div className="w-3 h-3 rounded-full bg-white/10"></div>
                  <div className="w-3 h-3 rounded-full bg-white/10"></div>
                  <div className="w-3 h-3 rounded-full bg-white/10"></div>
                </div>
                <div className="text-xs font-mono text-foreground/40 font-medium">instance_v8.run()</div>
              </div>
              
              <div className="space-y-4">
                {[
                  { id: '# 89a1x', status: 'Accepted', time: '12ms', width: 'w-full', color: 'bg-primary shadow-[0_0_10px_rgba(255,90,0,0.5)]' },
                  { id: '# c72x5', status: 'Compiling', time: '--', width: 'w-2/3', color: 'bg-white/20 animate-pulse' },
                  { id: '# a1b92', status: 'Accepted', time: '8ms', width: 'w-[90%]', color: 'bg-primary shadow-[0_0_10px_rgba(255,90,0,0.5)]' },
                  { id: '# z8d3e', status: 'Accepted', time: '14ms', width: 'w-full', color: 'bg-primary shadow-[0_0_10px_rgba(255,90,0,0.5)]' },
                ].map((row, i) => (
                  <div key={i} className="flex items-center justify-between py-2 px-1">
                    <span className="text-sm font-mono text-white/50 w-20">{row.id}</span>
                    <div className="flex-grow mx-6">
                      <div className="w-full h-1.5 bg-black/50 rounded-full border border-white/5 overflow-visible relative">
                        <div className={`absolute top-0 left-0 h-full rounded-full ${row.color} ${row.width} transition-all duration-1000`}></div>
                      </div>
                    </div>
                    <span className="text-xs font-mono text-white/80 w-12 text-right">{row.time}</span>
                  </div>
                ))}
              </div>
            </div>
            
            {/* Floating decoration element */}
            <div className="absolute -bottom-8 -left-8 bg-surface border border-white/10 p-5 rounded-2xl shadow-[0_20px_40px_-5px_rgba(0,0,0,0.8)] backdrop-blur-md animate-bounce" style={{ animationDuration: '4s' }}>
              <div className="flex items-center gap-4">
                <div className="bg-primary/10 border border-primary/20 rounded-xl p-3 text-primary shadow-[0_0_15px_rgba(255,90,0,0.2)]">
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                </div>
                <div>
                  <p className="text-sm font-bold text-white leading-tight">Zero Cold Starts</p>
                  <p className="text-xs text-white/50 mt-0.5">Always ready context.</p>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
  );
}
