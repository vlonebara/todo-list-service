// today возвращает сегодняшнюю дату браузера без преобразования в UTC.
export function today() {
    const date = new Date();
    return `${date.getFullYear()}${String(date.getMonth() + 1).padStart(2, "0")}${String(date.getDate()).padStart(2, "0")}`;
}

// inputDate переводит дату API в значение календарного поля.
export function inputDate(value) {
    return `${value.slice(0, 4)}-${value.slice(4, 6)}-${value.slice(6, 8)}`;
}

// displayDate показывает дату задачи в привычном формате.
export function displayDate(value) {
    return `${value.slice(6, 8)}.${value.slice(4, 6)}.${value.slice(0, 4)}`;
}

// repeatLabel описывает правило повторения для списка задач.
export function repeatLabel(rule) {
    const [mode, days, months] = rule.split(" ");
    if (mode === "d") return `Каждые ${days} дн.`;
    if (mode === "y") return "Ежегодно";
    if (mode === "w") {
        const names = ["", "пн", "вт", "ср", "чт", "пт", "сб", "вс"];
        return `Еженедельно: ${days.split(",").map(day => names[Number(day)]).join(", ")}`;
    }
    if (mode === "m") {
        const labels = days.split(",").map(day => ({ "-1": "последний день", "-2": "предпоследний день" })[day] || day);
        return `Ежемесячно: ${labels.join(", ")}${months ? `; месяцы: ${months}` : ""}`;
    }
    return "";
}
