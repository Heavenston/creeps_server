import { Vector2 } from "~src/utils/geom";

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
    [
      "/viewer/kenney_micro_roguelike/grass1.png",
      "/viewer/kenney_micro_roguelike/grass2.png",
      "/viewer/kenney_micro_roguelike/grass3.png",
      "/viewer/kenney_micro_roguelike/grass4.png",
      "/viewer/kenney_micro_roguelike/grass5.png",
    ],
    // Water
    "/viewer/kenney_micro_roguelike/water.png",
    // Stone
    [
      "/viewer/tinyranch/rock1.png",
      "/viewer/tinyranch/rock2.png",
      "/viewer/tinyranch/rock3.png",
    ],
    // Tree
    [
      "/viewer/kenney_micro_roguelike/tree1.png",
      "/viewer/kenney_micro_roguelike/tree2.png",
    ],
    // Bush
    [
      "/viewer/bush.png",
      "/viewer/bush2.png",
    ],
    // Oil
    "/viewer/kenney_micro_roguelike/oil.png",
    // TownHall
    "/viewer/kenney_micro_roguelike/castle.png",
    // Household
    "/viewer/kenney_micro_roguelike/house.png",
    // Smeltery
    "/viewer/smeltery.png",
    // SawMill
    "/viewer/sawmill.png",
    // RaiderCamp
    [
      "/viewer/raidercamp1.png",
      "/viewer/raidercamp2.png",
    ],
    // RaiderBorder
    null,
    // Road
    "/viewer/road.png",
  ];
  private unitsUrlTable: {[key: string]: string|string[]|undefined} = {
    "citizen": [
      "/viewer/kenney_micro_roguelike/citizen_basic.png",
      "/viewer/kenney_micro_roguelike/citizen_woman.png",
    ],
    "citizen:upgraded": [
      "/viewer/kenney_micro_roguelike/citizen_adventurer.png",
      "/viewer/kenney_micro_roguelike/citizen_woman2.png",
    ],
    "turret": "/viewer/kenney_micro_roguelike/red_robot.png",
    "raider": "/viewer/kenney_micro_roguelike/zombie.png",
    "bomber-bot": [
      "/viewer/bomberbot1.png",
      "/viewer/bomberbot2.png",
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

    return this.loadingTexture;
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

  public getUnitTexture(opcode: string, unitId: string, upgraded: boolean): ImageBitmap {
    let url = this.unitsUrlTable[opcode];
    if (upgraded && this.unitsUrlTable[opcode + ":upgraded"])
      url = this.unitsUrlTable[opcode + ":upgraded"];
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
