import { MinzeElement } from "minze"
import * as api from "~/src/api"

import { Renderer } from "./worldRenderer"

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
  private worldRenderer: Renderer | null = null;
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

    api.addEventListener("connection_event", c => {
      if (c.isConnected) {
        if (this.canvas != null)
          this.worldRenderer = new Renderer(this.canvas);
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
}).define("creeps-canvas")
