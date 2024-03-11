import { MinzeElement, Reactive } from "minze";

export interface SidePanelComp {
  panelOpened: boolean;
}

export class SidePanelComp extends MinzeElement {
  reactive: Reactive = [["panelOpened", true]];

  html = () => `
    <button class="closePanel" on:click="handleTogglePanel">
      <span> ${this.panelOpened ? '>' : "<"} </span>
    </button>
    <div class="internal">
      <slot></slot>
    </div>
  `;

  css = () => `
  :host {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 0;

    background-color: var(--dark-two);

    display: flex;
    flex-direction: row;
  }

  .internal {
    overflow: hidden;
    transition: width 150ms;
  }

  :host(.closed) .internal {
    width: 0;
  }

  :host(.opened) .internal {
    width: 500px;
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
  `;

  afterRender() {
    if (this.panelOpened) {
      this.classList.add("opened");
      this.classList.remove("closed");
    } else {
      this.classList.remove("opened");
      this.classList.add("closed");
    }
  }

  handleTogglePanel() {
    this.panelOpened = !this.panelOpened;
  }
}

SidePanelComp.define("side-panel");
