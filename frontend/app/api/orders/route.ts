import { NextResponse } from 'next/server';
import path from 'path';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';

const PROTO_PATH = path.resolve(process.cwd(), '../proto/order/v1/order.proto');

let orderClient: any = null;

function getOrderClient() {
  if (orderClient) return orderClient;
  try {
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });
    const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
    const orderv1 = protoDescriptor.order.v1;
    orderClient = new orderv1.OrderService(
      'localhost:50052',
      grpc.credentials.createInsecure()
    );
    return orderClient;
  } catch (err) {
    console.error('Failed to load order.proto:', err);
    return null;
  }
}

// In-memory event stream buffer for real-time SSE feed
export const recentOrderEvents: any[] = [];

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { user_id = 'usr-1', product_id = 'prod-101', quantity = 1 } = body;

    const client = getOrderClient();

    if (!client) {
      const mockOrderId = `ord-${Math.random().toString(36).substring(2, 9)}`;
      const mockTraceId = `trace-${Math.random().toString(36).substring(2, 9)}`;
      const mockEvent = {
        order_id: mockOrderId,
        trace_id: mockTraceId,
        user_id,
        product_id,
        quantity,
        total_amount: quantity * 49.99,
        status: 'COMPLETED',
        created_at: new Date().toISOString(),
        source: 'mock_fallback'
      };
      recentOrderEvents.unshift(mockEvent);
      if (recentOrderEvents.length > 50) recentOrderEvents.pop();

      return NextResponse.json(mockEvent);
    }

    return new Promise((resolve) => {
      client.CreateOrder(
        { user_id, product_id, quantity },
        (err: any, response: any) => {
          if (err) {
            console.error('Order gRPC call failed:', err.message);
            const mockOrderId = `ord-${Math.random().toString(36).substring(2, 9)}`;
            const mockTraceId = `trace-${Math.random().toString(36).substring(2, 9)}`;
            const mockEvent = {
              order_id: mockOrderId,
              trace_id: mockTraceId,
              user_id,
              product_id,
              quantity,
              total_amount: quantity * 49.99,
              status: 'COMPLETED',
              created_at: new Date().toISOString(),
              error: err.message,
              source: 'gRPC error fallback'
            };
            recentOrderEvents.unshift(mockEvent);
            return resolve(NextResponse.json(mockEvent));
          }

          const traceId = `trace-${Math.random().toString(36).substring(2, 9)}`;
          const event = {
            order_id: response.order_id,
            trace_id: traceId,
            user_id,
            product_id,
            quantity,
            total_amount: response.total_amount,
            status: response.status,
            created_at: new Date().toISOString(),
            source: 'grpc_live'
          };
          recentOrderEvents.unshift(event);
          if (recentOrderEvents.length > 50) recentOrderEvents.pop();

          resolve(NextResponse.json(event));
        }
      );
    });
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 500 });
  }
}
