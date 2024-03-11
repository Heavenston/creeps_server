import { Tilemap } from "./map";
import * as vmod from "~/src/models/viewer"
import * as emod from "~/src/models/epita"

const RETRY_INTERVAL: number = 5000;

export function isMoveReport(report: (emod.Report & {opcode: string})): report is emod.MoveReport {
  return report.opcode.startsWith("move:");
}

export class MessageEvent extends Event {
  public readonly message: vmod.S2CMessage;

  constructor(message: vmod.S2CMessage) {
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
  #initMessage: vmod.S2CInit | null = null;
  get initMessage(): vmod.S2CInit | null {
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
        const c = JSON.parse(e.data) as vmod.S2CMessage;
        console.debug("message", c)
        if (!("kind" in c)) {
          throw new Error("invalid input, missing kind");
        }
        if (!("content" in c)) {
          throw new Error("invalid input, missing content");
        }
        if (c.kind == "init")
          this.#initMessage = c.content;
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

  getActionCost(opcode: string): emod.CostResponse | null {
    const initMsg = this.#initMessage;
    if (!initMsg)
      return null;

    let name: keyof emod.CostsResponse|null = null;
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
    return initMsg.costs[name] ?? null;
  }

  secondsPerTicks(): number {
    if (!this.#initMessage)
      return 1;
    return 1 / this.#initMessage.setup.ticksPerSecond;
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

  sendMessage(message: vmod.C2SMessage) {
    if (this.#ws == null) {
      console.warn("could not send message, not connected", message);
      return;
    }
  
    this.#ws.send(JSON.stringify(message));
  }
}
