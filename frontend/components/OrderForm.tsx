'use client';

import { useState } from 'react';
import { ShoppingCart, CheckCircle2, AlertCircle, X, ArrowRight, Server } from 'lucide-react';

interface OrderFormProps {
  product: {
    id: string;
    name: string;
    price: number;
    available_quantity: number;
  } | null;
  onClose: () => void;
  onOrderSuccess: (orderData: any) => void;
}

export default function OrderForm({ product, onClose, onOrderSuccess }: OrderFormProps) {
  const [quantity, setQuantity] = useState(1);
  const [userId, setUserId] = useState('usr-demo-101');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!product) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const res = await fetch('/api/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: userId,
          product_id: product.id,
          quantity,
        }),
      });

      const data = await res.json();
      if (!res.ok || data.error) {
        throw new Error(data.error || 'Failed to place order');
      }

      onOrderSuccess(data);
      onClose();
    } catch (err: any) {
      setError(err.message || 'An error occurred during order processing.');
    } finally {
      setSubmitting(false);
    }
  };

  const total = (product.price * quantity).toFixed(2);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md">
      <div className="glass-panel w-full max-w-md rounded-2xl p-6 relative border border-white/10 shadow-2xl animate-in fade-in zoom-in duration-200">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute right-4 top-4 p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-white/10 transition-colors"
        >
          <X className="w-5 h-5" />
        </button>

        {/* Header */}
        <div className="flex items-center space-x-3 mb-6">
          <div className="p-3 rounded-xl bg-brand-600/20 text-brand-400 border border-brand-500/30">
            <ShoppingCart className="w-6 h-6" />
          </div>
          <div>
            <h3 className="text-xl font-bold text-white">Checkout Order</h3>
            <span className="text-xs text-slate-400">Order Service gRPC :50052</span>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-sm flex items-center space-x-2">
            <AlertCircle className="w-4 h-4 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-mono text-slate-400 mb-1">
              PRODUCT
            </label>
            <input
              type="text"
              readOnly
              value={`${product.name} (${product.id})`}
              className="w-full px-4 py-2.5 rounded-xl bg-surface-900/80 border border-white/10 text-white font-medium text-sm focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-mono text-slate-400 mb-1">
              USER ID
            </label>
            <input
              type="text"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              className="w-full px-4 py-2.5 rounded-xl bg-surface-900/80 border border-white/10 text-white text-sm focus:border-brand-500 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="block text-xs font-mono text-slate-400 mb-1">
              QUANTITY (Max: {product.available_quantity})
            </label>
            <input
              type="number"
              min="1"
              max={product.available_quantity}
              value={quantity}
              onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
              className="w-full px-4 py-2.5 rounded-xl bg-surface-900/80 border border-white/10 text-white text-sm focus:border-brand-500 focus:outline-none"
              required
            />
          </div>

          <div className="p-4 rounded-xl bg-surface-900/60 border border-white/5 flex justify-between items-center my-4">
            <span className="text-slate-400 text-sm">Total Amount</span>
            <span className="text-2xl font-bold text-white">${total}</span>
          </div>

          {/* Service Chain Indicator */}
          <div className="p-3 rounded-xl bg-indigo-500/10 border border-indigo-500/20 text-[11px] font-mono text-indigo-300 flex items-center justify-between">
            <div className="flex items-center space-x-1">
              <Server className="w-3.5 h-3.5" />
              <span>gRPC Order</span>
            </div>
            <ArrowRight className="w-3 h-3 text-slate-500" />
            <span>PostgreSQL</span>
            <ArrowRight className="w-3 h-3 text-slate-500" />
            <span>Kafka + RabbitMQ</span>
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="w-full py-3.5 px-4 rounded-xl font-medium text-sm flex items-center justify-center space-x-2 bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white shadow-lg shadow-brand-600/30 disabled:opacity-50 transition-all"
          >
            {submitting ? (
              <span>Processing gRPC Order...</span>
            ) : (
              <>
                <CheckCircle2 className="w-4 h-4" />
                <span>Confirm Order (${total})</span>
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
}
