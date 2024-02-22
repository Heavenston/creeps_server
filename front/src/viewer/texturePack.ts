import { Vector2 } from "~src/geom";

export class TexturePack {
  public readonly size: number;
  // color used as background of for all tiles
  public readonly fillColor: string = "#2c9c50";

  private defaultTexture: ImageBitmap;
  private urlTable: (string|string[]|null)[] = [
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
  private textureCache = new Map<string, "loading" | ImageBitmap>();

  public constructor() {
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
  }

  public getTileTexture(tileKind: number, tilePos: Vector2): ImageBitmap {
    const url = this.urlTable[tileKind];
    if (!url)
      return this.defaultTexture;
    let realUrl: string;
    if (Array.isArray(url)) {
      realUrl = url[Math.abs(Math.abs(tilePos.x) ^ Math.abs(tilePos.y)) % url.length];
    }
    else {
      realUrl = url;
    }

    const cached = this.textureCache.get(realUrl);
    if (cached instanceof ImageBitmap)
      return cached;
    if (cached == "loading")
      return this.defaultTexture;

    console.log("loading", url);
    this.textureCache.set(realUrl, "loading");
    const image = new Image();
    image.addEventListener("load", () => {
      createImageBitmap(image).then(i => {
        this.textureCache.set(realUrl, i);
      }).catch(e => {
        console.error(e);
      });
    });
    image.addEventListener("error", e => {
      console.error(`could not load image ${realUrl}`, e);
    });
    image.src = realUrl;

    return this.defaultTexture;
  }
}
