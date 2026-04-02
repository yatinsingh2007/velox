import React from "react";
import { Button } from "./ui/Button";

const features = [
  {
    name: 'Instant Feedback Loops',
    description:
      'Engineered specifically to remove wait times. Users get results as fast as they can type, enhancing the developer experience and ensuring a flow state during contests.',
    icon: (props: React.SVGProps<SVGSVGElement>) => (
      <svg fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" {...props}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
      </svg>
    ),
  },
  {
    name: 'Bulletproof Fairness',
    description:
      'We isolate execution entirely so that memory constraints and runtime metrics are measured precisely identically for every single submission, guaranteeing fair tournament outcomes.',
    icon: (props: React.SVGProps<SVGSVGElement>) => (
      <svg fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" {...props}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
      </svg>
    ),
  },
  {
    name: 'Engineered for Scale',
    description:
      'Traffic spikes during coding competitions are massive. Our auto-scaling architecture dynamically allocates backend resources without a blip in user experience.',
    icon: (props: React.SVGProps<SVGSVGElement>) => (
      <svg fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" {...props}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418" />
      </svg>
    ),
  },
]

export default function Features() {
  return (
    <div className="bg-background py-24 sm:py-32 relative overflow-hidden" id="features">
      
      {/* Background Radial Drop */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[300px] bg-secondary/5 rounded-full blur-[120px] pointer-events-none"></div>

      <div className="mx-auto max-w-7xl px-6 lg:px-8 relative z-10">
        <div className="mx-auto max-w-2xl lg:text-center">
          <h2 className="text-sm font-bold leading-7 text-white/50 tracking-widest uppercase mb-4 border border-white/10 rounded-full inline-block px-4 py-1 bg-white/5 backdrop-blur-sm shadow-xl">Uncompromising Quality</h2>
          <p className="mt-4 text-4xl font-extrabold tracking-tight text-white sm:text-5xl drop-shadow-lg">
            Focus on algorithms.<br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-[#ff8c00]">We handle the execution.</span>
          </p>
          <p className="mt-6 text-xl leading-8 text-foreground/60 max-w-xl mx-auto font-mono">
            Velox shifts the focus away from slow servers and unreadable logs, enabling a fluid platform built entirely for developer productivity.
          </p>
        </div>
        
        <div className="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-none">
          <dl className="grid max-w-xl grid-cols-1 gap-8 lg:max-w-none lg:grid-cols-3">
            {features.map((feature) => (
              <div key={feature.name} className="dev-card p-8 flex flex-col group h-full justify-between !bg-surface/50 backdrop-blur-md">
                <div>
                  <dt className="text-xl font-bold leading-7 text-white tracking-wide flex flex-col gap-6">
                    <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 group-hover:border-primary/50 group-hover:shadow-[0_0_20px_rgba(255,90,0,0.3)] transition-all duration-300">
                      <feature.icon className="h-6 w-6 text-primary group-hover:scale-110 transition-transform duration-300" aria-hidden="true" />
                    </div>
                    {feature.name}
                  </dt>
                  <dd className="mt-4 text-sm leading-relaxed text-white/60 font-mono">{feature.description}</dd>
                </div>
              </div>
            ))}
          </dl>
          
          <div className="mt-20 flex justify-center">
            <Button href="#features" variant="outline" size="lg" className="backdrop-blur-sm">
              View All Capabilities
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
