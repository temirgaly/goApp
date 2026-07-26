# Step 6: Frontend Application (Next.js) (`todo-step6.md`)

## Objective
Build the **Next.js** frontend application for **OrderPulse**. The web app will allow users to place orders (calling the Order Service API), view active inventory levels, and monitor live streaming order events and task queue updates in real-time.

---

## 📋 Task Checklist

- [x] **1. Initialize Next.js Project (App Router)**
  - Create directory: `frontend/`
  - Bootstrap Next.js with TypeScript, Tailwind CSS, and App Router.
  - Install dependencies: `lucide-react`, `recharts`, `@grpc/grpc-js`, `@grpc/proto-loader`.

- [x] **2. Configure API Routes / BFF (Backend for Frontend)**
  - Build Next.js Route Handlers (`app/api/orders/route.ts`) to bridge HTTP POST requests from the browser to the Go **Order Service** gRPC/HTTP endpoints.
  - Build Route Handlers (`app/api/inventory/route.ts`) to query stock status from the **Inventory Service**.

- [x] **3. Customer Storefront UI (`app/page.tsx`)**
  - Product catalog grid showcasing items with real-time stock indicators.
  - "Buy Now" interactive checkout trigger.
  - Toast notifications confirming order placement with `order_id` and `trace_id`.

- [x] **4. Seller / Admin Dashboard (`app/dashboard/page.tsx`)**
  - **Live Order Feed**: Real-time log/event list of incoming orders.
  - **Queue Metrics**: Visual status of background processing tasks (RabbitMQ invoice generation).
  - **Analytics Chart**: Real-time sales velocity graph populated via event stream updates.

- [x] **5. Real-Time Status Updates (SSE / WebSockets)**
  - Implement a Server-Sent Events (SSE) route handler subscribing to real-time order & task events.
  - Connect Next.js frontend components to auto-refresh UI state when an order is placed and processed.

---

## 📁 Directory Layout Template

```
orderpulse/
├── frontend/
│   ├── app/
│   │   ├── api/
│   │   │   ├── inventory/route.ts
│   │   │   └── orders/route.ts
│   │   ├── dashboard/
│   │   │   └── page.tsx
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── components/
│   │   ├── InventoryCard.tsx
│   │   ├── OrderForm.tsx
│   │   └── RealtimeFeed.tsx
│   ├── lib/
│   │   └── grpc-client.ts
│   ├── package.json
│   └── tailwind.config.js
```

---

## 🛠️ Verification Commands

```bash
# 1. Navigate to frontend directory
cd frontend

# 2. Install dependencies
npm install

# 3. Start Next.js development server
npm run dev

# 4. Open browser at http://localhost:3000
# Test ordering a product and verify live UI updates on the Dashboard
```

---

## Definition of Done (DoD)
- Next.js frontend renders cleanly with modern Tailwind CSS components.
- Submitting an order successfully calls the Go Order Service, reserves inventory, and displays a confirmed `order_id`.
- The Admin Dashboard receives real-time event updates as RabbitMQ background workers process invoices.
