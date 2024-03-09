import { vec, Vector2 } from "~/src/utils/geom"
import { Api } from "./api";

function remEuclid(a: number, b: number): number {
  return ((a % b) + b) % b;
}

export class Chunk {
  private tiles: Uint8Array;

  public constructor(public readonly map: Tilemap, data: string, public readonly pos: Vector2) {
    this.tiles = Uint8Array.from(atob(data), c => c.charCodeAt(0));
  }

  public get chunkSize(): number {
    return this.map.chunkSize;
  }

  private getIndex(pos: Vector2): number {
    return (Math.floor(pos.x) + Math.floor(pos.y) * this.chunkSize) * 2;
  }

  public getTileKind(pos: Vector2): number {
    const index = this.getIndex(pos);
    return this.tiles[index];
  }

  public getTileValue(pos: Vector2): number {
    const index = this.getIndex(pos);
    return this.tiles[index+1];
  }

  public updateTile(pos: Vector2, kind: number, value: number) {
    const index = this.getIndex(pos);
    this.tiles[index] = kind;
    this.tiles[index+1] = value;
  }
}

type VecKey = `${number}_${number}`;
function key(pos: Vector2): VecKey {
  return `${pos.x}_${pos.y}`;
}
function fromKey(key: VecKey): Vector2 {
  const [a, b] = key.split("_").map(x => +x);
  return vec(a, b);
}

export class Tilemap {
  #subs = new Set<VecKey>();
  #chunks: Map<VecKey, Chunk> = new Map();

  constructor(public readonly api: Api) {
    api.addEventListener("connection_event", ce => {
      if (!ce.isConnected) {
        this.#subs.clear();
        this.#chunks.clear();
      }
    });

    api.addEventListener("message", event => {
      if (event.message.kind == "fullchunk") {
        const content = event.message.content;
        const pos = vec(content.chunkPos);
        this.#chunks.set(key(pos), new Chunk(this, content.tiles, pos));
      }
      if (event.message.kind == "tileChange") {
        const content = event.message.content;
        const tilePos = this.global2ChunkSubCoords(vec(content.tilePos));
        const chunkPos = this.global2ContainingChunkCoords(vec(content.tilePos));
        const chunk = this.#chunks.get(key(chunkPos))
        if (!chunk) {
          console.warn("received tile for unkown chunk");
          return
        }
        chunk.updateTile(tilePos, content.kind, content.value);
      }
    })
  }

  public get chunkSize(): number {
    return this.api.initMessage?.content.chunkSize ?? 8;
  }

  global2ContainingChunkCoords(global: Vector2): Vector2 {
    return global
      .mapped(Math.floor)
      .mapped(x => Math.floor(x / this.chunkSize));
  }

  global2ChunkSubCoords(global: Vector2): Vector2 {
    return global
      .mapped(Math.floor)
      .mapped(a => remEuclid(a, this.chunkSize))
  }

  subscribe(pos: Vector2) {
    const k = key(pos);

    if (!this.api.isConnected) {
      console.warn("Tried to subscribe to", pos, "while not connected")
      return;
    }
    if (this.#subs.has(k)) {
      console.warn("Tried to resubscribe to", pos)
      return;
    }

    this.#subs.add(k);
    this.api.sendMessage({
      kind: "subscribe",
      content: {
        chunkPos: { x: pos.x, y: pos.y }
      }
    });
  }

  unsubscribe(pos: Vector2) {
    const k = key(pos);

    if (!this.api.isConnected) {
      console.warn("Tried to unsubscribe to", pos, "while not connected")
      return;
    }
    if (!this.#subs.has(k)) {
      console.warn("Tried to reunsubscribe to", pos)
      return;
    }

    this.#subs.delete(k);
    this.#chunks.delete(k);
    this.api.sendMessage({
      kind: "unsubscribe",
      content: {
        chunkPos: { x: pos.x, y: pos.y }
      }
    });
  }

  // Make sure the only subscribed chunks are the one listed
  setSubscribed(subed: Vector2[]) {
    const toUnsub = new Set(this.#subs);
    for (const sub of subed) {
      toUnsub.delete(key(sub));
      if (!this.#subs.has(key(sub)))
        this.subscribe(sub);
    }

    for (const unsub of toUnsub) {
      this.unsubscribe(fromKey(unsub));
    }
  }

  getChunk(chunkPos: Vector2): Chunk | null {
    return this.#chunks.get(key(chunkPos)) ?? null;
  }

  /// get the tile at the given global position or -1 if unavailable
  getTileKind(globalPos: Vector2): number {
    const chunkPos = this.global2ContainingChunkCoords(globalPos);
    const chunkSubPos = this.global2ChunkSubCoords(globalPos);
    console.log({
      globalPos: key(globalPos),
      chunkPos: key(chunkPos),
      chunkSubPos: key(chunkSubPos),
    })

    const chunk = this.#chunks.get(key(chunkPos));
    if (!chunk)
      return -1;

    return chunk.getTileKind(chunkSubPos);
  }

  /// get the tile at the given global position or -1 if unavailable
  getTileValue(globalPos: Vector2): number {
    const chunkPos = this.global2ContainingChunkCoords(globalPos);
    const chunkSubPos = this.global2ChunkSubCoords(globalPos);

    const chunk = this.#chunks.get(key(chunkPos));
    if (!chunk)
      return -1;

    return chunk.getTileValue(chunkSubPos);
  }
}


