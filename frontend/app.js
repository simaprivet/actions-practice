// Загружаем записи с backend и выводим списком
async function load() {
  const list = document.getElementById("items");
  const status = document.getElementById("status");
  try {
    const res = await fetch("/api/items");
    if (!res.ok) throw new Error("HTTP " + res.status);
    const items = await res.json();
    list.innerHTML = "";
    for (const it of items) {
      const li = document.createElement("li");
      li.textContent = it.id + ". " + it.name;
      list.appendChild(li);
    }
    status.textContent = items.length ? "Всего записей: " + items.length : "Записей пока нет";
  } catch (e) {
    status.textContent = "Ошибка загрузки: " + e.message;
  }
}

load();
