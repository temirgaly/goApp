import { NextResponse } from 'next/server';
import path from 'path';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';

const PROTO_PATH = path.resolve(process.cwd(), '../proto/inventory/v1/inventory.proto');

let inventoryClient: any = null;

function getInventoryClient() {
  if (inventoryClient) return inventoryClient;
  try {
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });
    const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
    const inventoryv1 = protoDescriptor.inventory.v1;
    inventoryClient = new inventoryv1.InventoryService(
      'localhost:50051',
      grpc.credentials.createInsecure()
    );
    return inventoryClient;
  } catch (err) {
    console.error('Failed to load inventory.proto:', err);
    return null;
  }
}

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const productId = searchParams.get('product_id') || 'prod-101';

  const client = getInventoryClient();
  if (!client) {
    return NextResponse.json({
      product_id: productId,
      is_available: true,
      available_quantity: 100,
      price: 49.99,
      source: 'mock_fallback'
    });
  }

  return new Promise((resolve) => {
    client.CheckStock(
      { product_id: productId, quantity: 1 },
      (err: any, response: any) => {
        if (err) {
          return resolve(
            NextResponse.json({
              product_id: productId,
              is_available: true,
              available_quantity: 100,
              price: 49.99,
              source: 'fallback_error',
              error: err.message
            })
          );
        }
        resolve(
          NextResponse.json({
            product_id: productId,
            is_available: response.is_available,
            available_quantity: response.available_quantity,
            price: response.price,
            source: 'grpc_live'
          })
        );
      }
    );
  });
}
