export const WEBSOCKET_URL: string =
  process.env.PREVIEW_WEBSOCKET_URL ?? "ws://localhost:1234/websocket";

let events: EventTarget = new EventTarget();
const RETRY_INTERVAL: number = 5000;
let ws: WebSocket | null = null;

let initMessage: InitMessage | null = null;
export function getInitMessage(): InitMessage | null {
  return initMessage;
}

export type InitMessage = {
  kind: "init",
  content: {
    chunkSize: number,
    // TODO: describe
    costs: unknown,
    // TODO: describe
    setup: unknown,
  }
}

export type Message = InitMessage

export class MessageEvent extends Event {
  public readonly message: Message;

  constructor(message: Message) {
    super("message");
    this.message = message;
  }
}

export class ConnectionEvent extends Event {
  public readonly isConnected: boolean;
  public readonly message: string;

  constructor(isConnected: boolean, message: string) {
    super("connection_stage");
    this.isConnected = isConnected;
    this.message = message;
  }
}

export function addEventListener(name: "disconnected", cb: () => void): void;
export function addEventListener(name: "connected", cb: () => void): void;
export function addEventListener(name: "connection_stage", cb: (e: ConnectionEvent) => void): void;
export function addEventListener(name: "message", cb: (e: MessageEvent) => void): void;

export function addEventListener(name: string, cb: (e: any) => void): void {
  events.addEventListener(name, cb);
}

export function removeEventListener(name: string, cb: (e: any) => void): void {
  events.removeEventListener(name, cb);
}

connect();
function connect() {
  if (ws != null)
    return;

  events.dispatchEvent(new ConnectionEvent(false, "connecting..."));
  console.log("connecting to websocket");

  try {
    ws = new WebSocket(WEBSOCKET_URL);
  }
  catch(e) {
    ws = null;
    events.dispatchEvent(new ConnectionEvent(false, "connect error, retry in 5s"));
    console.error("connect error", e, "retry in 5s")
    setTimeout(connect, RETRY_INTERVAL);
    return;
  }

  ws.addEventListener("open", (e) => {
    events.dispatchEvent(new ConnectionEvent(true, "connected!"));
    console.info("connected to websocket", e);
  });

  ws.addEventListener("message", (e) => {
    try {
      const c = JSON.parse(e.data);
      console.log("message",c);
      if (!("kind" in c)) {
        throw new Error("invalid input, missing kind");
      }
      if (!("content" in c)) {
        throw new Error("invalid input, missing content");
      }
      if (c.kind == "init")
        initMessage = c;
      events.dispatchEvent(new MessageEvent(c));
    }
    catch (e) {
      console.warn("invalid json received", e);
    }
  });

  ws.addEventListener("error", (e) => {
    events.dispatchEvent(new ConnectionEvent(false, "connection error, reconnecting in 5s"));
    if (ws == null)
      return;
    console.error("websocket error", e, "reconnecting retry in 5s");
    ws = null;
    setTimeout(connect, RETRY_INTERVAL);
  });
  ws.addEventListener("close", (e) => {
    events.dispatchEvent(new ConnectionEvent(false, "connection error, reconnecting in 5s"));
    if (ws == null)
      return;
    console.info("dirconnected from the websocket", e, "reconnecting in 5s");
    ws = null;
    setTimeout(connect, RETRY_INTERVAL);
  });
}
