const API_BASE_URL = process.env.MANAGER_API_BASE_URL ?? "http://localhost:16969/api";
const LOGIN_URl = process.env.LOGIN_URL ?? "http://example.com";

let loginPromise: Promise<void> | null = null;

export type User = {
  id: number,
  discord_id: string,
  discord_tag: string,
  avatar_url: string,
  username: string,
};

export async function logout() {
  localStorage.removeItem("token");
  setTimeout(() => document.location.reload());
  await new Promise(() => {})
}

export function isLoggedIn(): boolean {
  try {
    const usp = new URLSearchParams(document.location.search);
    const token = usp.get("token");
    if (token != null)
    {
      localStorage.setItem("token", token);
      document.location.search = "";
    }
  }
  catch (e) { }

  return localStorage.getItem("token") != null;
}

async function makeLogin() {
  try {
    const usp = new URLSearchParams(document.location.search);
    const token = usp.get("token");
    if (token != null)
      localStorage.setItem("token", token);
  }
  catch (e) { }

  const token = localStorage.getItem("token");
  if (token == null) {
    document.location.href = LOGIN_URl;
    await new Promise(() => {});
  }
}

export function login(): Promise<void> {
  if (loginPromise != null)
    return loginPromise;
  loginPromise = makeLogin();
  return loginPromise;
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
    await logout();
  }
  if (!resp.ok) {
    throw new Error("req error");
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
