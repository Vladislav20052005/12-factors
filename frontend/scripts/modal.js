// Закрыть окно при клике вовне
function closeModalOnOverlay(event) {
    if (event.target.id === 'postModal') {
        toggleModal(false);
    }
}

// Отобразить имя выбранного файла
function updateFileName(input) {
    const fileNameDisplay = document.getElementById('fileNameDisplay');
    if (input.files && input.files[0]) {
        fileNameDisplay.textContent = `Выбран файл: ${input.files[0].name}`;
    } else {
        fileNameDisplay.textContent = 'Файл не выбран';
    }
}

// Айди изменяемой карты
let currentEditingCardId = null;

// Открыть/закрыть окно
function toggleModal(show, cardId = null) {
    const modal = document.getElementById('postModal');
    const form = document.getElementById('createPostForm');
    const modalTitle = modal.querySelector('.modal-header h3');
    const submitBtn = modal.querySelector('.submit-post-btn');
    const fileInput = document.getElementById('fileInput');

    if (show) {
        modal.classList.add('active');
        
        if (cardId) {
            currentEditingCardId = cardId;
            modalTitle.textContent = 'Редактировать карточку';
            submitBtn.textContent = 'Сохранить изменения';
            
            fileInput.removeAttribute('required');
            fileInput.removeAttribute('required');

            const cardData = cards.find(c => c.adId === cardId); 
            
            if (cardData) {
                document.getElementById('textInput').value = cardData.url;
                document.getElementById('fileNameDisplay').textContent = 'Оставьте пустой, чтобы не менять фото';
            }
        } else {
            currentEditingCardId = null;
            modalTitle.textContent = 'Новая карточка';
            submitBtn.textContent = 'Опубликовать карточку';
            fileInput.setAttribute('required', 'required');
            
            form.reset();
            document.getElementById('fileNameDisplay').textContent = 'Файл не выбран';
        }
    } else {
        modal.classList.remove('active');
        form.reset();
        currentEditingCardId = null;
    }
}

// Отправить форму
async function handleFormSubmit(event) {
    event.preventDefault();

    const fileInput = document.getElementById('fileInput');
    const textContent = document.getElementById('textInput').value;


    const file = fileInput.files[0] ? fileInput.files[0] : null;

    const formData = new FormData();
    formData.append('Image', file);
    formData.append('URL', textContent);

    if (currentEditingCardId) {
        editCard(formData, currentEditingCardId)
    } else {
        addCard(formData)
    }

    toggleModal(false);
}
