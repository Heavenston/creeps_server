import { vec, Vector2 } from "~/src/geom"
import * as api from "~/src/api"
import { IRenderer, Renderer } from "./worldRenderer";

export class UnitRenderer implements IRenderer {
  private readonly renderer: Renderer;

  private eventAbort = new AbortController();

  private lastUnitMessage: Map<string, api.UnitMessage> = new Map();
  private unitsPositions: Map<string, Vector2> = new Map();

  public cleanup() {
    this.eventAbort.abort();
  }

  public constructor(renderer: Renderer) {
    this.renderer = renderer;

    api.addEventListener("message", event => {
      if (event.message.kind != "unit")      
        return;
      this.lastUnitMessage.set(event.message.content.unitId, event.message);
    }, {
      signal: this.eventAbort.signal,
    });

    api.addEventListener("message", event => {
      if (event.message.kind != "unitMovement")      
        return;
      const unit = this.lastUnitMessage.get(event.message.content.unitId);
      if (!unit) {
        console.warn("received unit movement for unkown unit ", event.message);
        return;
      }
      unit.content.position = event.message.content.new;
    }, {
      signal: this.eventAbort.signal,
    });

    api.addEventListener("message", event => {
      if (event.message.kind != "unitDespawned")      
        return;
      this.unitsPositions.delete(event.message.content.unitId);
      this.lastUnitMessage.delete(event.message.content.unitId);
    }, {
      signal: this.eventAbort.signal,
    });
  }

  private update(dt: number) {
    for (const unitId of this.lastUnitMessage.keys()) {
      const unit = this.lastUnitMessage.get(unitId);
      if (!unit)
        continue;
      let pos = this.unitsPositions.get(unitId);
      if (!pos) {
        pos = vec(unit.content.position);
        this.unitsPositions.set(unitId, pos);
      }

      pos.lerp(40 * dt, vec(unit.content.position));
    }
  }

  private renderUnit(unit: api.UnitMessage) {
    const pos = this.unitsPositions.get(unit.content.unitId)
      ?? vec(unit.content.position);

    const texture = this.renderer.texturePack.getUnitTexture(unit.content.opCode, unit.content.unitId);
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

