export const WEBSOCKET_URL: string =
  process.env.PREVIEW_WEBSOCKET_URL ?? "ws://localhost:1665/websocket";

let events: EventTarget = new EventTarget();
const RETRY_INTERVAL: number = 5000;
let ws: WebSocket | null = null;
let isWsConnected: boolean = false;

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

export type SubscribeMessage = {
  kind: "subscribe",
  content: {
    chunkPos: { x: number, y: number }
  }
}

export type UnsubscribeMessage = {
  kind: "unsubscribe",
  content: {
    chunkPos: { x: number, y: number }
  }
}

export type FullchunkMessage = {
  kind: "fullchunk",
  content: {
    chunkPos: { x: number, y: number },
    tiles: string,
  }
}

export type TileChangeMessage = {
  kind: "tileChange",
  content: {
  	tilePos: { x: number, y: number },
  	kind: number,
  	value: number,
  },
}

export type UnitMessage = {
  kind: "unit",
  content: {
    opCode: string,
    unitId: string,
    owner: string,
    position: { x: number, y: number },
  }
}

export type UnitMovementMessage = {
  kind: "unitMovement",
  content: {
    unitId: string,
    "new": { x: number, y: number },
  }
}

export type UnitDespawnedMessage = {
  kind: "unitDespawned",
  content: {
    unitId: string,
  }
}

export type PlayerSpawnMessage = {
  kind: "playerSpawn",
  content: {
  	id: string,
  	spawnPosition: { x: number, y: number },
  	username: string,
    // TODO: Describe
  	resources: unknown,
  }
}

export type PlayerDespawnMessage = {
  kind: "playerDespawn",
  content: {
  	id: string,
  }
}

export type RecvMessage =
  | InitMessage
  | FullchunkMessage
  | TileChangeMessage
  | UnitMessage
  | UnitMovementMessage
  | UnitDespawnedMessage
  | PlayerSpawnMessage
  | PlayerDespawnMessage;
export type SendMessage = SubscribeMessage | UnsubscribeMessage;

export class MessageEvent extends Event {
  public readonly message: RecvMessage;

  constructor(message: RecvMessage) {
    super("message");
    this.message = message;
  }
}

export class ConnectionEvent extends Event {
  public readonly isConnected: boolean;
  public readonly message: string;

  constructor(isConnected: boolean, message: string) {
    super("connection_event");
    this.isConnected = isConnected;
    this.message = message;
  }
}

export function addEventListener(
  name: "connection_event",
  cb: (e: ConnectionEvent) => void,
  cfg?: AddEventListenerOptions
): void;
export function addEventListener(
  name: "message",
  cb: (e: MessageEvent) => void,
  cfg?: AddEventListenerOptions
): void;

export function addEventListener(name: string, cb: (e: any) => void, cfg?: AddEventListenerOptions): void {
  events.addEventListener(name, cb, cfg);
}

export function removeEventListener(name: string, cb: (e: any) => void): void {
  events.removeEventListener(name, cb);
}

export function sendMessage(message: SendMessage) {
  if (ws == null)
  {
    console.warn("could not send message, not connected", message);
    return;
  }

  ws.send(JSON.stringify(message));
}

export function isConnected(): boolean {
  return isWsConnected;
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
    isWsConnected = true;
    events.dispatchEvent(new ConnectionEvent(true, "connected!"));
    console.info("connected to websocket", e);
  });

  ws.addEventListener("message", (e) => {
    try {
      const c = JSON.parse(e.data);
      console.debug("message", c)
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
    isWsConnected = false;
    events.dispatchEvent(new ConnectionEvent(false, "connection error, reconnecting in 5s"));
    if (ws == null)
      return;
    console.error("websocket error", e, "reconnecting retry in 5s");
    ws = null;
    setTimeout(connect, RETRY_INTERVAL);
  });
  ws.addEventListener("close", (e) => {
    isWsConnected = false;
    events.dispatchEvent(new ConnectionEvent(false, "connection error, reconnecting in 5s"));
    if (ws == null)
      return;
    console.info("dirconnected from the websocket", e, "reconnecting in 5s");
    ws = null;
    setTimeout(connect, RETRY_INTERVAL);
  });
}
