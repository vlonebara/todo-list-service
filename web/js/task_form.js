import { request, query } from "./api.js";
import { today, inputDate } from "./dates.js";
import { element, runAction, showMessage } from "./ui.js";

// createChoices создаёт одинаковые группы флажков для дней недели и месяцев.
function createChoices(container, name, labels, checked) {
    labels.forEach((text, index) => {
        const label = element("label");
        const input = document.createElement("input");
        input.type = "checkbox";
        input.name = name;
        input.value = String(index + 1);
        input.defaultChecked = checked;
        label.append(input, document.createTextNode(text));
        container.append(label);
    });
}

// selectedValues читает отмеченные значения группы флажков.
function selectedValues(form, name) {
    return [...form.querySelectorAll(`[name="${name}"]:checked`)].map(input => input.value);
}

// repeatRule формирует проверенное правило из полей формы.
export function repeatRule(form) {
    const fields = form.elements;
    const mode = fields.mode.value;
    if (!mode || mode === "y") return mode;
    if (mode === "d") {
        const interval = Number(fields.interval.value);
        if (!Number.isInteger(interval) || interval < 1 || interval > 400) throw new Error("day interval must be between 1 and 400");
        return `d ${interval}`;
    }
    if (mode === "w") {
        const days = selectedValues(form, "weekdays");
        if (!days.length) throw new Error("select at least one weekday");
        return `w ${days.join(",")}`;
    }
    const days = fields.monthdays.value.trim() ? fields.monthdays.value.split(",").map(value => Number(value.trim())) : [];
    if (days.some(day => !Number.isInteger(day) || day < 1 || day > 31)) throw new Error("month days must be between 1 and 31");
    if (fields.last.checked) days.push(-1);
    if (fields.penultimate.checked) days.push(-2);
    if (!days.length) throw new Error("select at least one month day");
    const months = selectedValues(form, "months");
    if (!months.length) throw new Error("select at least one month");
    return `m ${[...new Set(days)].join(",")}${months.length === 12 ? "" : ` ${months.join(",")}`}`;
}

// fillRepeat заполняет правило только после полного сброса формы.
function fillRepeat(form, rule) {
    const [mode = "", values = "", months] = rule.split(" ");
    const fields = form.elements;
    fields.mode.value = mode;
    if (mode === "d") fields.interval.value = values;
    if (mode === "w") {
        for (const input of form.querySelectorAll('[name="weekdays"]')) input.checked = values.split(",").map(Number).includes(Number(input.value));
    }
    if (mode === "m") {
        const days = values.split(",").map(Number);
        fields.monthdays.value = days.filter(day => day > 0).join(",");
        fields.last.checked = days.includes(-1);
        fields.penultimate.checked = days.includes(-2);
        if (months) {
            for (const input of form.querySelectorAll('[name="months"]')) input.checked = months.split(",").map(Number).includes(Number(input.value));
        }
    }
}

// updateRepeatFields показывает поля только выбранного правила.
function updateRepeatFields(form) {
    for (const group of form.querySelectorAll("[data-repeat]")) {
        group.hidden = group.dataset.repeat !== form.elements.mode.value;
        group.disabled = group.hidden;
    }
}

// createTaskForm связывает календарь, повторения и сохранение задачи.
export function createTaskForm(onSaved) {
    const dialog = document.querySelector("#task-dialog");
    const form = document.querySelector("#task-form");
    const error = document.querySelector("#form-error");
    const save = document.querySelector("#save-task");
    const close = document.querySelector("#close-task");
    const fieldsGroup = document.querySelector("#task-fields");
    let taskID = "";
    let opening = 0;
    createChoices(form.querySelector("#weekdays"), "weekdays", ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"], false);
    createChoices(form.querySelector("#months"), "months", ["Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"], true);
    form.elements.mode.addEventListener("change", () => updateRepeatFields(form));
    document.querySelector("#close-task").addEventListener("click", () => dialog.close());
    dialog.addEventListener("cancel", event => { if (save.disabled) event.preventDefault(); });
    form.addEventListener("submit", event => {
        event.preventDefault();
        runAction(save, error, async () => {
            const fields = form.elements;
            const task = {
                id: taskID,
                date: fields.date.value.replaceAll("-", ""),
                title: fields.title.value,
                comment: fields.comment.value,
                repeat: repeatRule(form),
            };
            fieldsGroup.disabled = true;
            close.disabled = true;
            try {
                await request("task", { method: taskID ? "PUT" : "POST", body: task });
                dialog.close();
                await onSaved();
            } finally {
                fieldsGroup.disabled = false;
                close.disabled = false;
            }
        });
    });

    // open всегда создаёт чистое состояние перед заполнением выбранной задачи.
    async function open(id = "") {
        const version = ++opening;
        const task = id ? await request(`task?${query({ id })}`) : { date: today(), title: "", comment: "", repeat: "" };
        if (version !== opening) return;
        form.reset();
        taskID = id;
        form.elements.date.value = inputDate(task.date);
        form.elements.title.value = task.title;
        form.elements.comment.value = task.comment;
        fillRepeat(form, task.repeat);
        updateRepeatFields(form);
        showMessage(error, "");
        document.querySelector("#form-title").textContent = id ? "Карточка задачи" : "Новая задача";
        dialog.showModal();
        form.elements.title.focus();
    }
    return { open };
}
