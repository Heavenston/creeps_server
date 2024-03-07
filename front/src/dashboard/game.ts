import { MinzeElement, Reactive } from "minze";
import * as mapi from "~/src/manager_api"
import "./dashboard"
import "../viewer/canvas"
import { createPopup } from "~src/popup";

export interface GameComp {
  game: mapi.Game | null;
}

export class GameComp extends MinzeElement {
  reactive: Reactive = [["game", null]];
  
  html = () => `
    ${this.game != null && this.game.viewer_port ? `
      <creeps-canvas url="ws://localhost:${this.game.viewer_port}/websocket">
      </creeps-canvas>
    ` : ``}
  `;

  css = () => `
    :host {
      flex-grow: 1;
      position: relative;
    }
  `

  onReady() {
    const gameId = document.location.hash.slice(1);
    mapi.getGame(gameId).then(game => {
      this.game = game;
    }).catch(e => {
      if (e instanceof mapi.RequestError) {
        e.response.json().then(body => {
          createPopup("error", body["message"] ?? body["error"] ?? "An error occured");
        });
      }
    });
  }
}

GameComp.define("creeps-game");

