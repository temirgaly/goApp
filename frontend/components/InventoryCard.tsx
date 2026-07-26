'use client';

import { Package, ShieldCheck, Zap, RefreshCw } from 'lucide-react';

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  available_quantity: number;
  is_available: boolean;
  source?: string;
}

interface InventoryCardProps {
  product: Product;
  onBuy: (product: Product) => void;
  loading: boolean;
  onRefresh: () => void;
}

export default function InventoryCard({ product, onBuy, loading, onRefresh }: InventoryCardProps) {
  return (
    <div className="glass-card rounded-2xl p-6 relative overflow-hidden group hover:border-brand-500/40 transition-all duration-300 shadow-xl">
      {/* Decorative Gradient Overlay */}
      <div className="absolute -right-16 -top-16 w-32 h-32 bg-brand-500/10 rounded-full blur-2xl group-hover:bg-brand-500/20 transition-all" />

      {/* Header Badge */}
      <div className="flex justify-between items-start mb-4">
        <div className="p-3 rounded-xl bg-surface-700/80 border border-white/10 text-brand-400 group-hover:scale-110 transition-transform">
          <Package className="w-6 h-6" />
        </div>
        <div className="flex items-center space-x-2">
          {product.source === 'grpc_live' ? (
            <span className="px-2.5 py-1 rounded-full text-[10px] font-mono bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 flex items-center space-x-1">
              <Zap className="w-3 h-3" />
              <span>LIVE gRPC</span>
            </span>
          ) : (
            <span className="px-2.5 py-1 rounded-full text-[10px] font-mono bg-amber-500/10 border border-amber-500/30 text-amber-400">
              FALLBACK
            </span>
          )}
          <button
            onClick={onRefresh}
            className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-white/10 transition-colors"
            title="Refresh Stock"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${loading ? 'animate-spin' : ''}`} />
          </button>
        </div>
      </div>

      {/* Title & Desc */}
      <h3 className="text-xl font-bold text-white mb-1 group-hover:text-indigo-300 transition-colors">
        {product.name}
      </h3>
      <p className="text-slate-400 text-sm mb-6 leading-relaxed">
        {product.description}
      </p>

      {/* Stock & Price Stats */}
      <div className="grid grid-cols-2 gap-3 p-3 rounded-xl bg-surface-900/60 border border-white/5 mb-6">
        <div>
          <span className="block text-[11px] font-mono text-slate-400 uppercase tracking-wider">
            Price
          </span>
          <span className="text-xl font-bold text-white">
            ${product.price.toFixed(2)}
          </span>
        </div>
        <div>
          <span className="block text-[11px] font-mono text-slate-400 uppercase tracking-wider">
            Available Stock
          </span>
          <span
            className={`text-xl font-bold ${
              product.available_quantity > 0 ? 'text-emerald-400' : 'text-rose-400'
            }`}
          >
            {product.available_quantity} units
          </span>
        </div>
      </div>

      {/* Action Button */}
      <button
        onClick={() => onBuy(product)}
        disabled={!product.is_available || product.available_quantity <= 0 || loading}
        className="w-full py-3 px-4 rounded-xl font-medium text-sm flex items-center justify-center space-x-2 bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white shadow-lg shadow-brand-600/30 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform active:scale-98"
      >
        <ShieldCheck className="w-4 h-4" />
        <span>Buy Now (gRPC Reserve)</span>
      </button>
    </div>
  );
}
