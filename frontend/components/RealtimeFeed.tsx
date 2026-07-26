'use client';

import { useEffect, useState } from 'react';
import { Radio, Database, CheckCircle2, ShieldCheck, Mail, Cpu } from 'lucide-react';

interface OrderEvent {
  order_id: string;
  trace_id: string;
  user_id: string;
  product_id: string;
  quantity: number;
  total_amount: number;
  status: string;
  created_at: string;
  source?: string;
}

export default function RealtimeFeed() {
  const [events, setEvents] = useState<OrderEvent[]>([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const eventSource = new EventSource('/api/events');

    eventSource.onopen = () => {
      setConnected(true);
    };

    eventSource.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data);
        setEvents(data);
      } catch (err) {
        console.error('SSE parse error:', err);
      }
    };

    eventSource.onerror = () => {
      setConnected(false);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="glass-panel rounded-2xl p-6 border border-white/10">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <div className="p-2.5 rounded-xl bg-brand-600/20 text-brand-400 border border-brand-500/30">
            <Radio className="w-5 h-5 animate-pulse" />
          </div>
          <div>
            <h3 className="text-lg font-bold text-white">Live Event & Task Stream</h3>
            <p className="text-xs text-slate-400">Kafka order-events & RabbitMQ invoice_queue</p>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <span
            className={`w-2 h-2 rounded-full ${
              connected ? 'bg-emerald-400 animate-ping' : 'bg-rose-500'
            }`}
          />
          <span className="text-xs font-mono text-slate-400">
            {connected ? 'SSE ACTIVE' : 'RECONNECTING'}
          </span>
        </div>
      </div>

      <div className="space-y-3 max-h-[460px] overflow-y-auto pr-2 custom-scrollbar">
        {events.length === 0 ? (
          <div className="p-8 text-center rounded-xl bg-surface-900/50 border border-white/5 text-slate-500 text-sm">
            No events streamed yet. Place an order from the storefront to see live pipeline execution!
          </div>
        ) : (
          events.map((evt, idx) => (
            <div
              key={evt.order_id + idx}
              className="p-4 rounded-xl bg-surface-900/80 border border-white/5 hover:border-brand-500/30 transition-all text-xs font-mono space-y-2"
            >
              {/* Top Bar */}
              <div className="flex items-center justify-between text-slate-300">
                <div className="flex items-center space-x-2">
                  <span className="px-2 py-0.5 rounded bg-brand-500/20 text-brand-300 font-bold">
                    ORDER PLACED
                  </span>
                  <span className="text-slate-400">{evt.order_id}</span>
                </div>
                <span className="text-slate-500">
                  {new Date(evt.created_at).toLocaleTimeString()}
                </span>
              </div>

              {/* Payload Metrics */}
              <div className="grid grid-cols-3 gap-2 py-1 text-slate-400">
                <div>
                  <span className="text-[10px] text-slate-500 block">PRODUCT</span>
                  <span className="text-slate-200">{evt.product_id} (x{evt.quantity})</span>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 block">TOTAL</span>
                  <span className="text-emerald-400 font-bold">${evt.total_amount.toFixed(2)}</span>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 block">TRACE ID</span>
                  <span className="text-indigo-400">{evt.trace_id || 'trace-autogen'}</span>
                </div>
              </div>

              {/* Pipeline Badges */}
              <div className="flex flex-wrap items-center gap-1.5 pt-1 border-t border-white/5 text-[10px]">
                <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 flex items-center space-x-1">
                  <ShieldCheck className="w-3 h-3" />
                  <span>gRPC Reserve</span>
                </span>
                <span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 border border-blue-500/20 flex items-center space-x-1">
                  <Database className="w-3 h-3" />
                  <span>Postgres Saved</span>
                </span>
                <span className="px-2 py-0.5 rounded bg-purple-500/10 text-purple-400 border border-purple-500/20 flex items-center space-x-1">
                  <Cpu className="w-3 h-3" />
                  <span>Kafka Streamed</span>
                </span>
                <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20 flex items-center space-x-1">
                  <Mail className="w-3 h-3" />
                  <span>RabbitMQ Processed</span>
                </span>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
