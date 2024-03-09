import { vec, Vector2 } from "~/src/utils/geom"
import * as map from "./map"
import { IRenderer, Renderer } from "./worldRenderer";
import { Api } from "./api";

export class TerrainRenderer implements IRenderer {
  private readonly renderer: Renderer;

  private eventAbort = new AbortController();

  public chunksOnCamera: Vector2[] = [];
  private chunksCurrentlyRendering: Set<map.Chunk> = new Set();
  private chunksCanvases: WeakMap<map.Chunk, OffscreenCanvas> = new WeakMap();
  private enableChunkBorder = false;

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(renderer: Renderer, public readonly api: Api) {
    this.renderer = renderer;

    document.body.addEventListener("keydown", k => {
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
      const chunk = api.tilemap.getChunk(pos)
      if (!chunk)
        return;
      // force redraw
      this.startRenderChunkCanvas(chunk);
    }, {
      signal: this.eventAbort.signal,
    });

    api.addEventListener("message", event => {
      if (event.message.kind != "tileChange")      
        return;
      const chunkPos = api.tilemap.global2ContainingChunkCoords(vec(event.message.content.tilePos));
      const chunk = api.tilemap.getChunk(chunkPos);
      if (!chunk)
        return;
      // force redraw
      this.startRenderChunkCanvas(chunk);
    }, {
      signal: this.eventAbort.signal,
    });

    this.renderer.texturePack.addEventListener("textureLoaded", () => {
      for (const chunkPos of this.chunksOnCamera)
      {
        const chunk = api.tilemap.getChunk(chunkPos);
        if (chunk != null)
          this.startRenderChunkCanvas(chunk);
      }
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
    for (cp.x = start.x; cp.x-this.api.tilemap.chunkSize < end.x; cp.x += this.api.tilemap.chunkSize) {
      for (cp.y = start.y; cp.y-this.api.tilemap.chunkSize < end.y; cp.y += this.api.tilemap.chunkSize) {
        chunksOnCamera.push(this.api.tilemap.global2ContainingChunkCoords(cp));
      }
    }

    this.api.tilemap.setSubscribed(chunksOnCamera)
  }

  private startRenderChunkCanvas(chunk: map.Chunk) {
    if (this.chunksCurrentlyRendering.has(chunk))
      return;
    this.chunksCurrentlyRendering.add(chunk);
    // schedule tack for later to avoid taking time during the render
    setTimeout(() => {
      const ts = this.renderer.texturePack.size;
      const canvas = new OffscreenCanvas(
        this.api.tilemap.chunkSize * ts,
        this.api.tilemap.chunkSize * ts,
      );
      const ctx = canvas.getContext("2d");
      if (ctx == undefined)
        throw new Error("unsupported device");
      ctx.imageSmoothingEnabled = false;

      this.chunksCanvases.set(chunk, canvas);

      ctx.fillStyle = this.renderer.texturePack.fillColor;
      ctx.fillRect(0, 0, canvas.width, canvas.height);

      for (let sx = 0; sx < this.api.tilemap.chunkSize; sx++) {
        for (let sy = 0; sy < this.api.tilemap.chunkSize; sy++) {
          const subTileCoord = vec(sx, sy);
          const globalTileCoord = chunk.pos.times(this.api.tilemap.chunkSize).plus(subTileCoord);

          const value = chunk.getTileKind(subTileCoord)

          const texture = this.renderer.texturePack.getTileTexture(value, globalTileCoord);
          ctx.drawImage(texture, subTileCoord.x * ts, subTileCoord.y * ts);
        }
      }

      this.chunksCurrentlyRendering.delete(chunk);
      this.chunksCanvases.set(chunk, canvas);
    });
  }

  private renderChunk(pos: Vector2) {
    // const start = this.screenTopLeftInWorldPos;
    // const end = this.screenBottomRightInWorldPos;

    const chunk = this.api.tilemap.getChunk(pos);
    if (chunk == null)
      return;

    let canvas = this.chunksCanvases.get(chunk);
    if (!canvas)
    {
      this.startRenderChunkCanvas(chunk);
      return;
    }

    const drawpos = pos.times(this.api.tilemap.chunkSize);
    // console.log(pos, drawpos);

    const ctx = this.renderer.ctx;

    ctx.imageSmoothingEnabled = false;
    ctx.drawImage(canvas, drawpos.x, drawpos.y, this.api.tilemap.chunkSize, this.api.tilemap.chunkSize);

    if (this.enableChunkBorder) {
      ctx.strokeStyle = "black";
      ctx.lineWidth = 0.1;
      ctx.strokeRect(drawpos.x, drawpos.y, this.api.tilemap.chunkSize, this.api.tilemap.chunkSize);
    }
  }

  public render(dt: number) {
    if (dt != 0)
      this.update(dt);

    for (const chunk of this.chunksOnCamera)
      this.renderChunk(chunk);
  }
}

