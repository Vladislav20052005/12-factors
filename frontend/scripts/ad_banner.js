// Рекламный баннер
function bannerHTML(uniqueId) {
    return `
    <div class="ads-banner" id="${uniqueId}">
        <a href="" target="_blank" rel="noopener noreferrer">
            <img src="" alt="Реклама">
        </a>
        <style>
            .ads-banner {
                width: 300px;
                height: 250px;
                margin: 20px auto;
                overflow: hidden;
                border-radius: 8px;
                box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
                background-color: #f9f9f9;
            }
            .ads-banner a {
                display: block;
                width: 100%;
                height: 100%;
            }
            .ads-banner img {
                width: 100%;
                height: 100%;
                object-fit: cover;
                transition: opacity 0.2s ease-in-out;
            }
            .ads-banner img:hover {
                opacity: 0.9;
                cursor: pointer;
            }
        </style>
    </div>
    `;
}

// Запрос на сервер для отображения баннера
async function fetchBannerData(uniqueId, backendUrl, token) {
    console.log("Получение баннера", uniqueId, backendUrl, token)
    const container = document.getElementById(uniqueId);
    if (!container) return;

    const pic = container.querySelector('img');
    const link = container.querySelector('a');

    try {
        const response = await fetch(`${backendUrl}/show/${token}`, {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Сервер ответил успехом:', result);
            pic.src = backendUrl + '/files/' + result.file;
            link.href = result.url;
        } else {
            const errorResult = await response.json().catch(() => ({}));
            console.log('Ошибка сервера:', errorResult);
        }
    } catch (err) {
        console.error('Ошибка сети:', err);
    }
}

async function banner() {
    const holder = document.getElementById('banner-holder');

    const uniqueId = 'banner_' + Math.floor(Math.random() * 1000000);
    
    if (holder) {
        holder.insertAdjacentHTML('beforeend', bannerHTML(uniqueId));
        
        fetchBannerData(uniqueId, url, TOKEN);
    }
    
}
setTimeout(banner, 500);
