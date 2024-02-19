import { MinzeElement } from "minze"
import { vec, Vector2 } from "./geom"
import * as api from "./api"
import * as map from "./map"

class WorldRenderer {
  private readonly canvas: HTMLCanvasElement;
  private readonly ctx: CanvasRenderingContext2D;

  // position of the center of the screen in world coordinate
  public cameraPos: Vector2 = vec(0, 0);
  // scale to go from screen pos to world pos
  public cameraScale: number = 25;

  // position of the mouse in screen coordinated
  private mousePos: Vector2 = vec(0, 0);

  private eventAbort = new AbortController();

  private chunksOnCamera: Vector2[] = [];
  private chunksCanvases: WeakMap<map.Chunk, OffscreenCanvas> = new WeakMap();

  private lastUnitMessage: Map<string, api.UnitMessage> = new Map();

  private get screenTopLeftInWorldPos(): Vector2 {
    return this.cameraPos
      .minus(vec(this.canvas.width, this.canvas.height).times(0.5).times(1/this.cameraScale));
  }

  private get screenBottomRightInWorldPos(): Vector2 {
    return this.cameraPos
      .plus(vec(this.canvas.width, this.canvas.height).times(0.5).times(1/this.cameraScale));
  }

  private get mouseWorldPos(): Vector2 {
    return this.mousePos
      .minus(vec(this.canvas.width, this.canvas.height).times(0.5))
      .times(1/this.cameraScale);
  }

