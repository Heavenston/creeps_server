import { MinzeElement } from "minze";
import "./dashboard"
import "../viewer/canvas"

(class GameComp extends MinzeElement {
  html = () => `
    <creeps-canvas/>
  `;

  css = () => `
    :host {
      flex-grow: 1;
      position: relative;
    }
  `
}).define("creeps-game")

