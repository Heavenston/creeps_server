export const WEBSOCKET_URL: string =
  process.env.PREVIEW_WEBSOCKET_URL ?? "ws://localhost:1234/websocket";

const RETRY_INTERVAL: number = 5000;

let ws: WebSocket | null = null;

connect();
function connect() {
  if (ws != null)
    return;

  console.log("connecting to websocket");

  try {
    ws = new WebSocket(WEBSOCKET_URL);
  }
  catch(e) {
    ws = null;
    console.error("connect error", e, "retry in 5s")
    setTimeout(connect, RETRY_INTERVAL);
    return;
  }

  ws.addEventListener("open", (e) => {
    console.info("connected to websocket", e);
  });

  ws.addEventListener("message", (e) => {
    console.info("message", e);
  });

  ws.addEventListener("error", (e) => {
    if (ws == null)
      return;
    console.error("websocket error", e, "reconnecting retry in 5s");
    ws = null;
    setTimeout(connect, RETRY_INTERVAL);
  });
  ws.addEventListener("close", (e) => {
    if (ws == null)
      return;
    console.info("dirconnected from the websocket", e, "reconnecting in 5s");
    ws = null;
    setTimeout(connect, RETRY_INTERVAL);
  });
}
