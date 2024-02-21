import { vec, Vector2 } from "~/src/geom"
import * as api from "~/src/api"
import * as map from "./map"
import { OverlayRenderer } from "./overlayRenderer";

export class Renderer {
  public readonly canvas: HTMLCanvasElement;
  public readonly ctx: CanvasRenderingContext2D;

  // position of the center of the screen in world coordinate
  public cameraPos: Vector2 = vec(0, 0);
  // scale to go from screen pos to world pos
  public cameraScale: number = 25;

  // position of the mouse in screen coordinated
  public mousePos: Vector2 = vec(0, 0);

  private eventAbort = new AbortController();

  public chunksOnCamera: Vector2[] = [];
  private chunksCanvases: WeakMap<map.Chunk, OffscreenCanvas> = new WeakMap();

  private lastUnitMessage: Map<string, api.UnitMessage> = new Map();

  private overlayRenderer = new OverlayRenderer(this);

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
      .times(1/this.cameraScale)
      .plus(this.screenTopLeftInWorldPos);
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
    this.overlayRenderer.cleanup();
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
      const tile = this.mouseWorldPos.mapped(Math.floor);
      console.log("Cliked tile: ", {
        position: [this.mouseWorldPos.x, this.mouseWorldPos.y].join(" "),
        flooredPosition: [tile.x, tile.y].join(" "),
        kind: map.getTileKind(tile),
        value: map.getTileValue(tile),
      });
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

    api.addEventListener("message", event => {
      if (event.message.kind != "unitDespawned")      
        return;
      this.lastUnitMessage.delete(event.message.content.unitId);
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
    // console.log(pos, drawpos);

    this.ctx.imageSmoothingEnabled = false;
    this.ctx.drawImage(canvas, drawpos.x, drawpos.y);
    // this.ctx.strokeStyle = "black"
    // this.ctx.lineWidth = 5 / this.cameraScale;
    // this.ctx.strokeRect(drawpos.x, drawpos.y, map.Chunk.chunkSize, map.Chunk.chunkSize)
  }

  private renderUnit(unit: api.UnitMessage) {
    const pos = vec(unit.content.position.x, unit.content.position.y);
    if (unit.content.opCode == "turret") {
      this.ctx.fillStyle = "magenta";
      this.ctx.fillRect(pos.x + 0.125, pos.y + 0.125, 0.75, 0.75);
    }
    else {
      this.ctx.fillStyle = "orange";
      this.ctx.fillRect(pos.x + 0.25, pos.y + 0.25, 0.5, 0.5);
    }
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

    for (const unit of this.lastUnitMessage.values())
      if (unit.content.opCode == "turret")
        this.renderUnit(unit);
    for (const unit of this.lastUnitMessage.values())
      if (unit.content.opCode != "turret")
        this.renderUnit(unit);

    this.overlayRenderer.render(dt);
  }
}