// showMessage выводит текст без интерпретации HTML.
export function showMessage(element, text, error = false) {
    if (element.textContent !== text) element.textContent = text;
    element.hidden = !text;
    element.classList.toggle("error", error);
}

// runAction блокирует повторный запрос и всегда возвращает кнопку в исходное состояние.
export async function runAction(button, message, action) {
    if (button.disabled) return;
    button.disabled = true;
    button.setAttribute("aria-busy", "true");
    showMessage(message, "");
    try {
        await action();
    } catch (error) {
        showMessage(message, error.message, true);
    } finally {
        button.disabled = false;
        button.removeAttribute("aria-busy");
    }
}

// element создаёт DOM-элемент с безопасным текстовым содержимым.
export function element(tag, text = "", className = "") {
    const node = document.createElement(tag);
    node.textContent = text;
    node.className = className;
    return node;
}
