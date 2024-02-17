
export function vec(other: Vector2): Vector2;
export function vec(x: number): Vector2;
export function vec(x: number, y: number): Vector2;

export function vec(other: Vector2 | number, other2?: number): Vector2 {
  if (other instanceof Vector2)
    return new Vector2(other.x, other.y);
  if (other2 == undefined)
    return new Vector2(other, other);
  return new Vector2(other, other2);
}

export class Vector2 {
  public constructor(public x: number, public y: number) {
  }

  public add(other: Vector2): void;
  public add(x: number): void;
  public add(x: number, y: number): void;

  public add(other: Vector2 | number, other2?: number): void {
    if (other instanceof Vector2) {
      if (other2 != undefined)
        throw new TypeError("Invalid input");
      this.x += other.x;
      this.y += other.y;
    }
    else if (other2 == undefined) {
      this.add(new Vector2(other, other));
    }
    else {
      this.add(new Vector2(other, other2));
    }
  }

  public plus(other: Vector2): Vector2;
  public plus(x: number): Vector2;
  public plus(x: number, y: number): Vector2;

  public plus(other: Vector2 | number, other2?: number): Vector2 {
    const nv = new Vector2(this.x, this.y);
    if (other instanceof Vector2)
      nv.add(other);
    else if (other2 == undefined)
      nv.add(other);
    else
      nv.add(other, other2);
    return nv;
  }

  public sub(other: Vector2): void;
  public sub(x: number): void;
  public sub(x: number, y: number): void;

  public sub(other: Vector2 | number, other2?: number): void {
    if (other instanceof Vector2) {
      if (other2 != undefined)
        throw new TypeError("Invalid input");
      this.x -= other.x;
      this.y -= other.y;
    }
    else if (other2 == undefined) {
      this.sub(new Vector2(other, other));
    }
    else {
      this.sub(new Vector2(other, other2));
    }
  }

  public minus(other: Vector2): Vector2;
  public minus(x: number): Vector2;
  public minus(x: number, y: number): Vector2;

  public minus(other: Vector2 | number, other2?: number): Vector2 {
    const nv = new Vector2(this.x, this.y);
    if (other instanceof Vector2)
      nv.sub(other);
    else if (other2 == undefined)
      nv.sub(other);
    else
      nv.sub(other, other2);
    return nv;
  }

  public mul(other: Vector2): void;
  public mul(x: number): void;
  public mul(x: number, y: number): void;

  public mul(other: Vector2 | number, other2?: number): void {
    if (other instanceof Vector2) {
      if (other2 != undefined)
        throw new TypeError("Invalid input");
      this.x *= other.x;
      this.y *= other.y;
    }
    else if (other2 == undefined) {
      this.mul(new Vector2(other, other));
    }
    else {
      this.mul(new Vector2(other, other2));
    }
  }

  public times(other: Vector2): Vector2;
  public times(x: number): Vector2;
  public times(x: number, y: number): Vector2;

  public times(other: Vector2 | number, other2?: number): Vector2 {
    const nv = new Vector2(this.x, this.y);
    if (other instanceof Vector2)
      nv.mul(other);
    else if (other2 == undefined)
      nv.mul(other);
    else
      nv.mul(other, other2);
    return nv;
  }
}
