import { vec, Vector2 } from "~/src/geom"
import * as api from "~/src/api"
import * as map from "./map"
import { IRenderer, Renderer } from "./worldRenderer";

export class TerrainRenderer implements IRenderer {
  private readonly renderer: Renderer;

  private eventAbort = new AbortController();

  public chunksOnCamera: Vector2[] = [];
  private chunksCanvases: WeakMap<map.Chunk, OffscreenCanvas> = new WeakMap();
  private enableChunkBorder = false;

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(renderer: Renderer) {
    this.renderer = renderer;

    document.body.addEventListener("keydown", k => {
      if (k.key == "r") {
        this.chunksCanvases = new WeakMap();
      }
      if (k.key == "c") {
        this.enableChunkBorder = !this.enableChunkBorder;
      }
    }, {
      signal: this.eventAbort.signal,
    });

    api.addEventListener("message", event => {
      if (event.message.kind != "fullchunk")      
        return;
      const pos = vec(event.message.content.chunkPos);
      const chunk = map.getChunk(pos)
      if (!chunk)
        return;
      // force redraw
      this.chunksCanvases.delete(chunk);
    }, {
      signal: this.eventAbort.signal,
    });

    api.addEventListener("message", event => {
      if (event.message.kind != "tileChange")      
        return;
      const chunkPos = map.global2ContainingChunkCoords(vec(event.message.content.tilePos));
      const chunk = map.getChunk(chunkPos);
      if (!chunk)
        return;
      // force redraw
      this.chunksCanvases.delete(chunk);
    }, {
      signal: this.eventAbort.signal,
    });

    this.renderer.texturePack.addEventListener("textureLoaded", () => {
      this.chunksCanvases = new WeakMap();
    }, {
      signal: this.eventAbort.signal,
    });
  }

  private lastChunkUpadeCameraPos = vec(-5888888, -588888);
  private update(_dt: number) {
    const cameraPos = this.renderer.cameraPos;

    // loaded chunks update
    if (this.lastChunkUpadeCameraPos.x == cameraPos.x && this.lastChunkUpadeCameraPos.y == cameraPos.y)
      return;
    this.lastChunkUpadeCameraPos = vec(cameraPos);

    const chunksOnCamera: Vector2[] = [];
    this.chunksOnCamera = chunksOnCamera;

    const start = this.renderer.screenTopLeftInWorldPos;
    const end = this.renderer.screenBottomRightInWorldPos;
    // console.log({start, end})
    const cp = vec(start);
    for (cp.x = start.x; cp.x-map.Chunk.chunkSize < end.x; cp.x += map.Chunk.chunkSize) {
      for (cp.y = start.y; cp.y-map.Chunk.chunkSize < end.y; cp.y += map.Chunk.chunkSize) {
        chunksOnCamera.push(map.global2ContainingChunkCoords(cp));
      }
    }

    map.setSubscribed(chunksOnCamera)
  }

  private renderChunkCanvas(chunk: map.Chunk): OffscreenCanvas {
    const ts = this.renderer.texturePack.size;
    const canvas = new OffscreenCanvas(
      map.Chunk.chunkSize * ts,
      map.Chunk.chunkSize * ts,
    );
    const ctx = canvas.getContext("2d");
    if (ctx == undefined)
      throw new Error("unsupported device");
    ctx.imageSmoothingEnabled = false;

    this.chunksCanvases.set(chunk, canvas);

    ctx.fillStyle = this.renderer.texturePack.fillColor;
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    for (let sx = 0; sx < map.Chunk.chunkSize; sx++) {
      for (let sy = 0; sy < map.Chunk.chunkSize; sy++) {
        const subTileCoord = vec(sx, sy);
        const globalTileCoord = chunk.pos.times(map.Chunk.chunkSize).plus(subTileCoord);

        const value = chunk.getTileKind(subTileCoord)

        const texture = this.renderer.texturePack.getTileTexture(value, globalTileCoord);
        ctx.drawImage(texture, subTileCoord.x * ts, subTileCoord.y * ts);
      }
    }

    return canvas;
  }

  private renderChunk(pos: Vector2) {
    // const start = this.screenTopLeftInWorldPos;
    // const end = this.screenBottomRightInWorldPos;

    const chunk = map.getChunk(pos);
    if (chunk == null)
      return;

    let canvas = this.chunksCanvases.get(chunk);
    if (!canvas)
      canvas = this.renderChunkCanvas(chunk);

    const drawpos = pos.times(map.Chunk.chunkSize);
    // console.log(pos, drawpos);

    const ctx = this.renderer.ctx;

    ctx.imageSmoothingEnabled = false;
    ctx.drawImage(canvas, drawpos.x, drawpos.y, map.Chunk.chunkSize, map.Chunk.chunkSize);

    if (this.enableChunkBorder) {
      ctx.strokeStyle = "black";
      ctx.lineWidth = 0.1;
      ctx.strokeRect(drawpos.x, drawpos.y, map.Chunk.chunkSize, map.Chunk.chunkSize);
    }
  }

  public render(dt: number) {
    if (dt != 0)
      this.update(dt);

    for (const chunk of this.chunksOnCamera)
      this.renderChunk(chunk);
  }
}

