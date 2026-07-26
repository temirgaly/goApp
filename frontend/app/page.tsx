'use client';

import { useEffect, useState } from 'react';
import InventoryCard from '@/components/InventoryCard';
import OrderForm from '@/components/OrderForm';
import { ShoppingCart, CheckCircle2, Sparkles, Server, Zap, ShieldCheck } from 'lucide-react';

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  available_quantity: number;
  is_available: boolean;
  source?: string;
}

const CATALOG: Product[] = [
  {
    id: 'prod-101',
    name: 'Wireless Noise-Canceling Headphones',
    description: 'Premium spatial audio with active noise cancellation and 40h battery life.',
    price: 49.99,
    available_quantity: 100,
    is_available: true,
  },
  {
    id: 'prod-102',
    name: 'Mechanical Gaming Keyboard',
    description: 'RGB mechanical switches with hot-swappable keycaps and ultra-low latency.',
    price: 89.99,
    available_quantity: 45,
    is_available: true,
  },
  {
    id: 'prod-103',
    name: 'Ergonomic Wireless Mouse',
    description: 'Precision optical tracking with dual-mode bluetooth and fast USB-C charging.',
    price: 39.99,
    available_quantity: 60,
    is_available: true,
  },
];

export default function StorefrontPage() {
  const [products, setProducts] = useState<Product[]>(CATALOG);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(false);
  const [lastConfirmedOrder, setLastConfirmedOrder] = useState<any | null>(null);

  const fetchLiveStock = async (productId: string) => {
    setLoading(true);
    try {
      const res = await fetch(`/api/inventory?product_id=${productId}`);
      const data = await res.json();
      setProducts((prev) =>
        prev.map((p) =>
          p.id === productId
            ? {
                ...p,
                available_quantity: data.available_quantity,
                price: data.price || p.price,
                is_available: data.is_available,
                source: data.source,
              }
            : p
        )
      );
    } catch (err) {
      console.error('Failed to fetch stock:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLiveStock('prod-101');
  }, []);

  return (
    <div className="space-y-10">
      {/* Hero Section */}
      <div className="relative glass-panel rounded-3xl p-8 md:p-12 overflow-hidden border border-white/10 shadow-2xl">
        <div className="absolute right-0 top-0 w-96 h-96 bg-brand-600/15 rounded-full blur-3xl" />
        <div className="max-w-2xl relative z-10 space-y-4">
          <div className="inline-flex items-center space-x-2 px-3 py-1.5 rounded-full bg-brand-500/10 border border-brand-500/30 text-brand-300 text-xs font-mono">
            <Sparkles className="w-3.5 h-3.5" />
            <span>REAL-TIME gRPC & ASYNC PIPELINE STOREFRONT</span>
          </div>

          <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight text-white leading-tight">
            Next-Gen Microservices <br />
            <span className="bg-gradient-to-r from-indigo-400 via-brand-400 to-purple-400 bg-clip-text text-transparent">
              Order Fulfillment Engine
            </span>
          </h1>

          <p className="text-slate-400 text-sm md:text-base leading-relaxed">
            Every order placed triggers a high-performance gRPC stock reservation on the{' '}
            <strong className="text-white">Inventory Service (:50051)</strong>, persists to{' '}
            <strong className="text-white">PostgreSQL</strong>, streams an event to{' '}
            <strong className="text-white">Kafka</strong>, and queues background invoice jobs into{' '}
            <strong className="text-white">RabbitMQ</strong>.
          </p>

          <div className="pt-2 flex flex-wrap items-center gap-4 text-xs font-mono text-slate-400">
            <div className="flex items-center space-x-1.5">
              <Zap className="w-4 h-4 text-emerald-400" />
              <span>Sub-millisecond gRPC</span>
            </div>
            <div className="flex items-center space-x-1.5">
              <Server className="w-4 h-4 text-indigo-400" />
              <span>Uber Fx & Zap Logs</span>
            </div>
            <div className="flex items-center space-x-1.5">
              <ShieldCheck className="w-4 h-4 text-purple-400" />
              <span>Postgres DB Locks</span>
            </div>
          </div>
        </div>
      </div>

      {/* Confirmed Order Toast Banner */}
      {lastConfirmedOrder && (
        <div className="p-4 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 flex items-center justify-between animate-in fade-in slide-in-from-top duration-300">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-xl bg-emerald-500/20 text-emerald-400">
              <CheckCircle2 className="w-5 h-5" />
            </div>
            <div>
              <span className="font-bold text-sm text-white">
                Order Confirmed: {lastConfirmedOrder.order_id}
              </span>
              <span className="block text-xs font-mono text-emerald-400">
                Trace ID: {lastConfirmedOrder.trace_id || 'trace-live'} | Status: {lastConfirmedOrder.status}
              </span>
            </div>
          </div>
          <button
            onClick={() => setLastConfirmedOrder(null)}
            className="text-xs font-mono text-emerald-400 hover:text-white underline"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Product Catalog Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white">Live Product Catalog</h2>
          <p className="text-xs text-slate-400">Real-time stock indicators synced via gRPC & Redis</p>
        </div>
      </div>

      {/* Product Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {products.map((product) => (
          <InventoryCard
            key={product.id}
            product={product}
            onBuy={(p) => setSelectedProduct(p)}
            loading={loading}
            onRefresh={() => fetchLiveStock(product.id)}
          />
        ))}
      </div>

      {/* Order Modal */}
      {selectedProduct && (
        <OrderForm
          product={selectedProduct}
          onClose={() => setSelectedProduct(null)}
          onOrderSuccess={(orderData) => {
            setLastConfirmedOrder(orderData);
            fetchLiveStock(selectedProduct.id);
          }}
        />
      )}
    </div>
  );
}
