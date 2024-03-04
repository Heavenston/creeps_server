import { vec } from "~/src/utils/geom"
import { IRenderer, Renderer } from "./worldRenderer";
import { Api, PlayerSpawnMessage } from "./api";

export class OverlayRenderer implements IRenderer {
  private readonly renderer: Renderer;

  private eventAbort = new AbortController();
  private players = new Map<string, PlayerSpawnMessage>();

  private renderPlayerUsernames = true;

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(renderer: Renderer, api: Api) {
    this.renderer = renderer;

    document.body.addEventListener("keypress", event => {
      if (event.key == "p") {
        this.renderPlayerUsernames = !this.renderPlayerUsernames;
      }
    })

    api.addEventListener("message", event => {
      const message = event.message;
      if (message.kind != "playerSpawn")
        return;
      this.players.set(message.content.id, message);
    }, {
      signal: this.eventAbort.signal,
    });

    api.addEventListener("message", event => {
      const message = event.message;
      if (message.kind != "playerDespawn")
        return;
      this.players.delete(message.content.id);
    }, {
      signal: this.eventAbort.signal,
    });
  }

  private update(_dt: number) {
  }

  public render(dt: number) {
    if (dt != 0)
      this.update(dt);

    const ctx = this.renderer.ctx;

    if (!this.renderPlayerUsernames) {
      return;
    }

    for (const player of this.players.values()) {
      const sp = vec(player.content.spawnPosition).plus(0.5);

      ctx.textAlign = "center";
      ctx.textBaseline = "middle";

      ctx.strokeStyle = "rgba(0, 0, 0, 0.75)";
      ctx.lineWidth = 3 / this.renderer.cameraScale;
      ctx.strokeText(player.content.username, sp.x, sp.y);
      ctx.font = `${18 / this.renderer.cameraScale}px arial`;
      ctx.fillStyle = "rgba(255, 255, 255, 0.75)";
      ctx.fillText(player.content.username, sp.x, sp.y);
    }
  }
}
