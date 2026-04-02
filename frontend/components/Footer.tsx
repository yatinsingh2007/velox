'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navigation = {
  product: [
    { name: 'Velox Engine', href: '#' },
    { name: 'Documentation', href: '/docs' },
    { name: 'API Reference', href: '/docs#api-reference' },
    { name: 'System Status', href: '#' },
  ],
  company: [
    { name: 'About Us', href: '#' },
    { name: 'Contact Us', href: 'mailto:contact@velox.dev' },
    { name: 'Blog', href: '#' },
  ],
  legal: [
    { name: 'Privacy Policy', href: '#' },
    { name: 'Terms of Service', href: '#' },
  ],
};

export default function Footer() {
  const pathname = usePathname();
  
  // Hide footer on pure interface pages
  if (pathname === '/login' || pathname === '/signup' || pathname === '/editor') {
    return null;
  }

  return (
    <footer className="bg-background mt-auto border-t border-white/5" aria-labelledby="footer-heading">
      <h2 id="footer-heading" className="sr-only">
        Footer
      </h2>
      <div className="mx-auto max-w-7xl px-6 pb-8 pt-16 sm:pt-24 lg:px-8 lg:pt-32">
        <div className="xl:grid xl:grid-cols-3 xl:gap-8">
          <div className="space-y-8">
            <Link href="/" className="flex items-center gap-2 group">
              <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center font-bold text-primary tracking-tighter border border-primary/20 shadow-[0_0_15px_rgba(255,90,0,0.1)] transition-all">
                <svg className="w-5 h-5 text-primary transform -skew-x-12" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M4.5 4L10.5 19h3.5L8 4H4.5z" />
                  <path d="M12.5 4L18.5 19h3.5L16 4h-3.5z" className="opacity-40" />
                </svg>
              </div>
              <span className="font-extrabold text-xl tracking-wide text-white ml-2">
                Velox
              </span>
            </Link>
            <p className="text-sm leading-6 text-foreground/50 max-w-xs font-mono">
              The high-performance remote code execution engine built for scale, fairness, and absolute velocity.
            </p>
          </div>
          <div className="mt-16 grid grid-cols-2 gap-8 xl:col-span-2 xl:mt-0">
            <div className="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 className="text-sm font-bold leading-6 text-white uppercase tracking-wider">Product</h3>
                <ul role="list" className="mt-6 space-y-4">
                  {navigation.product.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm leading-6 text-white/50 hover:text-primary transition-colors">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
              <div className="mt-10 md:mt-0">
                <h3 className="text-sm font-bold leading-6 text-white uppercase tracking-wider">Company</h3>
                <ul role="list" className="mt-6 space-y-4">
                  {navigation.company.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm leading-6 text-white/50 hover:text-primary transition-colors">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
            <div className="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 className="text-sm font-bold leading-6 text-white uppercase tracking-wider">Legal</h3>
                <ul role="list" className="mt-6 space-y-4">
                  {navigation.legal.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm leading-6 text-white/50 hover:text-primary transition-colors">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div className="mt-16 border-t border-white/5 pt-8 sm:mt-20 lg:mt-24 flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-xs leading-5 text-white/40">
            &copy; {new Date().getFullYear()} Velox Code Judger. All rights reserved.
          </p>
          {/* <div className="flex items-center gap-2 text-xs font-mono">
            <div className="w-2 h-2 rounded-full bg-success animate-pulse shadow-[0_0_10px_rgba(34,197,94,0.5)]"></div>
            <span className="text-success">All systems operational</span>
          </div> */}
        </div>
      </div>
    </footer>
  );
}
