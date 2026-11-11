import { request, query, pageSize } from "./api.js";

// createPagination задаёт общую загрузку страниц для списка задач и истории.
export function createPagination(path, key, render) {
    let items = [];
    let page = 0;
    let generation = 0;

    // load добавляет страницу или заменяет список, игнорируя запоздавшие ответы.
    async function load(reset = true, params = {}) {
        const version = ++generation;
        const nextPage = reset ? 1 : page + 1;
        const result = await request(`${path}?${query({ ...params, page: nextPage })}`);
        if (version !== generation) return;
        const received = result[key];
        items = reset ? received : [...items, ...received];
        page = nextPage;
        render(items, received.length >= pageSize);
    }

    return load;
}
