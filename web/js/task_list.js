import { displayDate, inputDate, repeatLabel, today } from "./dates.js";
import { element } from "./ui.js";

// actionButton создаёт доступную с клавиатуры кнопку действия.
function actionButton(label, action) {
    const button = element("button", label, "quiet");
    button.type = "button";
    button.addEventListener("click", () => action(button));
    return button;
}

// renderTasks группирует задачи по датам и выводит пользовательские данные текстом.
export function renderTasks(container, tasks, actions) {
    container.replaceChildren();
    if (!tasks.length) {
        container.append(element("p", "Задач не найдено. Измените фильтр или добавьте новую задачу.", "empty"));
        return;
    }
    let date = "";
    for (const task of tasks) {
        if (task.date !== date) {
            date = task.date;
            const heading = element("h2", displayDate(date), "date-heading");
            if (date < today()) heading.append(element("span", "Просрочено", "overdue"));
            container.append(heading);
        }
        const card = element("article", "", "task");
        const content = element("div", "", "task-content");
        content.append(element("h3", task.title));
        if (task.comment) content.append(element("p", task.comment, "comment"));
        if (task.repeat) content.append(element("p", repeatLabel(task.repeat), "muted"));
        const time = element("time", displayDate(task.date), "visually-hidden");
        time.dateTime = inputDate(task.date);
        content.append(time);
        const buttons = element("div", "", "task-actions");
        buttons.append(
            actionButton("Выполнить", button => actions.done(task, button)),
            actionButton("Изменить", button => actions.edit(task, button)),
            actionButton("Удалить", button => actions.remove(task, button)),
        );
        card.append(content, buttons);
        container.append(card);
    }
}
