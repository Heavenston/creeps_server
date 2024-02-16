import { MinzeElement } from "minze"

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
  private ctx: CanvasRenderingContext2D | null = null;
  private animationFrameId: number = -1;

  private resizeCanvas() {
    if (this.canvas == null)
      return;
    this.canvas.width = this.clientWidth;
    this.canvas.height = this.clientHeight;
    this.renderCanvas(this.lastTime);
  }

  private lastTime = 0;

  private x: number = 0;
  private y: number = 0;
  private dx: number = 0.5;
  private dy: number = 0.5;
  private readonly w: number = 50;
  private readonly h: number = 50;

  private renderCanvas(time: number) {
    if (this.canvas == null || this.ctx == null)
      return;
    const dt = time - this.lastTime;
    this.lastTime = time;

    this.x += this.dx * dt;
    if (this.x < 0) {
      this.x = 0;
      this.dx *= -1;
    }
    if (this.x+this.w > this.canvas.width) {
      this.x = this.canvas.width-this.w;
      this.dx *= -1;
    }
    this.y += this.dy * dt;
    if (this.y < 0) {
      this.y = 0;
      this.dy *= -1;
    }
    if (this.y+this.h > this.canvas.height) {
      this.y = this.canvas.height-this.h;
      this.dy *= -1;
    }

    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    this.ctx.fillStyle = "white";
    this.ctx.fillRect(this.x, this.y, this.w, this.h);

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
    this.ctx = ctx;
    this.resizeCanvas();

    new ResizeObserver(() => {
      this.resizeCanvas();
    }).observe(this);
  }

  onDestroy() {
    cancelAnimationFrame(this.animationFrameId);
  }
}).define("creeps-canvas")
