import type { Metadata } from 'next';
import './globals.css';
import Navbar from '@/components/Navbar';

export const metadata: Metadata = {
  title: 'OrderPulse — Microservices Platform',
  description: 'High performance gRPC microservice storefront and analytics dashboard.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className="antialiased">
        <Navbar />
        <main className="max-w-7xl mx-auto px-6 pb-16">
          {children}
        </main>
      </body>
    </html>
  );
}
