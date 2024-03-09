import { MinzeElement } from "minze"
import * as mapi from "~/src/manager_api"

(class HeaderComp extends MinzeElement {
  #user: mapi.User | null = null;

  attrs = [["no-css-reset", ""]] as any;
  
  html = () => `
    <a class="start" href="/">
      <h1>Heav's Creeps</h1>
    </a>
    <div style="flex-grow: 1;"></div>
    ${mapi.isLoggedIn() && this.#user ? `
      <button class="end" on:click="handleLogout">
        <span class="username">${this.#user.username}</span>
        <span class="image" style="background-image: url(${this.#user.avatar_url});"></span>
      </button>
    ` : `
      <button class="end" on:click="handleLogin">
        Login
      </button>
    `}
  `;

  css = () => `
    :host {
      background-color: rgba(255, 255, 255, 0.07);

      display: flex;
      align-items: center;

      vertical-align: middle;
    }

    :host h1 {
      font-size: 1.2rem;
    }

    :host>* {
      display: flex;
      align-items: center;

      padding-left: 1rem;
      padding-right: 1rem;
      padding-top: 0.7rem;
      padding-bottom: 0.7rem;
    }

    .end {
      gap: 0.5rem;
      cursor: pointer;
    }

    .end:hover {
      background-color: rgba(255, 255, 255, 0.07);
    }

    .image {
      display: inline-block;
      background-size: contain;
      height: 2rem;
      width: 2rem;
      border-radius: 10%;
    }
  `;

  onReady() {
    if (mapi.isLoggedIn()) {
      mapi.getUserSelf().then(user => {
        this.#user = user;
        this.rerender();
      })
    }
  }

  handleLogin() {
    mapi.login();
  }
  handleLogout() {
    mapi.logout();
  }
}).define("creeps-header");

(class DashboardComp extends MinzeElement {
  html = () => `
    <creeps-header></creeps-header>
    <slot></slot>
  `

  css = () => `
    :host {
      display: flex;
      height: 100%;
      flex-direction: column;

      justify-content: flex-start;
      align-items: stretch;
    }
  `
}).define("creeps-dashboard");
