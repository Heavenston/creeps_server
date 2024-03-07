import { MinzeElement, Reactive } from "minze";
import * as mapi from "~/src/manager_api"
import "./dashboard"
import "../viewer/canvas"
import { createPopup } from "~src/popup";

export interface GameComp {
  game: mapi.Game | null;
  panelOpened: boolean;
}

export class GameComp extends MinzeElement {
  reactive: Reactive = [["game", null], ["panelOpened", false]];
  
  html = () => `
    ${this.game != null && this.game.viewer_port ? `
      <creeps-canvas url="ws://localhost:${this.game.viewer_port}/websocket">
      </creeps-canvas>
      <div class="panel ${this.panelOpened ? "open" : "close"}">
        <button class="closePanel" on:click="handleTogglePanel">
          <span> ${this.panelOpened ? '>' : "<"} </span>
        </button>
        <div class="internal">
        </div>
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

    background-color: rgba(255, 255, 255, 0.07);
    background-color: #212121;

    display: flex;
    flex-direction: row;
  }

  .panel .internal {
    overflow: hidden;
    transition: width 150ms;
  }

  .panel.close .internal {
    width: 0rem;
  }

  .panel.open .internal {
    width: 30rem;
  }

  .closePanel {
    cursor: pointer;
    width: 2rem;

    display: flex;
    justify-content: center;
    align-items: center;
  }

  .closePanel > span {
    transform: scaleY(10);
  }

  .internal {
    flex-grow: 1;
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

  handleTogglePanel() {
    this.panelOpened = !this.panelOpened;
  }
}

GameComp.define("creeps-game");

