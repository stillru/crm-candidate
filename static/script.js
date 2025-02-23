document.addEventListener("DOMContentLoaded", function () {
    // Обработка кликов по ссылкам в меню
    document.querySelectorAll(".sidebar a").forEach(link => {
        link.addEventListener("click", function (event) {
            event.preventDefault(); // Отменяем стандартное поведение ссылки
            const url = this.getAttribute("href"); // Получаем URL из атрибута href

            // Загружаем контент с сервера
            fetch(url)
                .then(response => response.text())
                .then(html => {
                    // Обновляем только рабочую область
                    document.querySelector(".content").innerHTML = html;
                })
                .catch(error => console.error("Ошибка при загрузке контента:", error));
        });
    });
});
