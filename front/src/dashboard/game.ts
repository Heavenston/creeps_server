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
      <div class="panel">
        <button class="close">
          X
        </button>
      </div>
    ` : ``}
  `;

  css = () => `
  :host {
    flex-grow: 1;
    position: relative;
  }

  .panel {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 0;
    width: 30rem;
    background-color: rgba(255, 255, 255, 0.07);
    background-color: #212121;
  }

  .close {
    display: block;

    position: absolute;
    top: 0;
    left: 0;

    height: 2rem;
    width: 2rem;

    background-color: #C23B22;
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
      else {
        createPopup("error", "An error occured");
      }
    });
  }
}

GameComp.define("creeps-game");

