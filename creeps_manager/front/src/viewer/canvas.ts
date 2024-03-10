import { Attrs, MinzeElement, Reactive } from "minze"

import { Renderer } from "./worldRenderer"
import { Api } from "./api"

export interface CreepsCanvasComp {
  url: string | null;
  panelOpened: boolean;
}

export class CreepsCanvasComp extends MinzeElement {
  reactive: Reactive = [["panelOpened", false]];
  attrs: Attrs = [["url", null]]
  
  // html template
  html = () => `
    <canvas></canvas>
    <div class="panel ${this.panelOpened ? "open" : "close"}">
      <button class="closePanel" on:click="handleTogglePanel">
        <span> ${this.panelOpened ? '>' : "<"} </span>
      </button>
      <div class="internal">
      </div>
    </div>
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

  .panel {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 0;

    background-color: rgba(255, 255, 255, 0.07);
    background-color: #212121;

    display: flex;
    flex-direction: row;
  }

  .panel .internal {
    overflow: hidden;
    transition: width 150ms;
  }

  .panel.close .internal {
    width: 0rem;
  }

  .panel.open .internal {
    width: 30rem;
  }

  .closePanel {
    cursor: pointer;
    width: 2rem;

    display: flex;
    justify-content: center;
    align-items: center;
  }

  .closePanel > span {
    transform: scaleY(10);
  }

  .internal {
    flex-grow: 1;
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

  handleTogglePanel() {
    this.panelOpened = !this.panelOpened;
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
