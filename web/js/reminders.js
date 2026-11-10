import { request, query } from "./api.js";
import { showMessage } from "./ui.js";

// preference читает настройку даже при запрещённом хранилище браузера.
function preference(key) {
    try { return localStorage.getItem(key); } catch { return null; }
}

// savePreference сохраняет только настройку, без токенов и текста задач.
function savePreference(key, value) {
    try { localStorage.setItem(key, value); } catch { /* Напоминания продолжат работать в этой вкладке. */ }
}

// setupReminders напоминает о сегодняшних задачах, пока вкладка открыта.
export function setupReminders() {
    const toggle = document.querySelector("#reminders");
    const message = document.querySelector("#reminder-message");
    toggle.checked = preference("todo.reminders") === "true";
    let checking = false;

    // check обновляет сообщение во вкладке и не создаёт параллельных запросов.
    async function check() {
        if (!toggle.checked || checking) return;
        checking = true;
        try {
            const { tasks } = await request(`tasks?${query({ period: "today" })}`);
            if (!toggle.checked) return;
            showMessage(message, tasks.length ? "На сегодня есть задачи. Откройте фильтр «Сегодня»." : "На сегодня задач нет.");
            document.title = tasks.length ? "Есть задачи на сегодня — Планировщик" : "Планировщик задач";

        } catch (error) {
            showMessage(message, error.message, true);
        } finally {
            checking = false;
        }
    }

    toggle.addEventListener("change", async () => {
        savePreference("todo.reminders", String(toggle.checked));
        if (!toggle.checked) {
            showMessage(message, "");
            document.title = "Планировщик задач";
            return;
        }

        await check();
    });
    const timer = setInterval(check, 60000);
    window.addEventListener("pagehide", () => clearInterval(timer), { once: true });
    check();
    return check;
}
