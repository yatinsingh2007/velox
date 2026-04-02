import React from 'react';
import { CodeBlock } from '@/components/docs/CodeBlock';

const quickStartCode = `const response = await fetch('https://api.velox.dev/v1/jobs/submit', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_API_KEY'
  },
  body: JSON.stringify({
    language: 'python:3.11',
    code: 'def solve(arr):\\n  return sum(arr)\\nprint(solve([1,2,3]))',
    stdin: '',
    cpu_time_limit: 1000 // In milliseconds
  })
});

const data = await response.json();
console.log(data);`;

const responseJson = `{
  "id": "job_01h6c4g6vz42d05f32axt",
  "status": "Accepted",
  "time_ms": 12,
  "memory_kb": 2048,
  "stdout": "6\\n",
  "stderr": "",
  "compile_output": null
}`;

export default function DocsPage() {
  return (
    <div className="max-w-4xl text-white/70 leading-relaxed font-sans pb-32">
      
      {/* Introduction Section */}
      <div id="introduction">
        <h1 className="flex items-center gap-4 text-white font-extrabold text-4xl mb-6 font-outfit">
          <div className="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center border border-primary/20 shadow-[0_0_15px_rgba(255,90,0,0.1)]">
            <svg className="w-6 h-6 text-primary" fill="currentColor" viewBox="0 0 24 24">
              <path d="M4.5 4L10.5 19h3.5L8 4H4.5z" />
              <path d="M12.5 4L18.5 19h3.5L16 4h-3.5z" className="opacity-40" />
            </svg>
          </div>
          Velox Engine Documentation
        </h1>
        <p className="mt-6 text-lg">
          Welcome to the official documentation for Velox. The Velox Engine is a distributed, high-performance Remote Code Execution (RCE) environment designed specifically for competitive programming, online IDEs, and judging platforms.
        </p>
      </div>

      <hr className="my-12 border-white/10" />

      {/* Architecture Section */}
      <div id="architecture" className="scroll-mt-24">
        <h2 className="text-white text-3xl font-extrabold font-outfit mb-6">Architecture Overview</h2>
        <p className="mb-6 text-lg">
          Unlike traditional judging platforms that rely on single monolithic databases to queue and process code submissions (which inevitably buckle under contest traffic loads), Velox utilizes a deeply decoupled architecture:
        </p>
        <div className="bg-surface/50 border border-white/10 rounded-2xl p-6 my-6 shadow-[0_4px_20px_rgba(0,0,0,0.5)]">
          <ul className="text-white/70 space-y-4 list-disc list-inside m-0 pl-2">
            <li><strong className="text-white">API Gateway (Go):</strong> Ingests massive concurrency streams and handles routing.</li>
            <li><strong className="text-white">Redis Queue:</strong> Operates strictly as a blazingly fast in-memory task broker ensuring <code className="bg-white/10 px-1.5 py-0.5 rounded font-mono text-sm text-primary">0% packet loss</code>.</li>
            <li><strong className="text-white">Worker Fleet (Docker):</strong> Instantly scales up independent, isolated sub-containers that execute untrusted code securely using Linux namespaces and cgroups.</li>
          </ul>
        </div>
        <p className="mb-6 text-lg">
          Because execution environments are pre-warmed, cold starts are effectively zero. Your users receive feedback in the exact time it takes to compile and run their script—no overhead.
        </p>
      </div>

      <hr className="my-12 border-white/10" />

      {/* Quick Start Section */}
      <div id="quick-start" className="scroll-mt-24">
        <h2 className="text-white text-3xl font-extrabold font-outfit mb-6">Quick Start</h2>
        <p className="mb-6 text-lg">
          Submitting code to Velox requires a single HTTP request. Our engine bypasses traditional complex WebSocket setups for simple polling loops, or you can supply a webhook URL for immediate asynchronous callbacks.
        </p>
        
        <h3 className="text-white/90 mt-8 mb-4">Submitting a Job</h3>
        <p>To run your first code snippet, send a <code className="bg-primary/20 text-primary border border-primary/20 px-1.5 py-0.5 rounded font-mono text-sm">POST</code> request to the job submittal endpoint.</p>
        
        <CodeBlock code={quickStartCode} language="javascript / node.js" />
        
        <div className="bg-primary/5 border border-primary/20 rounded-xl p-4 my-6 flex items-start gap-3">
           <svg className="w-5 h-5 text-primary mt-0.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p className="m-0 text-sm text-white/80 leading-relaxed">
              <strong>Note:</strong> Always ensure you are limiting the <code className="bg-white/10 px-1 py-0.5 rounded">cpu_time_limit</code> to prevent malicious users from triggering infinite loops on your workers. Max allowed value without an Enterprise license is 15000ms.
            </p>
        </div>

        <h3 className="text-white/90 mt-8 mb-4">Understanding the Response</h3>
        <p>Once the worker successfully executes your payload inside an isolated container, the engine will return detailed CPU metrics, memory traces, and output streams.</p>
        
        <CodeBlock code={responseJson} language="json (response)" />
        
        <p className="mt-8">
          The <code className="bg-white/10 px-1 py-0.5 rounded text-primary">status</code> field will return one of the following exact competitive-programming standard verdicts: <br />
           <span className="inline-block mt-3 space-x-2">
            <span className="bg-success/20 text-success border border-success/30 px-2 py-1 rounded text-xs font-mono font-bold">Accepted</span>
            <span className="bg-red-500/20 text-red-400 border border-red-500/30 px-2 py-1 rounded text-xs font-mono font-bold">Wrong Answer</span>
            <span className="bg-secondary/20 text-secondary border border-secondary/30 px-2 py-1 rounded text-xs font-mono font-bold">Time Limit Exceeded</span>
            <span className="bg-zinc-700 text-white/80 border border-zinc-600 px-2 py-1 rounded text-xs font-mono font-bold">Runtime Error</span>
           </span>
        </p>
      </div>

    </div>
  );
}
