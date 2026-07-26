import { NextResponse } from 'next/server';
import { recentOrderEvents } from '../orders/route';

export const dynamic = 'force-dynamic';

export async function GET() {
  const encoder = new TextEncoder();

  const stream = new ReadableStream({
    start(controller) {
      const sendEvents = () => {
        const data = JSON.stringify(recentOrderEvents);
        controller.enqueue(encoder.encode(`data: ${data}\n\n`));
      };

      sendEvents();
      const interval = setInterval(sendEvents, 2000);

      return () => {
        clearInterval(interval);
      };
    },
  });

  return new Response(stream, {
    headers: {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache, no-transform',
      'Connection': 'keep-alive',
    },
  });
}
