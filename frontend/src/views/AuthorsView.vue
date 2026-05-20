<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../api/index.js';

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

const filteredAuthors = computed(() => {
  const q = searchQuery.value.trim().toLowerCase();
  if (!q) return authors.value;
  return authors.value.filter(a => fullName(a).toLowerCase().includes(q));
});
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
      <p class="count">{{ filteredAuthors.length }} художник{{ filteredAuthors.length === 1 ? '' : filteredAuthors.length < 5 ? 'а' : 'ов' }}</p>
    </div>

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
      </div>
    </section>

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

@media (max-width: 1100px) {
  .grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 750px) {
  .container { padding: 0 20px; }
  .grid { grid-template-columns: repeat(2, 1fr); gap: 16px; }
  .card { padding: 32px 20px 28px; }
}

@media (max-width: 480px) {
  .grid { grid-template-columns: 1fr; }
  .page-head { flex-direction: column; align-items: flex-start; gap: 12px; }
  .search-wrap { max-width: 100%; width: 100%; }
}
</style>