  // changes the scale but also changes the cameraPos making sure the mousePos
  // doesn't change what it is pointing at
  private changeScale(val: number) {
    const adjustedPos = this.mousePos.minus(vec(this.canvas.width, this.canvas.height).times(0.5));
    const prevGobal = adjustedPos.times(1/this.cameraScale).plus(this.cameraPos);
    const newGlobal = adjustedPos.times(1/val).plus(this.cameraPos);

    this.cameraPos.sub(newGlobal.minus(prevGobal));
    this.cameraScale = val;
  }

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      alert("device/browser not supported");
      throw new Error("could not create context");
    }

    let clickMouseStart: Vector2 | null = null;
    let clickCameraStart: Vector2 | null = null;
    this.canvas.addEventListener("mousedown", ev => {
      console.log(this.mouseWorldPos.mapped(Math.floor));
      clickMouseStart = vec(ev.clientX, ev.clientY);
      clickCameraStart = vec(this.cameraPos);
    }, {
      signal: this.eventAbort.signal,
    });

    this.canvas.addEventListener("mousemove", ev => {
      this.mousePos = vec(ev.clientX, ev.clientY);
      if (clickMouseStart == null || clickCameraStart == null)
        return;
      const diff = clickMouseStart.minus(ev.clientX, ev.clientY);
      this.cameraPos = clickCameraStart.plus(diff.times(1 / this.cameraScale));
    }, {
      signal: this.eventAbort.signal,
    })

    this.canvas.addEventListener("mouseup", () => {
      clickCameraStart = null;
      clickMouseStart = null;
    }, {
      signal: this.eventAbort.signal,
    });

    this.canvas.addEventListener("mouseleave", () => {
      clickCameraStart = null;
      clickMouseStart = null;
    }, {
      signal: this.eventAbort.signal,
    });

    this.canvas.addEventListener("wheel", e => {
      const sign = e.deltaY < 0 ? -1 : 1;
      if (sign > 0)
        this.changeScale(this.cameraScale * 0.8);
      else
        this.changeScale(this.cameraScale * 1.2);
    }, {
      signal: this.eventAbort.signal,
    });

    document.body.addEventListener("keydown", k => {
      if (k.key == "r") {
        this.cameraPos = vec(0, 0);
        this.cameraScale = 25;
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
      if (event.message.kind != "unit")      
        return;
      this.lastUnitMessage.set(event.message.content.unitId, event.message);
    }, {
      signal: this.eventAbort.signal,
    });

    this.ctx = ctx;
  }

  private lastChunkUpadeCameraPos = vec(-5888888, -588888);
  private update(_dt: number) {
    // this.cameraPos = vec(this.canvas.width, this.canvas.height).times(0.5);

    if (this.lastChunkUpadeCameraPos.x == this.cameraPos.x && this.lastChunkUpadeCameraPos.y == this.cameraPos.y)
      return;
    this.lastChunkUpadeCameraPos = vec(this.cameraPos);

    const chunksOnCamera: Vector2[] = [];
    this.chunksOnCamera = chunksOnCamera;

    const start = this.screenTopLeftInWorldPos;
    const end = this.screenBottomRightInWorldPos;
    // console.log({start, end})
    const cp = vec(start);
    for (cp.x = start.x; cp.x-map.Chunk.chunkSize < end.x; cp.x += map.Chunk.chunkSize) {
      for (cp.y = start.y; cp.y-map.Chunk.chunkSize < end.y; cp.y += map.Chunk.chunkSize) {
        chunksOnCamera.push(map.global2ContainingChunkCoords(cp));
      }
    }

    // console.log("----");
    // for (const c of chunksOnCamera)
    //   console.log(c);
    // console.log("----");
    map.setSubscribed(chunksOnCamera)
  }

  private renderChunkCanvas(chunk: map.Chunk): OffscreenCanvas {
    const canvas = new OffscreenCanvas(map.Chunk.chunkSize, map.Chunk.chunkSize);
    const ctx = canvas.getContext("2d");
    if (ctx == undefined)
      throw new Error("unsupported device");
    ctx.imageSmoothingEnabled = false;

    this.chunksCanvases.set(chunk, canvas);

    for (let sx = 0; sx < map.Chunk.chunkSize; sx++) {
      for (let sy = 0; sy < map.Chunk.chunkSize; sy++) {
        const subTileCoord = vec(sx, sy);

        const value = chunk.getTileKind(subTileCoord)

        let style: string;
        switch (value) {
        case 0:
          style = "green";
          break;
        case 1:
          style = "blue";
          break;
        case 2:
          style = "gray";
          break;
        case 3:
          style = "lime";
          break;
        case 4:
          style = "red";
          break;
        case 5:
          style = "black";
          break;
        default:
          style = "yellow";
          break;
        }

        ctx.fillStyle = style;
        ctx.fillRect(
          subTileCoord.x, subTileCoord.y,
          1, 1,
        );
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

    this.ctx.imageSmoothingEnabled = false;
    this.ctx.drawImage(canvas, drawpos.x, drawpos.y);
  }

  public render(dt: number) {
    if (this.canvas == null || this.ctx == null)
      return;
    if (dt != 0)
      this.update(dt);

    this.ctx.resetTransform();

    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    this.ctx.translate(
      this.canvas.width/2,
      this.canvas.height/2,
    );
    this.ctx.transform(
      this.cameraScale, 0, 0, this.cameraScale, 0, 0
    );
    this.ctx.translate(
      -this.cameraPos.x,
      -this.cameraPos.y,
    );

    for (const chunk of this.chunksOnCamera)
      this.renderChunk(chunk);

    for (const unit of this.lastUnitMessage.values()) {
      const pos = vec(unit.content.position.x, unit.content.position.y);
      this.ctx.fillStyle = "orange";
      this.ctx.fillRect(pos.x + 0.25, pos.y + 0.25, 0.5, 0.5);
    }
  }
}

(class extends MinzeElement {
  // html template
  html = () => `<canvas/> `

  // scoped stylesheet
  css = () => `
  :host {
    height: 100%;
  }
  `

  private canvas: HTMLCanvasElement | null = null;
  private renderer: WorldRenderer | null = null;
  private animationFrameId: number = -1;

  private lastTime = 0;

  private eventAbort = new AbortController();

  private resizeCanvas() {
    if (this.canvas == null)
      return;
    this.canvas.width = this.clientWidth;
    this.canvas.height = this.clientHeight;

    this.renderCanvas(this.lastTime);
  }

  private renderCanvas(time: number) {
    if (this.renderer == null)
      return;
    this.animationFrameId = -1;

    this.renderer?.render(this.lastTime - time);
    this.lastTime = time;

    if (this.animationFrameId != -1)
      cancelAnimationFrame(this.animationFrameId);
    this.animationFrameId = requestAnimationFrame(this.renderCanvas.bind(this));
  }

  onReady() {
    this.canvas = this.select("canvas") ?? document.createElement("canvas");
    const ctx = this.canvas.getContext("2d")
    if (ctx == null) {
      alert("unsupported device");
      return;
    }

    api.addEventListener("connection_event", c => {
      if (c.isConnected) {
        if (this.canvas != null)
          this.renderer = new WorldRenderer(this.canvas);
        this.resizeCanvas();
      }
      else {
        this.renderer?.cleanup();
        this.renderer = null;
      }
    }, {
      signal: this.eventAbort.signal
    });

    const ro = new ResizeObserver(() => {
      this.resizeCanvas();
    });
    ro.observe(this);
    this.eventAbort.signal.addEventListener("abort", ro.disconnect.bind(ro));
  }

  onDestroy() {
    if (this.renderer)
      this.renderer.cleanup();
    cancelAnimationFrame(this.animationFrameId);
    this.eventAbort.abort();
  }
}).define("creeps-canvas")
