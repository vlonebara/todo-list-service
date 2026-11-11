import { request, query } from "./api.js";
import { createTaskForm } from "./task_form.js";
import { renderTasks } from "./task_list.js";
import { renderHistory } from "./history.js";
import { setupReminders } from "./reminders.js";
import { createPagination } from "./pagination.js";
import { runAction, showMessage } from "./ui.js";

const message = document.querySelector("#message");
const taskContainer = document.querySelector("#tasks");
const historyContainer = document.querySelector("#history");
const more = document.querySelector("#more-tasks");
const moreHistory = document.querySelector("#more-history");
const search = document.querySelector("#search");
const period = document.querySelector("#period");
let historyVisible = false;
const checkReminders = setupReminders();
const taskForm = createTaskForm(refresh);

const loadTaskPage = createPagination("tasks", "tasks", (tasks, hasMore) => {
    renderTasks(taskContainer, tasks, actions);
    more.hidden = !hasMore;
    document.querySelector("#task-count").textContent = `Показано задач: ${tasks.length}`;
});

const loadHistory = createPagination("history", "history", (changes, hasMore) => {
    renderHistory(historyContainer, changes);
    moreHistory.hidden = !hasMore;
});

// loadTasks передаёт текущие фильтры общей функции загрузки страниц.
function loadTasks(reset = true) {
    return loadTaskPage(reset, { search: search.value.trim(), period: period.value });
}

// refresh обновляет данные после одного завершённого действия.
async function refresh() {
    await loadTasks();
    if (historyVisible) await loadHistory();
    await checkReminders();
    showMessage(message, "Изменения сохранены.");
}

const actions = {
    edit: (task, button) => runAction(button, message, () => taskForm.open(task.id)),
    done: (task, button) => runAction(button, message,
        () => changeTask(`task/done?${query({ id: task.id })}`, { method: "POST", body: {} })),
    remove: (task, button) => {
        if (!confirm(`Удалить задачу «${task.title}»?`)) return;
        return runAction(button, message,
            () => changeTask(`task?${query({ id: task.id })}`, { method: "DELETE" }));
    },
};

// changeTask выполняет изменение и обновляет связанные списки и напоминания.
async function changeTask(path, options) {
    await request(path, options);
    await refresh();
}

// reportError показывает ошибку загрузки в общем сообщении страницы.
function reportError(error) {
    showMessage(message, error.message, true);
}

// bindAction связывает кнопку с общей блокировкой и обработкой ошибок.
function bindAction(selector, action) {
    const button = document.querySelector(selector);
    button.addEventListener("click", () => runAction(button, message, () => action(button)));
}

// toggleHistory переключает список задач и историю действий.
async function toggleHistory(button) {
    historyVisible = !historyVisible;
    document.querySelector("#task-view").hidden = historyVisible;
    document.querySelector("#history-view").hidden = !historyVisible;
    button.textContent = historyVisible ? "К задачам" : "История";
    if (historyVisible) await loadHistory();
}

// undoLast восстанавливает данные после последнего действия.
async function undoLast() {
    await changeTask("task/undo", { method: "POST", body: {} });
    showMessage(message, "Последнее действие отменено.");
}

// signout завершает сеанс текущего браузера.
async function signout() {
    await request("signout", { method: "POST", body: {} });
    location.assign("/login.html");
}

// exportTasks сохраняет все задачи в читаемый JSON-файл.
async function exportTasks() {
    const data = await request("tasks/export");
    const url = URL.createObjectURL(new Blob([JSON.stringify(data, null, 2)], { type: "application/json" }));
    const link = document.createElement("a");
    link.href = url;
    link.download = "tasks.json";
    link.hidden = true;
    document.body.append(link);
    link.click();
    link.remove();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
}

bindAction("#add-task", () => taskForm.open());
bindAction("#show-history", toggleHistory);
bindAction("#undo", undoLast);
bindAction("#signout", signout);
bindAction("#export", exportTasks);
bindAction("#more-tasks", () => loadTasks(false));
bindAction("#more-history", () => loadHistory(false));
document.querySelector("#search-form").addEventListener("submit", event => {
    event.preventDefault();
    runAction(document.querySelector("#find"), message, () => loadTasks());
});
period.addEventListener("change", () => loadTasks().catch(reportError));
search.addEventListener("input", () => {
    if (!search.value) loadTasks().catch(reportError);
});
request("session").then(session => {
    document.querySelector("#signout").hidden = !session.required;
}).catch(reportError);
loadTasks().catch(reportError);
