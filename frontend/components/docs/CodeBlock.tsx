'use client';

import React, { useState } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

interface CodeBlockProps {
  code: string;
  language: string;
}

export function CodeBlock({ code, language }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  // Normalize language for the syntax highlighter (e.g. "javascript / node.js" -> "javascript")
  const parsedLang = language.split(' ')[0].toLowerCase() === 'json' ? 'json' : 'javascript';

  const copyToClipboard = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative group rounded-xl overflow-hidden border border-white/10 bg-[#0a0a0a] my-6 shadow-[0_4px_20px_rgba(0,0,0,0.5)]">
      {/* Header bar */}
      <div className="flex justify-between items-center px-4 py-2 bg-white/5 border-b border-white/5">
        <span className="text-xs font-mono text-white/50">{language}</span>
        <button 
          onClick={copyToClipboard}
          className="text-xs font-bold text-white/40 hover:text-white transition-colors flex items-center gap-1 bg-transparent hover:bg-white/10 px-2 py-1 rounded"
        >
          {copied ? (
            <span className="text-success">Copied!</span>
          ) : (
            <>
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              Copy
            </>
          )}
        </button>
      </div>
      
      {/* Syntax Highlighted Content */}
      <div className="overflow-x-auto text-[13px] leading-relaxed">
        <SyntaxHighlighter 
          language={parsedLang} 
          style={vscDarkPlus} 
          customStyle={{
            margin: 0,
            padding: '1rem',
            background: 'transparent',
            fontSize: 'inherited',
            fontFamily: 'var(--font-jetbrains-mono), monospace'
          }}
          wrapLines={true}
        >
          {code}
        </SyntaxHighlighter>
      </div>
    </div>
  );
}
