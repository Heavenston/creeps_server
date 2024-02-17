import { MinzeElement } from "minze"
import { vec, Vector2 } from "./geom"

class WorldRenderer {
  private readonly canvas: HTMLCanvasElement;
  private readonly ctx: CanvasRenderingContext2D;

  // position of the center of the screen in world coordinate
  public cameraPos: Vector2 = vec(0, 0);
  // scale to go from screen pos to world pos
  public cameraScale: number = 1;

  // position of the mouse in screen coordinated
  private mousePos: Vector2 = vec(0, 0);

  private eventAbort = new AbortController();

  private pos = vec(0, 0);
  private vel = vec(0.5, 0.5);
  private readonly size = vec(50, 50);
  private minpos = vec(-400, -300);
  private maxpos = vec( 400,  300);

  // changes the scale but also changes the cameraPos making sure the mousePos
  // doesn't change what it is pointing at
  private changeScale(val: number) {
    const adjustedPos = this.mousePos.minus(vec(this.canvas.width, this.canvas.height).times(0.5));
    const prevGobal = adjustedPos.times(1/this.cameraScale).plus(this.cameraPos);
    const newGlobal = adjustedPos.times(1/val).plus(this.cameraPos);

    this.cameraPos.sub(newGlobal.minus(prevGobal));
    console.log(prevGobal, newGlobal);
    this.cameraScale = val;
  }

  private update(dt: number) {
    // this.cameraPos = vec(this.canvas.width, this.canvas.height).times(0.5);

    this.pos.add(this.vel.times(dt));
    if (this.pos.x < this.minpos.x) {
      this.pos.x = this.minpos.x;
      this.vel.x *= -1;
    }
    if (this.pos.x+this.size.x > this.maxpos.x) {
      this.pos.x = this.maxpos.x-this.size.x;
      this.vel.x *= -1;
    }
    if (this.pos.y < this.minpos.y) {
      this.pos.y = this.minpos.y;
      this.vel.y *= -1;
    }
    if (this.pos.y+this.size.y > this.maxpos.y) {
      this.pos.y = this.maxpos.y-this.size.y;
      this.vel.y *= -1;
    }
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

    this.canvas.addEventListener("wheel", e => {
      const sign = e.deltaY < 0 ? -1 : 1;
      if (sign > 0)
        this.changeScale(this.cameraScale * 0.8);
      else
        this.changeScale(this.cameraScale * 1.2);
    }, {
      signal: this.eventAbort.signal,
    });

    this.ctx = ctx;
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

    this.ctx.beginPath();
    this.ctx.moveTo(this.minpos.x, this.minpos.y);
    this.ctx.lineTo(this.maxpos.x, this.minpos.y);
    this.ctx.lineTo(this.maxpos.x, this.maxpos.y);
    this.ctx.lineTo(this.minpos.x, this.maxpos.y);
    this.ctx.lineTo(this.minpos.x, this.minpos.y);
    this.ctx.lineTo(this.maxpos.x, this.minpos.y);
    this.ctx.strokeStyle = "white";
    this.ctx.lineWidth = 10;
    this.ctx.stroke();

    this.ctx.fillStyle = "white";
    this.ctx.fillRect(this.pos.x, this.pos.y, this.size.x, this.size.y);
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

  private resizeCanvas() {
    if (this.canvas == null)
      return;
    this.canvas.width = this.clientWidth;
    this.canvas.height = this.clientHeight;

    this.renderCanvas(this.lastTime);
  }

  private renderCanvas(time: number) {
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
    this.renderer = new WorldRenderer(this.canvas);
    this.resizeCanvas();

    new ResizeObserver(() => {
      this.resizeCanvas();
    }).observe(this);
  }

  onDestroy() {
    cancelAnimationFrame(this.animationFrameId);
  }
}).define("creeps-canvas")
