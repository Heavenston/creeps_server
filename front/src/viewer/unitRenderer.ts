import { vec, Vector2 } from "~/src/geom"
import * as api from "~/src/api"
import { IRenderer, Renderer } from "./worldRenderer";

type RunningAction = {
  action: api.Action,
  elapsed: number,
} & ({
  state: "running",
} | {
  state: "finished" | "error",
  finishedSince: number,
});

export class UnitRenderer implements IRenderer {
  private readonly renderer: Renderer;

  private eventAbort = new AbortController();

  private lastUnitMessage: Map<string, api.UnitMessage> = new Map();
  private unitsActions: Map<string, RunningAction> = new Map();

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(renderer: Renderer) {
    this.renderer = renderer;

    api.addEventListener("message", event => {
      switch (event.message.kind) {
        case "unit": {
          this.lastUnitMessage.set(event.message.content.unitId, event.message);
          break;
        }
        case "unitDespawned": {
          this.unitsActions.delete(event.message.content.unitId);
          this.lastUnitMessage.delete(event.message.content.unitId);
          break;
        }
        case "unitMovement": {
          console.log("movement")
          const unit = this.lastUnitMessage.get(event.message.content.unitId);
          if (!unit) {
            console.warn("received unit movement for unkown unit ", event.message);
            break;
          }
          unit.content.position = event.message.content.new;
          break;
        }
        case "unitStartedAction": {
          if (event.message.content.action.actionOpCode.startsWith("move:"))
            console.log("move actino", event.message.content)
          this.unitsActions.set(event.message.content.unitId, {
            action: event.message.content.action,
            elapsed: 0,
            state: "running",
          });
          break;
        }
        case "unitFinishedAction": {
          if (event.message.content.action.actionOpCode.startsWith("move:"))
            console.log("finished move action ", event.message.content)
          const act = this.unitsActions.get(event.message.content.unitId);
          if (!act) {
            console.warn("received finished action for unkown action", event.message);
            break;
          }
          this.unitsActions.set(event.message.content.unitId, {
            action: event.message.content.action,
            elapsed: act.elapsed,
            state: event.message.content.success ? "finished" : "error",
            finishedSince: 0,
          })
          break;
        }
      }
    }, {
      signal: this.eventAbort.signal,
    });
  }

  private update(dt: number) {
    for (const unitId of this.lastUnitMessage.keys()) {
      const unit = this.lastUnitMessage.get(unitId);
      if (!unit)
        continue;
      const action = this.unitsActions.get(unitId);
      if (!action)
        continue;

      action.elapsed += dt;
    }
  }

  private renderUnit(unit: api.UnitMessage) {
    let pos = vec(unit.content.position);
    const action = this.unitsActions.get(unit.content.unitId);
    const cost = action == null ? null : api.getActionCost(action.action.actionOpCode);
    
    if (action != null && cost != null && action.state == "running") {
      let prop = action.elapsed / (cost.cast * api.getSecondsPerTicks());
      prop = Math.min(Math.max(prop, 0), 1);
      switch (action.action.actionOpCode) {
        case "move:left":
          pos.x -= prop;
          break;
        case "move:right":
          pos.x += prop;
          break;
        case "move:up":
          pos.y += prop;
          break;
        case "move:down":
          pos.y -= prop;
          break;
      }
    }


    const texture = this.renderer.texturePack.getUnitTexture(
      unit.content.opCode,
      unit.content.unitId,
      unit.content.upgraded,
    );
    this.renderer.ctx.drawImage(texture, pos.x, pos.y, 1, 1);
  }

  public render(dt: number) {
    if (dt != 0)
      this.update(dt);

    for (const unit of this.lastUnitMessage.values())
      if (unit.content.opCode == "turret")
        this.renderUnit(unit);
    for (const unit of this.lastUnitMessage.values())
      if (unit.content.opCode != "turret")
        this.renderUnit(unit);
  }
}

