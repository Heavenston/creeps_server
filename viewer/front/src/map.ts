import { vec, Vector2 } from "./geom"
import * as api from "./api"

export class Chunk {
  private tiles: Uint8Array;

  public constructor(data: string) {
    this.tiles = Uint8Array.from(atob(data), c => c.charCodeAt(0));
  }

  public static get chunkSize(): number {
    return api.getInitMessage()?.content.chunkSize ?? 8;
  }

  private getIndex(pos: Vector2): number {
    return (Math.floor(pos.x) + Math.floor(pos.y) * Chunk.chunkSize) * 2;
  }

  public getTileKind(pos: Vector2): number {
    const index = this.getIndex(pos);
    return this.tiles[index];
  }

  public getTileValue(pos: Vector2): number {
    const index = this.getIndex(pos);
    return this.tiles[index+1];
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

const subs = new Set<VecKey>()
const chunks = new Map<VecKey, Chunk>()

api.addEventListener("connection_event", ce => {
  if (!ce.isConnected) {
    subs.clear();
    chunks.clear();
  }
});

api.addEventListener("message", event => {
  if (event.message.kind == "fullchunk") {
    const content = event.message.content;
    const pos = vec(content.chunkPos.x, content.chunkPos.y);
    chunks.set(key(pos), new Chunk(content.tiles));
  }
})

export function subscribe(pos: Vector2) {
  if (!api.isConnected()) {
    console.warn("Tried to subscribe to", pos, "while not connected")
    return;
  }
  if (subs.has(key(pos))) {
    console.warn("Tried to resubscribe to", pos)
    return;
  }

  subs.add(key(pos));
  api.sendMessage({
    kind: "subscribe",
    content: {
      chunkPos: { x: pos.x, y: pos.y }
    }
  });
}

export function unsubscribe(pos: Vector2) {
  if (!api.isConnected()) {
    console.warn("Tried to unsubscribe to", pos, "while not connected")
    return;
  }
  if (!subs.has(key(pos))) {
    console.warn("Tried to reunsubscribe to", pos)
    return;
  }

  subs.delete(key(pos));
  api.sendMessage({
    kind: "unsubscribe",
    content: {
      chunkPos: { x: pos.x, y: pos.y }
    }
  });
}

// Make sure the only subscribed chunks are the one listed
export function setSubscribed(subed: Vector2[]) {
  const toUnsub = new Set(subs);
  for (const sub of subed) {
    toUnsub.delete(key(sub));
    if (!subs.has(key(sub)))
      subscribe(sub);
  }

  for (const unsub of toUnsub) {
    unsubscribe(fromKey(unsub));
  }
}

// euclidian remained
function remEuclid(a: number, b: number): number {
  return ((a % b) + b) % b;
}

// see the go server same function lol
export function global2ContainingChunkCoords(global: Vector2): Vector2 {
  return global
    .mapped(Math.round)
    .mapped(x => Math.floor(x / Chunk.chunkSize));
}

export function global2ChunkSubCoords(global: Vector2): Vector2 {
  return global
    .mapped(Math.round)
    .mapped(a => remEuclid(a, Chunk.chunkSize))
}

/// get the tile at the given global position or -1 if unavailable
export function getTileKind(globalPos: Vector2): number {
  const chunkPos = global2ContainingChunkCoords(globalPos);
  const chunkSubPos = global2ChunkSubCoords(globalPos);

  const chunk = chunks.get(key(chunkPos));
  if (!chunk)
    return -1;

  return chunk.getTileKind(chunkSubPos);
}

/// get the tile at the given global position or -1 if unavailable
export function getTileValue(globalPos: Vector2): number {
  const chunkPos = global2ContainingChunkCoords(globalPos);
  const chunkSubPos = global2ChunkSubCoords(globalPos);

  const chunk = chunks.get(key(chunkPos));
  if (!chunk)
    return -1;

  return chunk.getTileValue(chunkSubPos);
}
