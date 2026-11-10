import { request } from "./api.js";
import { runAction, showMessage } from "./ui.js";

const form = document.querySelector("#login-form");
const submit = document.querySelector("#signin");
const message = document.querySelector("#login-error");

// signin отправляет пароль, а защищённую cookie устанавливает сервер.
function signin(event) {
    event.preventDefault();
    runAction(submit, message, async () => {
        await request("signin", { method: "POST", body: { password: form.elements.password.value }, redirect: false });
        location.assign("/");
    });
}

form.addEventListener("submit", signin);

request("session", { redirect: false }).then(session => {
    if (session.authenticated) location.replace("/");
}).catch(error => showMessage(message, error.message, true));
