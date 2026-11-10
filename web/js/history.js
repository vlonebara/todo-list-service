import { displayDate } from "./dates.js";
import { element } from "./ui.js";

const actionNames = { add: "Создана", edit: "Изменена", complete: "Выполнена", delete: "Удалена" };

// snapshot показывает сохранённые поля задачи до или после действия.
function snapshot(label, task) {
    const block = element("section");
    block.append(element("h4", label));
    if (!task) {
        block.append(element("p", "Задачи нет."));
        return block;
    }
    block.append(element("p", task.title), element("p", displayDate(task.date)));
    if (task.comment) block.append(element("p", task.comment, "comment"));
    block.append(element("p", task.repeat ? `Повторение: ${task.repeat}` : "Без повторения", "muted"));
    return block;
}

// renderHistory показывает снимки задач и результат отмены действия.
export function renderHistory(container, changes) {
    container.replaceChildren();
    if (!changes.length) {
        container.append(element("p", "История появится после первого действия с задачами.", "empty"));
        return;
    }
    for (const change of changes) {
        const task = change.after || change.before;
        const item = element("article", "", "history-item");
        const title = element("h3", `${actionNames[change.action] || change.action}: ${task.title}`);
        if (change.undone) title.append(element("span", "Отменено", "muted"));
        const timestamp = new Date(change.created_at).toLocaleString("ru-RU");
        const before = change.before ? displayDate(change.before.date) : "";
        const after = change.after ? displayDate(change.after.date) : "";
        const dates = before && after && before !== after ? `${before} → ${after}` : before || after;
        item.append(title, element("p", `${timestamp}; дата задачи: ${dates}`, "muted"));
        const details = element("details");
        details.append(element("summary", "Данные задачи"), snapshot("До действия", change.before), snapshot("После действия", change.after));
        item.append(details);
        container.append(item);
    }
}
