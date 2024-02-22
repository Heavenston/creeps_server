import { Vector2 } from "~src/geom";

function splitmix32(a: number) {
    a |= 0; a = a + 0x9e3779b9 | 0;
    let t = a ^ a >>> 16; t = Math.imul(t, 0x21f0aaad);
        t = t ^ t >>> 15; t = Math.imul(t, 0x735a2d97);
    return (t = t ^ t >>> 15) >>> 0;
}

export class TexturePack extends EventTarget {
  public readonly size: number;
  // color used as background of for all tiles
  public readonly fillColor: string = "#2c9c50";

  private defaultTexture: ImageBitmap;
  private loadingTexture: ImageBitmap;

  private tilesUrlTable: (string|string[]|null)[] = [
    // Grass
    "/kenney_micro_roguelike/grass1.png",
    // Water
    "/kenney_micro_roguelike/water.png",
    // Stone
    [
      "/tinyranch/rock1.png",
      "/tinyranch/rock2.png",
      "/tinyranch/rock3.png",
    ],
    // Tree
    [
      "/kenney_micro_roguelike/tree1.png",
      "/kenney_micro_roguelike/tree2.png",
    ],
    // Bush
    [
      "/bush.png",
      "/bush2.png",
    ],
    // Oil
    "/kenney_micro_roguelike/oil.png",
    // TownHall
    "/kenney_micro_roguelike/castle.png",
    // Household
    "/kenney_micro_roguelike/house.png",
    // Smeltery
    null,
    // SawMill
    null,
    // RaiderCamp
    null,
    // RaiderBorder
    null,
    // Road
    "/kenney_micro_roguelike/road.png",
  ];
  private unitsUrlTable: {[key: string]: string|string[]|undefined} = {
    "citizen": "/kenney_micro_roguelike/citizen_basic.png",
    "turret": "/kenney_micro_roguelike/red_robot.png",
    "raider": "/kenney_micro_roguelike/zombie.png",
    "bomber-bot": undefined,
  };
  private textureCache = new Map<string, "loading" | ImageBitmap>();

  public constructor() {
    super();

    this.size = 8;

    const canvas = new OffscreenCanvas(8, 8);
    const ctx = canvas.getContext("2d");
    if (!ctx)
      throw new Error("unsupported");

    ctx.fillStyle = "magenta";
    ctx.fillRect(0, 0, 4, 4);
    ctx.fillRect(4, 4, 4, 4);
    ctx.fillStyle = "black";
    ctx.fillRect(0, 4, 4, 4);
    ctx.fillRect(4, 0, 4, 4);

    this.defaultTexture = canvas.transferToImageBitmap();

    ctx.clearRect(0, 0, 8, 8);

    this.loadingTexture = canvas.transferToImageBitmap();
  }

  private getTexture(url: string): ImageBitmap {
    const cached = this.textureCache.get(url);
    if (cached instanceof ImageBitmap)
      return cached;
    if (cached == "loading")
      return this.loadingTexture;

    console.log("loading", url);
    this.textureCache.set(url, "loading");
    const image = new Image();
    image.addEventListener("load", () => {
      createImageBitmap(image).then(i => {
        this.textureCache.set(url, i);
        this.dispatchEvent(new Event("textureLoaded"));
      }).catch(e => {
        console.error(e);
      });
    });
    image.addEventListener("error", e => {
      console.error(`could not load image ${url}`, e);
    });
    image.src = url;

    return this.defaultTexture;
  }

  public getTileTexture(tileKind: number, tilePos: Vector2): ImageBitmap {
    const url = this.tilesUrlTable[tileKind];
    if (!url)
      return this.defaultTexture;
    let realUrl: string;
    if (Array.isArray(url)) {
      const a = splitmix32(tilePos.x);
      const b = splitmix32(tilePos.y);
      realUrl = url[Math.abs(a ^ b ^ 951274213) % url.length];
    }
    else {
      realUrl = url;
    }

    return this.getTexture(realUrl);
  }

  public getUnitTexture(opcode: string, unitId: string): ImageBitmap {
    const url = this.unitsUrlTable[opcode];
    if (!url)
      return this.defaultTexture;
    let realUrl: string;
    if (Array.isArray(url)) {
      let sum = 0;
      for (let i = 0; i < unitId.length; i++)
        sum += unitId.charCodeAt(i);
      realUrl = url[sum % url.length];
    }
    else
      realUrl = url;

    return this.getTexture(realUrl);
  }
}
