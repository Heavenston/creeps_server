import { Attrs, MinzeElement } from "minze"

import { Renderer } from "./worldRenderer"
import { Api } from "./api"

export interface CreepsCanvasComp {
  url: string | null;
}

export class CreepsCanvasComp extends MinzeElement {
  attrs: Attrs = [["url", null]]
  
  // html template
  html = () => `
    <canvas></canvas>
  `

  // scoped stylesheet
  css = () => `
  :host {
    flex-grow: 1;
    position: relative;
  }

  canvas {
    position: absolute;
    inset: 0;
  }
  `

  private canvas: HTMLCanvasElement | null = null;
  private worldRenderer: Renderer | null = null;
  private animationFrameId: number = -1;
  private api: Api | null = null;

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
    if (this.worldRenderer == null)
      return;
    this.animationFrameId = -1;

    this.worldRenderer?.render((time - this.lastTime) / 1000);
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

    this.api = new Api(this.url ?? "");

    this.api.addEventListener("connection_event", c => {
      if (!this.api)
        return;

      if (c.isConnected) {
        if (this.canvas != null)
          this.worldRenderer = new Renderer(this.canvas, this.api);
        this.resizeCanvas();
      }
      else {
        this.worldRenderer?.cleanup();
        this.worldRenderer = null;
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
    if (this.worldRenderer)
      this.worldRenderer.cleanup();
    cancelAnimationFrame(this.animationFrameId);
    this.eventAbort.abort();
  }
}

CreepsCanvasComp.define("creeps-canvas")
