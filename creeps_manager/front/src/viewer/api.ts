import { Tilemap } from "./map";

const RETRY_INTERVAL: number = 5000;

export type Point = {
  x: number,
  y: number,
}

export type Action = {
  actionOpCode: string,
  reportId: string,
  // hahahah...
  parameter?: any,
}

export type CreepsReport = {
  reportId: string,
  unitId: string,
  login: string,
  unitPosition: Point,
  status: "SUCCESS" | "ERROR",
}

export type MoveReport = CreepsReport & {
  opcode: `move:${string}`,
  newPosition: Point,
}

export function isMoveReport(report: (CreepsReport & {opcode: string})): report is MoveReport {
  return report.opcode.startsWith("move:");
}

export type AnyReport = MoveReport;

export type Resources = {
	rock: number,
	wood: number,
	food: number,
	oil: number,
	copper: number,
	woodPlank: number,
}

export type Cost = Resources & {
  cast: number,
}

export type Costs = {
  [action: string]: undefined | Cost,
}

export type InitMessage = {
  kind: "init",
  content: {
    chunkSize: number,
    costs: Costs,
    // TODO: describe
    setup: {
      ticksPerSecond: number,
    },
  }
}

export type SubscribeMessage = {
  kind: "subscribe",
  content: {
    chunkPos: Point
  }
}

export type UnsubscribeMessage = {
  kind: "unsubscribe",
  content: {
    chunkPos: Point
  }
}

export type FullchunkMessage = {
  kind: "fullchunk",
  content: {
    chunkPos: Point,
    tiles: string,
  }
}

export type TileChangeMessage = {
  kind: "tileChange",
  content: {
  	tilePos: Point,
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
    position: Point,
    upgraded: boolean,
  }
}

export type UnitDespawnedMessage = {
  kind: "unitDespawned",
  content: {
    unitId: string,
  }
}

export type UnitStartedActionMessage = {
  kind: "unitStartedAction",
  content: {
    unitId: string,
    action: Action,
  }
}

export type UnitFinishedActionMessage = {
  kind: "unitFinishedAction",
  content: {
    unitId: string,
    action: Action,
    report: AnyReport,
  }
}

export type PlayerSpawnMessage = {
  kind: "playerSpawn",
  content: {
  	id: string,
  	spawnPosition: Point,
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
  | UnitDespawnedMessage
  | UnitStartedActionMessage
  | UnitFinishedActionMessage
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

export class Api { 
  public readonly url: string;

  #closed: boolean = false;
  get closed(): boolean {
    return this.#closed;
  }

  #events: EventTarget = new EventTarget();
  #ws: WebSocket | null = null;
  #isConnected: boolean = false;
  get isConnected(): boolean {
    return this.#isConnected;
  }
  #initMessage: InitMessage | null = null;
  get initMessage(): InitMessage | null {
    return this.#initMessage;
  }
  #tilemap: Tilemap;
  get tilemap(): Tilemap {
    return this.#tilemap;
  }

  constructor(url: string) {
    this.url = url;
    this.connect();
    this.#tilemap = new Tilemap(this);
  }

  public close() {
    if (this.#closed) {
      console.warn("tried to double close api");
      return;
    }
    this.#ws?.close();
    this.#ws = null;
    this.#closed = true;
  }

  private connect() {
    if (this.#closed)
      return;
    if (this.#ws != null)
      return;

    this.#events.dispatchEvent(new ConnectionEvent(false, "connecting..."));
    console.log("connecting to websocket");

    try {
      this.#ws = new WebSocket(this.url);
    }
    catch(e) {
      this.#ws = null;
      this.#events.dispatchEvent(new ConnectionEvent(false, "connect error, retry in 5s"));
      console.error("connect error", e, "retry in 5s")
      setTimeout(this.connect.bind(this), RETRY_INTERVAL);
      return;
    }

    this.#ws.addEventListener("open", (e) => {
      this.#isConnected = true;
      this.#events.dispatchEvent(new ConnectionEvent(true, "connected!"));
      console.info("connected to websocket", e);
    });

    this.#ws.addEventListener("message", (e) => {
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
          this.#initMessage = c;
        this.#events.dispatchEvent(new MessageEvent(c));
      }
      catch (e) {
        console.warn("invalid json received", e);
      }
    });

    this.#ws.addEventListener("error", (e) => {
      this.#isConnected = false;
      this.#events.dispatchEvent(new ConnectionEvent(false, "connection error, reconnecting in 5s"));
      if (this.#ws == null)
        return;
      console.error("websocket error", e, "reconnecting retry in 5s");
      this.#ws = null;
      setTimeout(this.connect.bind(this), RETRY_INTERVAL);
    });

    this.#ws.addEventListener("close", (e) => {
      this.#isConnected = false;
      this.#events.dispatchEvent(new ConnectionEvent(false, "connection error, reconnecting in 5s"));
      if (this.#ws == null)
        return;
      console.info("dirconnected from the websocket", e, "reconnecting in 5s");
      this.#ws = null;
      setTimeout(this.connect.bind(this), RETRY_INTERVAL);
    });
  }

  getActionCost(opcode: string): Cost | null {
    const initMsg = this.#initMessage;
    if (!initMsg)
      return null;

    let name = null;
    switch (opcode) {
      case "move:left":
      case "move:right":
      case "move:up":
      case "move:down":
        name = "move";
        break
      case "fire:turret":
        name = "fireTurret";
        break
    }
    if (name == null)
      return null;
    return initMsg.content.costs[name] ?? null;
  }

  secondsPerTicks(): number {
    if (!this.#initMessage)
      return 1;
    return 1 / this.#initMessage.content.setup.ticksPerSecond;
  }

  public addEventListener(
    name: "connection_event",
    cb: (e: ConnectionEvent) => void,
    cfg?: AddEventListenerOptions
  ): void;
  public addEventListener(
    name: "message",
    cb: (e: MessageEvent) => void,
    cfg?: AddEventListenerOptions
  ): void;

  addEventListener(name: string, cb: (e: any) => void, cfg?: AddEventListenerOptions): void {
    this.#events.addEventListener(name, cb, cfg);
  }

  removeEventListener(name: string, cb: (e: any) => void): void {
    this.#events.removeEventListener(name, cb);
  }

  sendMessage(message: SendMessage) {
    if (this.#ws == null) {
      console.warn("could not send message, not connected", message);
      return;
    }
  
    this.#ws.send(JSON.stringify(message));
  }
}
