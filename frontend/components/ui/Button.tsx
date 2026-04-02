import React from 'react';
import Link from 'next/link';

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost';
type ButtonSize = 'sm' | 'md' | 'lg';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  href?: string;
  className?: string;
  children: React.ReactNode;
}

export function Button({ 
  children, 
  variant = 'primary', 
  size = 'md',
  href, 
  className = '', 
  ...props 
}: ButtonProps) {
  const baseStyles = "inline-flex items-center justify-center font-bold transition-all duration-300 active:scale-95 rounded-xl";
  
  const variants = {
    primary: "bg-primary text-black hover:bg-white hover:text-black border border-transparent shadow-[0_0_15px_rgba(255,90,0,0.3)] hover:shadow-[0_0_30px_rgba(255,255,255,0.5)] hover:-translate-y-0.5",
    secondary: "bg-surface text-foreground hover:bg-surface/80 border border-white/10 shadow-lg hover:-translate-y-0.5",
    outline: "bg-transparent text-foreground border border-white/20 hover:border-primary hover:text-primary shadow-sm hover:shadow-[0_0_15px_rgba(255,90,0,0.2)] hover:-translate-y-0.5",
    ghost: "bg-transparent text-foreground/80 hover:text-foreground hover:bg-white/5 border border-transparent",
  };

  const sizes = {
    sm: "text-xs px-4 py-2",
    md: "text-sm px-6 py-3",
    lg: "text-base px-8 py-4",
  };

  const cssClass = `${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`;

  if (href) {
    // If it's a relative link standard for next, use next/link but for simple anchor tags we use native a
    if (href.startsWith('#')) {
      return (
        <a href={href} className={cssClass}>
          {children}
        </a>
      );
    }
    return (
      <Link href={href} className={cssClass}>
        {children}
      </Link>
    );
  }

  return (
    <button className={cssClass} {...props}>
      {children}
    </button>
  );
}
