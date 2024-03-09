import { Attrs, MinzeElement } from "minze";

export interface PopupComp {
  kind: "info" | "error";
  duration: number;
}

export class PopupComp extends MinzeElement {
  attrs: Attrs = [["kind", "info"], ["duration", 5000], ["no-css-reset", ""]];

  createdAt: number = 0;
  forceHide: boolean = false;

  static popups: PopupComp[] = [];

  html = () => `
    <slot></slot>
  `

  css = () => `
    :host {
      display: block;

      position: absolute;
      top: 1rem;
      right: 1rem;

      z-index: 9999 !important;

      background-color: ${this.kind == "error" ? "#C23B22" : "#6C8DFA"};

      padding-left: 0.75rem;
      padding-right: 0.75rem;
      padding-top: 0.5rem;
      padding-bottom: 1rem;

      min-width: 10rem;
      height: 3rem;

      border-radius: 0.15rem;

      text-align: right;

      cursor: pointer;

      opacity: 0;
      transform: translateY(0);
      transition: opacity 150ms, filter 150ms, transform 150ms;
    }

    :host::after {
      content: "";

      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;

      background-color: rgba(0, 0, 0, 0.2);

      height: 0.5rem;
      width: var(--progress);
      filter: blur(0);
    }

    :host(.show) {
      opacity: 1;
    }

    :host(.hide) {
      opacity: 0;
      filter: blur(15px);
      pointer-events: none;
    }
  `

  updateIndex() {
    this.style.transform = `translateY(${this.myIndex()*4}rem)`;
  }

  myIndex() {
    return PopupComp.popups.findIndex(a => a === this);
  }

  frame() {
    const elapsed = Date.now() - this.createdAt;
    const progress = (elapsed/this.duration)*100;
    this.style.setProperty("--progress", progress+"%");

    if (this.forceHide || progress > 100) {
      this.classList.remove("show");
      this.classList.add("hide");
      setTimeout(() => {
        this.remove();
        PopupComp.popups.splice(this.myIndex(), 1);
        for (const popup of PopupComp.popups) {
          popup.updateIndex();
        }
      }, 175);
    }
    else {
      requestAnimationFrame(this.frame.bind(this));
    }
  }

  onReady() {
    setTimeout(() => {
      this.classList.add("show");
    }, 100);
    this.createdAt = Date.now();
    requestAnimationFrame(this.frame.bind(this));

    this.addEventListener("click", this.handleClick.bind(this));

    PopupComp.popups.push(this);
    for (const popup of PopupComp.popups) {
      popup.updateIndex();
    }
  }

  handleClick() {
    this.forceHide = true;
  }
}

PopupComp.define("popup-spawn")

export function createPopup(
  kind: "error" | "info",
  message: string,
  duration: number = 5000,
) {
  const popup = document.createElement("popup-spawn") as PopupComp;

  popup.appendChild(document.createTextNode(message));
  popup.kind = kind;
  popup.duration = duration;

  document.body.appendChild(popup);
}
