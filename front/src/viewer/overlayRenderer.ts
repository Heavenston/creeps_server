import { vec, Vector2 } from "~/src/geom"
import * as api from "~/src/api"
import * as map from "./map"
import { Renderer } from "./worldRenderer";

export class OverlayRenderer {
  private readonly renderer: Renderer;

  private eventAbort = new AbortController();
  private players = new Map<string, api.PlayerSpawnMessage>();

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(renderer: Renderer) {
    this.renderer = renderer;

    api.addEventListener("message", event => {
      const message = event.message;
      if (message.kind != "playerSpawn")
        return;
      this.players.set(message.content.id, message);
    });

    api.addEventListener("message", event => {
      const message = event.message;
      if (message.kind != "playerDespawn")
        return;
      this.players.delete(message.content.id);
    });
  }

  private update(_dt: number) {
  }

  public render(dt: number) {
    if (dt != 0)
      this.update(dt);

    const ctx = this.renderer.ctx;

    for (const player of this.players.values()) {
      const sp = player.content.spawnPosition;

      ctx.textAlign = "center";
      ctx.textBaseline = "middle";

      ctx.strokeStyle = "rgba(0, 0, 0, 1)";
      ctx.lineWidth = 3 / this.renderer.cameraScale;
      ctx.strokeText(player.content.username, sp.x, sp.y);
      ctx.font = `${18 / this.renderer.cameraScale}px arial`;
      ctx.fillStyle = "rgba(255, 255, 255, 1)";
      ctx.fillText(player.content.username, sp.x, sp.y);
    }
  }
}
