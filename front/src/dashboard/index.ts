import { MinzeElement } from "minze";
import "./dashboard"
import * as mapi from "~/src/manager_api"
import { createPopup } from "~/src/popup";

(class IndexComp extends MinzeElement {
  #games: mapi.Game[] = []

  attrs = [["no-css-reset", ""]] as any;

  html = () => `
    <section>
      <h1><span>Active games</span><span class="line"></span></h1>
      <div class="content">
        ${this.#games.map(game => `
          <button class="item item-game" on:click="join">
            <div>${game.name}</div>
            <div class="players">${game.players.length} players</div>
          </button>
        `).join("")}
        ${mapi.isLoggedIn() ? `
          <button class="item item-new" on:click="handleCreate">
            <span>+</span>
          </button>
        ` : ``}
      </div>
    </section>
    <section>
      <h1><span>Past games</span><span class="line"></span></h1>
      <div class="content">
      </div>
    </section>
  `

  css = () => `
    :host {
      flex-grow: 1;

      display: flex;
      height: 100%;
      flex-direction: column;

      padding: 1.2rem !important;
      gap: 1rem;
    }

    h1 {
      display: flex;
      font-size: 1.2rem;

      flex-direction: row;

      align-items: center;

      gap: 1rem;
    }

    h1 .line {
      flex-grow: 1;
      height: 0.2rem;

      background: rgba(255, 255, 255, 0.4);
    }

    section .content {
      display: flex;
      flex-direction: row;
      gap: 1rem;

      padding: 1rem;
    }

    .item {
      width: 12rem;
      height: 8rem;
      background: rgba(255, 255, 255, 0.07);
      border-radius: 0.3rem;

      display: flex;
      align-items: center;
      justify-content: center;

      transition: background 100ms;
    }

    .item:hover {
      background: rgba(255, 255, 255, 0.2);
    }

    .item-new {
      font-size: 2rem;
    }

    .item-game {
      position: relative;
    }

    .item-game .players {
      position: absolute;

      bottom: 0.5rem;
      right: 0.5rem;

      color: #777;
    }
  `

  updateGames(games: mapi.Game[]) {
    this.#games = games;
    this.rerender();
  }

  onReady() {
    mapi.getGames().then(this.updateGames.bind(this));
  }

  join() {
    document.location.href = "/game";
  }

  handleCreate() {
    mapi.createGame("lol")
      .then(() => mapi.getGames())
      .then(this.updateGames.bind(this))
      .catch((e) => {
        if (e instanceof mapi.RequestError) {
          e.response.json().then(body => {
            createPopup("error", body["message"] ?? body["error"] ?? `Could not create game (${e.response.status})`)
          });
        }
      });
  }
}).define("creeps-index");
