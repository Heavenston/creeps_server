import { MinzeElement } from "minze"

(class HeaderComp extends MinzeElement {
  attrs = [["no-css-reset", ""]] as any;
  
  html = () => `
  <div class="start">
    <h1>Heav's Creeps</h1>
  </div>
  <div style="flex-grow: 1;"></div>
  <div class="end">
    <span class="username">Heavenstone</span>
    <span class="image"></span>
  </div>
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
      background-image: url(https://0.gravatar.com/avatar/568a630daeef28607806fc21b9344800cbaf11120d53d503adf7dbb25e56b79c?size=256);
      background-size: contain;
      height: 2rem;
      width: 2rem;
      border-radius: 10%;
    }
  `;
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
