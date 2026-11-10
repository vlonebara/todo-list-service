// request отправляет JSON-запрос и возвращает понятную ошибку сервера.
export async function request(path, { method = "GET", body, redirect = true } = {}) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 10000);
    try {
        const response = await fetch(`/api/${path}`, {
            method,
            credentials: "same-origin",
            signal: controller.signal,
            headers: body === undefined ? {} : { "Content-Type": "application/json" },
            body: body === undefined ? undefined : JSON.stringify(body),
        });
        const data = await response.json();
        if (!response.ok) {
            if (response.status === 401 && redirect) location.assign("/login.html");
            const error = new Error(data.error || "request failed");
            error.status = response.status;
            throw error;
        }
        return data;
    } catch (error) {
        if (error.name === "AbortError") throw new Error("request timed out; refresh before trying again");
        if (error instanceof TypeError) throw new Error("could not connect to the server");
        throw error;
    } finally {
        clearTimeout(timeout);
    }
}

// query кодирует параметры, сохраняя плюсы, амперсанды и другие символы поиска.
export function query(values) {
    return new URLSearchParams(values).toString();
}
