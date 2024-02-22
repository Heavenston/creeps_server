import { Vector2 } from "~src/geom";

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
    "/tinyranch/rock1.png",
    // Tree
    [
      "/kenney_micro_roguelike/tree1.png",
      "/kenney_micro_roguelike/tree2.png",
    ],
    // Bush
    "/kenney_micro_roguelike/bush.png",
    // Oil
    null,
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
    null,
  ];
  private unitsUrlTable: {[key: string]: string|string[]|undefined} = {
    "citizen": [
      "/kenney_micro_roguelike/citizen_basic.png",
    ],
    "turret": [
      "/kenney_micro_roguelike/red_robot.png",
    ],
    "raider": [
      "/kenney_micro_roguelike/zombie.png",
    ],
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
      realUrl = url[Math.abs(Math.abs(tilePos.x) ^ Math.abs(tilePos.y)) % url.length];
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
