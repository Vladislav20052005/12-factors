

// Карточка
function cardHTML(link, text, id) {
    return `
    <div class="card">
        <div class="card-image">
            <img src="${link}" alt="Изображение">
        </div>
        <div class="card-body">
            <a href="${text}" target="_blank" rel="noopener noreferrer">${text}</a>
        </div>
        <div class="card-actions">
            <button class="card-btn" onclick="toggleModal(true, '${id}')">изменить</button>
            <button class="card-btn" onclick="deleteCard('${id}')">удалить</button>
        </div>
    </div>
    `;
}

// Все карты, отображаемые на странице
let cards = []
// Получить карты с сервера и отобразить
async function getCards() {
    console.log("Получение карточек")
    const container = document.getElementById('cards-container');
    
    try {
        const response = await fetch(`${url}/ads`, {
            method: 'GET',
            credentials: 'include'
        });
        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
            cards = result

            container.replaceChildren(); 

            if (cards.length === 0) {
                container.insertAdjacentHTML('beforeend', "<div>No cards</div>");
            } else {
                
                cards.forEach((element) => {
                    container.insertAdjacentHTML('beforeend', cardHTML(`${url}/files/${element.file}?v=` + new Date().getTime(), element.url, element.adId));
                });
                
            }
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log(errorResult)
            if (response.status === 401) {
                redirectToLogin()
            }
        }
    }
    catch (err) {
        console.error('Ошибка сети:', err);
    }
}

getCards()


// Удалить карту
async function deleteCard(id) { 
    console.log(`Удаление карточки: ${id}`)
    try {
        const response = await fetch(`${url}/ads/${id}`, {
            method: 'DELETE',
            credentials: 'include'
        });
    
        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
            getCards()
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log(errorResult)
            if (response.status === 401) {
                redirectToLogin()
            }
        }
    }
    catch (err) {
        console.error('Ошибка сети:', err);
    }
}

// Добавить карту в модальном окне
async function addCard(formData) {
    console.log(`Добавление карточки`)
    try {
        const response = await fetch(`${url}/ads`, {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);

            getCards()
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log(errorResult)
            if (response.status === 401) {
                redirectToLogin()
            }
        }
    } catch (err) {
        console.error('Ошибка сети:', err);
    }
}

// Изменить карту в модальном окне
async function editCard(formData, id) {
    console.log(`Изменение карточки: ${id}`)
    try {
        const response = await fetch(`${url}/ads/${id}`, {
            method: 'PUT',
            body: formData,
            credentials: 'include'
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);

            await getCards()
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log(errorResult)
            if (response.status === 401) {
                redirectToLogin()
            }
        }
    } catch (err) {
        console.error('Ошибка сети:', err);
    }
}