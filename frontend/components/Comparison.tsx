export default function Comparison() {
  return (
    <section className="bg-background py-24 relative overflow-hidden" id="compare">
      {/* Background divider glow */}
      <div className="absolute top-0 inset-x-0 h-[1px] bg-gradient-to-r from-transparent via-white/10 to-transparent"></div>
      
      <div className="mx-auto max-w-7xl px-6 lg:px-8 relative z-10">
        <div className="mx-auto max-w-2xl lg:text-center mb-20 text-center">
          <h2 className="text-sm font-bold leading-7 text-primary tracking-widest uppercase mb-3">The Unfair Advantage</h2>
          <p className="text-4xl font-extrabold tracking-tight text-white mb-6">
            Why wait in queue?
          </p>
          <p className="text-lg leading-8 text-foreground/60 max-w-xl mx-auto">
            Traditional judging platforms buckle under tournament traffic, leaving users waiting minutes just to see a compilation error. We fixed that.
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 lg:gap-12 items-center relative">
          
          {/* Subtle connecting lines */}
          <div className="hidden md:block absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-32 h-[1px] bg-gradient-to-r from-white/5 via-primary/30 to-white/5 z-0"></div>

          {/* Traditional Platform Mockup */}
          <div className="dev-card p-8 bg-surface/40 backdrop-blur-sm z-10">
            <div className="flex items-center justify-between border-b border-white/5 pb-4 mb-6">
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 rounded-full bg-red-400"></div>
                <h3 className="text-lg font-bold text-white/40 line-through">Traditional Platform</h3>
              </div>
              <span className="text-xs font-mono text-red-400/80 bg-red-400/10 px-2 py-1 rounded-md">Queue Delay</span>
            </div>
            
            <p className="mb-6 text-sm text-foreground/50">Your code is held up in massive, unscalable databases during contests.</p>
              
            <div className="bg-black/40 rounded-xl p-6 border border-white/5 relative overflow-hidden">
               <div className="absolute inset-0 bg-red-500/5 blur-xl pointer-events-none"></div>
              
              <div className="flex justify-between text-xs mb-3 text-white/50 font-mono">
                <span>Job #89211</span>
                <span>Pending Array</span>
              </div>
              
              {/* Slow loading bar */}
              <div className="w-full h-2 bg-white/5 rounded-full overflow-hidden relative">
                <div className="absolute top-0 left-0 h-full bg-gradient-to-r from-red-600 to-red-400 animate-fill-slow w-full shadow-[0_0_10px_rgba(248,113,113,0.5)]"></div>
              </div>
              
              <div className="mt-6 flex items-center justify-between">
                <div className="flex items-center gap-3 opacity-60">
                  <div className="animate-spin h-4 w-4 border-2 border-red-400 border-t-transparent rounded-full"></div>
                  <span className="text-sm font-mono text-red-200">Executing...</span>
                </div>
                <span className="text-xs font-mono text-red-300">Est: 4m 12s</span>
              </div>
            </div>
          </div>

          {/* Velox Mockup */}
          <div className="dev-card p-8 bg-surface/80 border-primary/30 shadow-[0_0_50px_rgba(255,90,0,0.1)] relative z-20 md:scale-105">
            {/* Inner primary glow */}
            <div className="absolute inset-x-0 top-0 h-1/2 bg-gradient-to-b from-primary/10 to-transparent pointer-events-none rounded-t-2xl"></div>
            
            <div className="flex items-center justify-between border-b border-white/10 pb-4 mb-6 relative">
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 rounded-full bg-primary animate-pulse shadow-[0_0_10px_rgba(255,90,0,0.8)]"></div>
                <h3 className="text-xl font-bold text-white tracking-wide">Velox Engine</h3>
              </div>
              <span className="text-xs font-mono text-primary bg-primary/10 px-3 py-1 rounded-md border border-primary/20 shadow-[0_0_10px_rgba(255,90,0,0.2)]">Instant Execution</span>
            </div>
            
            <p className="mb-6 text-sm text-foreground/80 leading-relaxed">Independent workers scale instantly to handle any traffic spike. Zero delays.</p>
              
            <div className="bg-black/60 rounded-xl p-6 border border-white/10 relative overflow-hidden shadow-inner">
              <div className="flex justify-between text-sm mb-3 text-white font-mono font-bold">
                <span>Job #00001</span>
                <span className="text-primary">Done</span>
              </div>
              
              {/* Fast loading bar */}
              <div className="w-full h-2 bg-white/5 rounded-full overflow-hidden relative border border-white/5">
                <div className="absolute top-0 left-0 h-full bg-gradient-to-r from-primary to-[#ff8c00] animate-fill-fast w-full shadow-[0_0_15px_rgba(255,90,0,0.8)] filter drop-shadow-lg"></div>
              </div>
              
              <div className="mt-6 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="h-5 w-5 bg-success/20 rounded-full flex items-center justify-center border border-success/30 shadow-[0_0_10px_rgba(34,197,94,0.3)]">
                    <svg className="w-3 h-3 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="3">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                  <span className="text-sm font-mono text-success font-semibold tracking-wide">Accepted</span>
                </div>
                <span className="text-sm font-mono text-primary font-bold">12ms</span>
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  )
}
