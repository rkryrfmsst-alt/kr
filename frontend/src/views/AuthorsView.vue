<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { api } from '../api/index.js';
import { useAuth } from '../composables/useAuth.js';

const router = useRouter();
const { adminMode } = useAuth();

const authors   = ref([]);
const isLoading = ref(true);
const loadError = ref('');

const searchQuery = ref('');

onMounted(async () => {
  try {
    authors.value = await api.getAuthors();
  } catch {
    loadError.value = 'Не удалось загрузить художников.';
  } finally {
    isLoading.value = false;
  }
});

const fullName = (a) =>
  [a.last_name, a.first_name, a.middle_name].filter(Boolean).join(' ');

const initials = (a) =>
  [a.first_name, a.last_name]
    .map(n => n?.[0] ?? '')
    .join('');

const goToWorks = (author) => {
  router.push({ path: '/', query: { author: fullName(author) } });
};

const filteredAuthors = computed(() => {
  const q = searchQuery.value.trim().toLowerCase();
  if (!q) return authors.value;
  return authors.value.filter(a => fullName(a).toLowerCase().includes(q));
});

// --- Add author modal ---
const showModal    = ref(false);
const addLoading   = ref(false);
const addError     = ref('');
const newFirstName = ref('');
const newLastName  = ref('');
const newMiddleName  = ref('');
const newDescription = ref('');

function openModal() {
  newFirstName.value  = '';
  newLastName.value   = '';
  newMiddleName.value  = '';
  newDescription.value = '';
  addError.value = '';
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
}

async function submitAddAuthor() {
  addError.value = '';
  addLoading.value = true;
  try {
    const author = await api.createAuthor({
      first_name:  newFirstName.value.trim(),
      last_name:   newLastName.value.trim(),
      middle_name: newMiddleName.value.trim() || null,
      description: newDescription.value.trim() || null,
    });
    authors.value.push(author);
    closeModal();
  } catch (e) {
    addError.value = e.message;
  } finally {
    addLoading.value = false;
  }
}

// --- Edit author modal ---
const showEditModal    = ref(false);
const editAuthorID     = ref(null);
const editFirstName    = ref('');
const editLastName     = ref('');
const editMiddleName   = ref('');
const editDescription  = ref('');
const editLoading      = ref(false);
const editError        = ref('');

function openEditModal(author) {
  editAuthorID.value    = author.id;
  editFirstName.value   = author.first_name;
  editLastName.value    = author.last_name;
  editMiddleName.value  = author.middle_name ?? '';
  editDescription.value = author.description ?? '';
  editError.value       = '';
  showEditModal.value   = true;
}

function closeEditModal() {
  showEditModal.value = false;
}

async function submitEditAuthor() {
  editError.value   = '';
  editLoading.value = true;
  try {
    const updated = await api.updateAuthor(editAuthorID.value, {
      first_name:  editFirstName.value.trim(),
      last_name:   editLastName.value.trim(),
      middle_name: editMiddleName.value.trim() || null,
      description: editDescription.value.trim() || null,
    });
    const idx = authors.value.findIndex(a => a.id === editAuthorID.value);
    if (idx !== -1) authors.value[idx] = updated;
    closeEditModal();
  } catch (e) {
    editError.value = e.message;
  } finally {
    editLoading.value = false;
  }
}

// --- Delete author ---
async function deleteAuthor(author) {
  try {
    await api.deleteAuthor(author.id);
    authors.value = authors.value.filter(a => a.id !== author.id);
  } catch (e) {
    loadError.value = e.message;
  }
}
</script>

