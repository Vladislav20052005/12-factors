// Копирование при нажатии на кнопку
function copyText() {
    const textToCopy = document.getElementById('tokenText').innerText;
    const btn = document.querySelector('.copy-btn');

    navigator.clipboard.writeText(textToCopy).then(() => {
        btn.innerText = 'Скопировано!';
        btn.classList.add('success');

        setTimeout(() => {
            btn.innerText = 'Копировать';
            btn.classList.remove('success');
        }, 2000);
    }).catch(err => {
        console.error('Ошибка копирования: ', err);
    });
}

// Получить информацию о юзере с сервера
async function getUser() {
    const usernameText = document.getElementById('usernameText');
    const tokenText = document.getElementById('tokenText');
    console.log(`Получение пользователя (имени и токена)`)

    
    try {
        const response = await fetch(`${url}/auth`, {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
            usernameText.textContent = result.username
            TOKEN = result.token
            console.log(TOKEN)
            tokenText.textContent = `${url}/show/${TOKEN}`
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

// Удалить юзера
async function deleteUser() {
    console.log(`Удаление аккаунта`)
    try {
        const response = await fetch(`${url}/auth`, {
            method: 'DELETE',
            credentials: 'include',
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
            
            redirectToLogin()
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

getUser()

