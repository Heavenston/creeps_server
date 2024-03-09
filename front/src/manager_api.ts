(window as any).process = { env: {} };
const API_BASE_URL = process.env.MANAGER_API_BASE_URL ?? "http://localhost:16969/api";
const LOGIN_URL = process.env.LOGIN_URL ?? "http://example.com";

let loginPromise: Promise<void> | null = null;

export type User = {
  id: number,
  discord_id: string,
  discord_tag: string,
  avatar_url: string,
  username: string,
};

export type GameConfig = {
  can_join_after_start: boolean,
  private: boolean,
  is_local: boolean,
};

export type Game = {
  id: number,
  name: string,

  creator: User,
  players: User[],

  config: GameConfig,

  started_at?: number,
  ended_at?: number,
  
  api_port?: number,
  viewer_port?: number,
};

export async function logout() {
  localStorage.removeItem("token");
  setTimeout(() => document.location.reload());
  await new Promise(() => { })
}

export function isLoggedIn(): boolean {
  try {
    const usp = new URLSearchParams(document.location.search);
    const token = usp.get("token");
    if (token != null) {
      localStorage.setItem("token", token);
      document.location.search = "";
    }
  }
  catch (e) { }

  return localStorage.getItem("token") != null;
}

async function makeLogin() {
  if (isLoggedIn())
    return;

  const token = localStorage.getItem("token");
  if (token == null) {
    document.location.href = LOGIN_URL;
    await new Promise(() => { });
  }
}

export function login(): Promise<void> {
  if (loginPromise != null)
    return loginPromise;
  loginPromise = makeLogin();
  return loginPromise;
}

export class RequestError extends Error {
  constructor(public readonly response: Response) {
    super("request error");
  }
}

async function get<T>(url: string): Promise<T> {
  const headers = new Headers();
  if (isLoggedIn()) {
    headers.set("Authorization", localStorage.getItem("token") ?? "");
  }
  const resp = await fetch(API_BASE_URL + url, {
    method: "GET",
    headers,
  });
  if (resp.status == 403) {
    if (isLoggedIn())
      await logout();
  }
  if (!resp.ok) {
    throw new RequestError(resp);
  }
  return resp.json();
}

async function post<T>(url: string, body: any): Promise<T> {
  const headers = new Headers();
  headers.set("content-type", "application/json");
  if (isLoggedIn()) {
    headers.set("Authorization", localStorage.getItem("token") ?? "");
  }
  const resp = await fetch(API_BASE_URL + url, {
    method: "POST",
    body: JSON.stringify(body),
    headers,
  });
  if (resp.status == 403) {
    if (isLoggedIn())
      await logout();
  }
  if (!resp.ok) {
    throw new RequestError(resp);
  }
  return resp.json();
}

export async function getUserSelf(): Promise<User> {
  if (!isLoggedIn())
    throw new Error("cannot get user if not logged in");

  if (sessionStorage.getItem("user")) {
    return JSON.parse(sessionStorage.getItem("user") ?? "");
  }

  const m: User = await get("/users/@me");
  sessionStorage.setItem("user", JSON.stringify(m));
  return m;
}

export async function getGame(id: string): Promise<Game> {
  return get(`/games/${id}`);
}

export async function getGames(): Promise<Game[]> {
  return get("/games");
}

export async function createGame(name: String): Promise<Game> {
  return post("/games", { name });
}