<template>
  <div class="container">

    <div class="page-head">
      <div class="search-wrap">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Поиск по имени..."
          class="search-input"
        />
      </div>
      <div class="page-head-right">
        <button v-if="adminMode" class="add-btn" @click="openModal">+ Добавить художника</button>
        <p class="count">{{ filteredAuthors.length }} художник{{ filteredAuthors.length === 1 ? '' : filteredAuthors.length < 5 ? 'а' : 'ов' }}</p>
      </div>
    </div>

    <teleport to="body">
      <div v-if="showModal" class="modal-backdrop" @click.self="closeModal">
        <div class="modal">
          <h2 class="modal-title">Новый художник</h2>

          <form class="modal-form" novalidate @submit.prevent="submitAddAuthor">
            <div class="modal-field">
              <label class="modal-label modal-label--req">Фамилия</label>
              <input v-model="newLastName" class="modal-input" type="text" placeholder="Иванов" required />
            </div>
            <div class="modal-field">
              <label class="modal-label modal-label--req">Имя</label>
              <input v-model="newFirstName" class="modal-input" type="text" placeholder="Иван" required />
            </div>
            <div class="modal-field">
              <label class="modal-label">Отчество</label>
              <input v-model="newMiddleName" class="modal-input" type="text" placeholder="Иванович" />
            </div>
            <div class="modal-field">
              <label class="modal-label">Описание</label>
              <textarea v-model="newDescription" class="modal-input modal-textarea" placeholder="Краткая биография..." rows="4"></textarea>
            </div>

            <p v-if="addError" class="modal-error">{{ addError }}</p>

            <div class="modal-actions">
              <button type="button" class="modal-cancel" @click="closeModal">Отмена</button>
              <button type="submit" class="modal-submit" :disabled="addLoading">
                {{ addLoading ? 'Сохранение...' : 'Добавить' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </teleport>

    <p v-if="isLoading" class="no-results">Загрузка...</p>
    <p v-else-if="loadError" class="no-results">{{ loadError }}</p>
    <p v-else-if="filteredAuthors.length === 0" class="no-results">
      Художник не найден.
    </p>

    <section class="grid">
      <div v-for="author in filteredAuthors" :key="author.id" class="card">
        <div class="avatar">
          <span class="initials">{{ initials(author) }}</span>
        </div>

        <div class="card-body">
          <h2 class="name">{{ fullName(author) }}</h2>
          <p class="description">{{ author.description ?? '—' }}</p>
        </div>
        <button class="works-btn" @click="goToWorks(author)">
          К работам
          <img src="../assets/icons/arrow.svg" class="works-arrow" alt="">
        </button>
        <div v-if="adminMode" class="card-admin">
          <button class="admin-action-btn" @click.stop="openEditModal(author)">Изменить</button>
          <button class="admin-action-btn admin-action-btn--del" @click.stop="deleteAuthor(author)">Удалить</button>
        </div>
      </div>
    </section>

    <teleport to="body">
      <div v-if="showEditModal" class="modal-backdrop" @click.self="closeEditModal">
        <div class="modal">
          <h2 class="modal-title">Редактировать художника</h2>
          <form class="modal-form" novalidate @submit.prevent="submitEditAuthor">
            <div class="modal-field">
              <label class="modal-label modal-label--req">Фамилия</label>
              <input v-model="editLastName" class="modal-input" type="text" placeholder="Иванов" required />
            </div>
            <div class="modal-field">
              <label class="modal-label modal-label--req">Имя</label>
              <input v-model="editFirstName" class="modal-input" type="text" placeholder="Иван" required />
            </div>
            <div class="modal-field">
              <label class="modal-label">Отчество</label>
              <input v-model="editMiddleName" class="modal-input" type="text" placeholder="Иванович" />
            </div>
            <div class="modal-field">
              <label class="modal-label">Описание</label>
              <textarea v-model="editDescription" class="modal-input modal-textarea" placeholder="Краткая биография..." rows="4"></textarea>
            </div>
            <p v-if="editError" class="modal-error">{{ editError }}</p>
            <div class="modal-actions">
              <button type="button" class="modal-cancel" @click="closeEditModal">Отмена</button>
              <button type="submit" class="modal-submit" :disabled="editLoading">
                {{ editLoading ? 'Сохранение...' : 'Сохранить' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </teleport>

  </div>
</template>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 50px;
}

.page-head {
  padding: 36px 0 28px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e8e8e8;
  margin-bottom: 40px;
}

.page-head-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.add-btn {
  background: #000;
  color: #fff;
  border: none;
  border-radius: 2px;
  padding: 10px 18px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  cursor: pointer;
  transition: background 0.2s;
  white-space: nowrap;
}
.add-btn:hover { background: #222; }

.search-wrap {
  flex: 1;
  max-width: 320px;
}

.search-input {
  width: 100%;
  background: #ebebeb;
  border: 1px solid transparent;
  border-radius: 2px;
  padding: 10px 14px;
  font-family: inherit;
  font-size: 0.9rem;
  color: #000;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}
.search-input:focus { border-color: #444; }
.search-input::placeholder { color: #aaa; }

.count {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #aaa;
  margin: 0;
}

.no-results {
  text-align: center;
  color: #666;
  padding: 80px 0;
  font-size: 0.9rem;
  letter-spacing: 1px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  padding-bottom: 60px;
}

.card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 44px 32px 36px;
  border: 1px solid #e8e8e8;
  transition: box-shadow 0.2s, border-color 0.2s;
}
.card:hover {
  border-color: #ccc;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.07);
}

.card-body {
  flex: 1;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: #1a1a1a;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  flex-shrink: 0;
}

.initials {
  font-size: 1.5rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.25);
  letter-spacing: 2px;
  user-select: none;
}

.card-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #000;
  margin: 0;
  line-height: 1.3;
}

.description {
  font-size: 0.875rem;
  color: #666;
  margin: 0;
  line-height: 1.65;
}

.works-btn {
  margin-top: 28px;
  background: none;
  border: 1px solid #000;
  border-radius: 2px;
  padding: 10px 27.5px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #000;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background 0.2s, color 0.2s;
}
.works-btn:hover {
  background: #000;
  color: #fff;
}
.works-btn:hover .works-arrow { filter: invert(1); }

.card-admin {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.admin-action-btn {
  flex: 1;
  background: none;
  border: 1px solid #ccc;
  border-radius: 2px;
  padding: 7px 4px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.68rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #666;
  cursor: pointer;
  transition: border-color 0.2s, color 0.2s;
}
.admin-action-btn:hover { border-color: #444; color: #000; }
.admin-action-btn--del:hover { border-color: #c00; color: #c00; }

.works-arrow {
  width: 13px;
  height: 13px;
  transform: rotate(-90deg);
  transition: filter 0.2s;
}

@media (max-width: 1100px) {
  .grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 750px) {
  .container { padding: 0 20px; }
  .grid { grid-template-columns: repeat(2, 1fr); gap: 16px; }
  .card { padding: 32px 20px 28px; }
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 20px;
}

.modal {
  background: #fff;
  width: 100%;
  max-width: 480px;
  padding: 40px;
  border-radius: 4px;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #000;
  margin: 0 0 28px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.modal-field {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.modal-label {
  font-family: 'Raleway', sans-serif;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #000;
}
.modal-label--req::after {
  content: ' *';
  color: #c00;
  font-size: 1rem;
  font-weight: 700;
  font-family: Arial, sans-serif;
  text-transform: none;
  letter-spacing: 0;
}

.modal-input {
  background: #ebebeb;
  border: 1px solid transparent;
  border-radius: 2px;
  padding: 11px 14px;
  font-family: 'Inter', sans-serif;
  font-size: 0.9rem;
  color: #000;
  outline: none;
  transition: border-color 0.2s;
  width: 100%;
  box-sizing: border-box;
}
.modal-input:focus { border-color: #444; }
.modal-input::placeholder { color: #aaa; }

.modal-textarea {
  resize: none;
  min-height: 130px;
}

.modal-error {
  font-size: 0.8rem;
  color: #c00;
  margin: 0;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}

.modal-cancel {
  background: none;
  border: 1px solid #ccc;
  border-radius: 2px;
  padding: 10px 20px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #666;
  cursor: pointer;
  transition: border-color 0.2s, color 0.2s;
}
.modal-cancel:hover { border-color: #999; color: #000; }

.modal-submit {
  background: #000;
  color: #fff;
  border: none;
  border-radius: 2px;
  padding: 10px 24px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  cursor: pointer;
  transition: background 0.2s;
}
.modal-submit:hover:not(:disabled) { background: #222; }
.modal-submit:disabled { opacity: 0.5; cursor: default; }

@media (max-width: 480px) {
  .grid { grid-template-columns: 1fr; }
  .page-head { flex-direction: column; align-items: flex-start; gap: 12px; }
  .search-wrap { max-width: 100%; width: 100%; }
}
</style>
