interface Message {
	id: string;
	payload: any;
}

interface EnqueueRequest {
	payload: any;
}

export interface Env {
	QUEUE_STORAGE: KVNamespace;
}

export default {
	async fetch(request: Request, env: Env): Promise<Response> {
		try {
			const url = new URL(request.url);

			if (request.method === "POST") {
				switch (url.pathname) {
					case "/enqueue": {
						const body = await request.json() as EnqueueRequest;
						const id = crypto.randomUUID();

						const message: Message = {
							id,
							payload: body.payload
						};

						// Get current queue
						const queueStr = await env.QUEUE_STORAGE.get("queue");
						const queue: string[] = queueStr ? JSON.parse(queueStr) : [];

						// Add message ID to queue
						queue.push(id);

						// Store message and updated queue
						await Promise.all([
							env.QUEUE_STORAGE.put(id, JSON.stringify(message)),
							env.QUEUE_STORAGE.put("queue", JSON.stringify(queue))
						]);

						return new Response(JSON.stringify({ id }), {
							headers: { "Content-Type": "application/json" }
						});
					}

					case "/dequeue": {
						// Get current queue
						const queueStr = await env.QUEUE_STORAGE.get("queue");
						const queue: string[] = queueStr ? JSON.parse(queueStr) : [];

						if (queue.length === 0) {
							return new Response(null, { status: 204 }); // No Content
						}

						// Get first message ID and remove it from queue
						const messageId = queue.shift();
						await env.QUEUE_STORAGE.put("queue", JSON.stringify(queue));

						// Get message data
						const messageStr = await env.QUEUE_STORAGE.get(messageId!);
						if (!messageStr) {
							return new Response(JSON.stringify({ error: "Message not found" }), {
								status: 404,
								headers: { "Content-Type": "application/json" }
							});
						}

						// Delete the message and return it
						await env.QUEUE_STORAGE.delete(messageId!);

						return new Response(messageStr, {
							headers: { "Content-Type": "application/json" }
						});
					}
				}
			}

			return new Response("Not Found", { status: 404 });
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred';
			return new Response(JSON.stringify({ error: errorMessage }), {
				status: 500,
				headers: { "Content-Type": "application/json" }
			});
		}
	}
};