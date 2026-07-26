'use client';

import { useEffect, useState } from 'react';
import RealtimeFeed from '@/components/RealtimeFeed';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { LayoutDashboard, TrendingUp, Layers, CheckCircle2, AlertTriangle, ExternalLink } from 'lucide-react';

const INITIAL_CHART_DATA = [
  { time: '10:00', sales: 120, orders: 3 },
  { time: '10:05', sales: 240, orders: 5 },
  { time: '10:10', sales: 180, orders: 4 },
  { time: '10:15', sales: 350, orders: 7 },
  { time: '10:20', sales: 490, orders: 10 },
  { time: '10:25', sales: 620, orders: 13 },
];

export default function DashboardPage() {
  const [chartData, setChartData] = useState(INITIAL_CHART_DATA);
  const [stats, setStats] = useState({
    totalOrders: 13,
    grossRevenue: 620.0,
    gRPCLatency: '1.2ms',
    rabbitmqTasksAcked: '100%',
  });

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div className="flex items-center space-x-2 text-brand-400 font-mono text-xs mb-1">
            <LayoutDashboard className="w-4 h-4" />
            <span>SYSTEM OVERVIEW & ANALYTICS</span>
          </div>
          <h1 className="text-3xl font-extrabold text-white">Admin Operations Dashboard</h1>
        </div>

        <div className="flex items-center space-x-3">
          <a
            href="http://localhost:3000"
            target="_blank"
            rel="noreferrer"
            className="px-4 py-2 rounded-xl glass-card text-xs font-mono text-indigo-300 hover:text-white border border-indigo-500/30 flex items-center space-x-2 transition-all"
          >
            <span>Open Grafana Logs (:3000)</span>
            <ExternalLink className="w-3.5 h-3.5" />
          </a>
          <a
            href="http://localhost:15672"
            target="_blank"
            rel="noreferrer"
            className="px-4 py-2 rounded-xl glass-card text-xs font-mono text-amber-300 hover:text-white border border-amber-500/30 flex items-center space-x-2 transition-all"
          >
            <span>RabbitMQ UI (:15672)</span>
            <ExternalLink className="w-3.5 h-3.5" />
          </a>
        </div>
      </div>

      {/* Metrics Row */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="glass-card rounded-2xl p-5 border border-white/10">
          <div className="flex justify-between items-center text-slate-400 mb-2">
            <span className="text-xs font-mono">GROSS REVENUE</span>
            <TrendingUp className="w-4 h-4 text-emerald-400" />
          </div>
          <span className="text-2xl font-bold text-white">${stats.grossRevenue.toFixed(2)}</span>
          <span className="block text-[11px] text-emerald-400 mt-1 font-mono">
            +18.4% from live orders
          </span>
        </div>

        <div className="glass-card rounded-2xl p-5 border border-white/10">
          <div className="flex justify-between items-center text-slate-400 mb-2">
            <span className="text-xs font-mono">TOTAL ORDERS</span>
            <Layers className="w-4 h-4 text-brand-400" />
          </div>
          <span className="text-2xl font-bold text-white">{stats.totalOrders}</span>
          <span className="block text-[11px] text-indigo-400 mt-1 font-mono">
            Processed via gRPC :50052
          </span>
        </div>

        <div className="glass-card rounded-2xl p-5 border border-white/10">
          <div className="flex justify-between items-center text-slate-400 mb-2">
            <span className="text-xs font-mono">gRPC AVG LATENCY</span>
            <CheckCircle2 className="w-4 h-4 text-purple-400" />
          </div>
          <span className="text-2xl font-bold text-white">{stats.gRPCLatency}</span>
          <span className="block text-[11px] text-purple-400 mt-1 font-mono">
            Direct HTTP/2 binary protocol
          </span>
        </div>

        <div className="glass-card rounded-2xl p-5 border border-white/10">
          <div className="flex justify-between items-center text-slate-400 mb-2">
            <span className="text-xs font-mono">RABBITMQ WORKER ACK</span>
            <CheckCircle2 className="w-4 h-4 text-amber-400" />
          </div>
          <span className="text-2xl font-bold text-white">{stats.rabbitmqTasksAcked}</span>
          <span className="block text-[11px] text-amber-400 mt-1 font-mono">
            Worker invoice_queue active
          </span>
        </div>
      </div>

      {/* Main Grid: Chart & Realtime Feed */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Sales Velocity Chart */}
        <div className="glass-panel rounded-2xl p-6 border border-white/10 space-y-4">
          <div>
            <h3 className="text-lg font-bold text-white">Live Sales Velocity</h3>
            <p className="text-xs text-slate-400">Order revenue over time</p>
          </div>

          <div className="h-72 w-full pt-4">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id="salesGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#6366f1" stopOpacity={0.4} />
                    <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" />
                <XAxis dataKey="time" stroke="#64748b" fontSize={11} />
                <YAxis stroke="#64748b" fontSize={11} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#1f2937',
                    borderColor: 'rgba(255,255,255,0.1)',
                    borderRadius: '0.75rem',
                    color: '#fff',
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="sales"
                  stroke="#6366f1"
                  strokeWidth={2}
                  fillOpacity={1}
                  fill="url(#salesGrad)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Real-time Order & Task Stream */}
        <RealtimeFeed />
      </div>
    </div>
  );
}
