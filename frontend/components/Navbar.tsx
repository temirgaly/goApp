'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Activity, ShoppingBag, LayoutDashboard, Database, Server } from 'lucide-react';

export default function Navbar() {
  const pathname = usePathname();

  return (
    <nav className="glass-panel sticky top-0 z-50 px-6 py-4 mb-8">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        {/* Brand */}
        <Link href="/" className="flex items-center space-x-3 group">
          <div className="p-2 rounded-xl bg-gradient-to-tr from-brand-600 to-indigo-500 text-white shadow-lg shadow-indigo-500/25 group-hover:scale-105 transition-transform">
            <Activity className="w-6 h-6 animate-pulse" />
          </div>
          <div>
            <span className="text-xl font-bold bg-gradient-to-r from-white via-slate-200 to-indigo-300 bg-clip-text text-transparent">
              OrderPulse
            </span>
            <span className="block text-[10px] text-indigo-400 font-mono tracking-wider">
              MICROSERVICES PLATFORM
            </span>
          </div>
        </Link>

        {/* Links */}
        <div className="flex items-center space-x-2 bg-surface-800/80 p-1.5 rounded-xl border border-white/5">
          <Link
            href="/"
            className={`flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              pathname === '/'
                ? 'bg-brand-600 text-white shadow-md shadow-brand-600/30'
                : 'text-slate-400 hover:text-white hover:bg-white/5'
            }`}
          >
            <ShoppingBag className="w-4 h-4" />
            <span>Storefront</span>
          </Link>
          <Link
            href="/dashboard"
            className={`flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              pathname === '/dashboard'
                ? 'bg-brand-600 text-white shadow-md shadow-brand-600/30'
                : 'text-slate-400 hover:text-white hover:bg-white/5'
            }`}
          >
            <LayoutDashboard className="w-4 h-4" />
            <span>Admin Dashboard</span>
          </Link>
        </div>

        {/* System Health Indicators */}
        <div className="hidden md:flex items-center space-x-4 text-xs font-mono text-slate-400">
          <div className="flex items-center space-x-1.5 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
            <Server className="w-3.5 h-3.5" />
            <span>gRPC :50051/:50052</span>
          </div>
          <div className="flex items-center space-x-1.5 px-3 py-1.5 rounded-full bg-indigo-500/10 border border-indigo-500/20 text-indigo-400">
            <Database className="w-3.5 h-3.5" />
            <span>Kafka & RabbitMQ</span>
          </div>
        </div>
      </div>
    </nav>
  );
}
