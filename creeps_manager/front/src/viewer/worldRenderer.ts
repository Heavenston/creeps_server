import { vec, Vector2 } from "~/src/utils/geom"
import { OverlayRenderer } from "./overlayRenderer";
import { TexturePack } from "./texturePack";
import { UnitRenderer } from "./unitRenderer";
import { TerrainRenderer } from "./terrainRenderer";
import { Api } from "./api";

export interface IRenderer {
  render(dt: number): void;
  cleanup(): void;
}

export class Renderer {
  public readonly canvas: HTMLCanvasElement;
  public readonly ctx: CanvasRenderingContext2D;
  public readonly texturePack = new TexturePack();

  // position of the center of the screen in world coordinate
  public cameraPos: Vector2 = vec(0, 0);
  // scale to go from screen pos to world pos
  public cameraScale: number = 25;
  // position of the mouse in screen coordinated
  public mousePos: Vector2 = vec(0, 0);

  private eventAbort = new AbortController();
  private renderers: IRenderer[] = [];

  public get screenTopLeftInWorldPos(): Vector2 {
    return this.cameraPos
      .minus(vec(this.canvas.width, this.canvas.height).times(0.5).times(1/this.cameraScale));
  }

  public get screenBottomRightInWorldPos(): Vector2 {
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
    for (const r of this.renderers)
      r.cleanup();
  }

  public constructor(canvas: HTMLCanvasElement, api: Api) {
    this.canvas = canvas;
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      alert("device/browser not supported");
      throw new Error("could not create context");
    }
    this.ctx = ctx;

    // order is significant: render order
    this.renderers.push(new TerrainRenderer(this, api));
    this.renderers.push(new UnitRenderer(this, api));
    this.renderers.push(new OverlayRenderer(this, api));
    
    let clickMouseStart: Vector2 | null = null;
    let clickCameraStart: Vector2 | null = null;
    this.canvas.addEventListener("mousedown", ev => {
      const tile = this.mouseWorldPos.mapped(Math.floor);
      console.log("Cliked tile: ", {
        position: [this.mouseWorldPos.x, this.mouseWorldPos.y].join(" "),
        flooredPosition: [tile.x, tile.y].join(" "),
        kind: api.tilemap.getTileKind(tile),
        value: api.tilemap.getTileValue(tile),
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
  }

  private update(_dt: number) {
    
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

    for (const renderer of this.renderers)
      renderer.render(dt);
  }
}
