import { MinzeElement } from "minze";
import "~/src/dashboard/dashboard.ts"
import * as mapi from "~/src/manager_api"
import { createPopup } from "~/src/popup";

(class CreateGameComp extends MinzeElement {
  attrs = [["no-css-reset", ""]] as any;

  html = () => `
    <form class="form" on:submit="handleSubmit">
      <div>
        <h1>Create Game</h1>
      </div>
      <div class="body">
        <label>
          <span>Name</span>
          <input type="text" name="name">
        </label>
        <button type="submit">
          CREATE
        </button>
      </div>
    </form>
  `

  css = () => `
    :host {
      flex-grow: 1;
      display: flex;
      justify-content: center;
      align-items: center;
    }

    h1 {
      font-size: 1.2rem;
      margin-bottom: 2rem;

      text-align: center;
    }

    .form {
      padding: 1rem;
      min-width: min(30rem, 100%);
    }

    .body {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    label {
      display: flex;
      flex-wrap: wrap;
      align-items: center;

      gap: 1rem;
    }

    input[type=text] {
      display: block;

      background: rgba(255, 255, 255, 0.1);
      padding: 0.5rem;
      border-radius: 0.2rem;
      flex-grow: 1;
      min-width: min(100%, 10rem);
    }

    button {
      width: 100%;
      background-color: #6C8DFA;
      padding: 0.5rem;
      border-radius: 0.2rem;
    }

    button:disabled {
      cursor: normal;
    }
  `

  handleSubmit(event: SubmitEvent) {
    const form = event.target as HTMLFormElement;
    const name = form.elements["name" as any] as HTMLInputElement;
    const button = this.select("button") as HTMLButtonElement;
    button.disabled = true;

    console.log(name);
    
    mapi.createGame(name.value).then(() => {
      document.location.href = "/";
    }).catch(e => {
      button.disabled = false;
      if (e instanceof mapi.RequestError) {
        e.response.json().then(body => {
          createPopup("error", body["message"] ?? body["error"] ?? `Could not create game (${e.response.status})`)
        });
      }
    });
    
    event.preventDefault();
  }
}).define("creeps-create-game");

