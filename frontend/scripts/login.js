
// Войти
async function login() {
    const loginInput = document.getElementById('username').value.trim();
    const passwordInput = document.getElementById('password').value.trim();

    console.log(`Вход в аккаунт : ${loginInput}, ***`)
    
    try {
        const response = await fetch(`${url}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: loginInput,
                password: passwordInput
            }),
            credentials: "include",
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
                
            redirectToDashboard(); 
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log(errorResult, response.status, response.statusText)
        }
    } catch (err) {
        console.error('Ошибка сети:', err);
    }
}


// Зарегистрироваться
async function signup() {
    const loginInput = document.getElementById('username').value.trim();
    const passwordInput = document.getElementById('password').value.trim();

    console.log(`Создание аккаунта : ${loginInput}, ***`)
    try {
        const response = await fetch(`${url}/auth/signup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: loginInput,
                password: passwordInput
            }),
            credentials: "include",
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
                
            redirectToDashboard(); 
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log(errorResult)
        }
    } catch (err) {
        console.error('Ошибка сети:', err);
    }
}