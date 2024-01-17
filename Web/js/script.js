document.addEventListener('DOMContentLoaded', function () {
    function handleFormSubmit(event) {
        event.preventDefault();

        const idValue = document.getElementById("idInput").value;

        fetch("/get-json/" + idValue)
            .then(response => response.json())
            .then(data => {
                console.log("Ответ от сервера:", data);

                document.getElementById("result").innerHTML = `
                    <p>Информация по ID ${idValue}:</p>
                    <pre>${JSON.stringify(data, null, 2)}</pre>
                `;
            })
            .catch(error => {
                console.error("Ошибка:", error);
            });
    }

    function handleAddInfoFormSubmit(event) {
        event.preventDefault();

        const jsonData = document.getElementById("inputJson").value;

        fetch("/add-json", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ jsonData }),
        })
            .then(response => response.json())
            .then(data => {
                console.log("Ответ от сервера:", data);

                document.getElementById("result").innerHTML = `
                    <p>Данные успешно добавлены:</p>
                    <pre>${JSON.stringify(data, null, 2)}</pre>
                `;
            })
            .catch(error => {
                console.error("Ошибка:", error);
            });
    }

    function handleGetAllIdsFormSubmit(event) {
        event.preventDefault();

        fetch("/get-all-ids")
            .then(response => response.json())
            .then(data => {
                console.log("Ответ от сервера:", data);

                // Вывод всех ID в блок all-ids
                document.getElementById("all-ids").innerHTML = `
                    <h4>Все ID:</h4>
                    <pre>${JSON.stringify(data, null, 2)}</pre>
                `;
            })
            .catch(error => {
                console.error("Ошибка:", error);
            });
    }

    const getInfoForm = document.getElementById('get-info-form');
    const addInfoForm = document.getElementById('add-info-form');
    const getAllIdsForm = document.getElementById('get-all-ids-form');

    if (getInfoForm) {
        getInfoForm.addEventListener('submit', handleFormSubmit);
    }

    if (addInfoForm) {
        addInfoForm.addEventListener('submit', handleAddInfoFormSubmit);
    }

    if (getAllIdsForm) {
        getAllIdsForm.addEventListener('submit', handleGetAllIdsFormSubmit);
    }
});
