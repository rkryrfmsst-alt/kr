<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import PaintingModal from '../components/PaintingModal.vue';
import { api, BASE } from '../api/index.js';

const selectedPainting = ref(null);

const openModal = (painting) => {
  selectedPainting.value = painting;
  const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
  document.body.style.paddingRight = scrollbarWidth + 'px';
  document.body.style.overflow = 'hidden';
};

const closeModal = () => {
  selectedPainting.value = null;
  document.body.style.overflow = '';
  document.body.style.paddingRight = '';
};

const handleKeydown = (e) => {
  if (e.key === 'Escape') closeModal();
};

const paintings  = ref([]);
const isLoading  = ref(true);
const loadError  = ref('');

const authorName = (a) => [a.last_name, a.first_name, a.middle_name].filter(Boolean).join(' ');

const primaryAuthor = (p) => p.authors?.length ? authorName(p.authors[0]) : '—';

const isFilterOpen    = ref(false);
const searchQuery     = ref('');
const selectedStyle   = ref('');
const selectedSubject = ref('');
const selectedAuthor  = ref('');
const selectedTechnique = ref('');
const yearFrom = ref('');
const yearTo   = ref('');

const resetFilters = () => {
  searchQuery.value       = '';
  selectedStyle.value     = '';
  selectedSubject.value   = '';
  selectedAuthor.value    = '';
  selectedTechnique.value = '';
  yearFrom.value          = '';
  yearTo.value            = '';
};

const hasActiveFilters = computed(() =>
  searchQuery.value || selectedStyle.value || selectedSubject.value ||
  selectedAuthor.value || selectedTechnique.value || yearFrom.value || yearTo.value
);

const styles = computed(() =>
  [...new Set(paintings.value.flatMap(p => (p.styles ?? []).map(s => s.name)))].sort()
);
const subjects = computed(() =>
  [...new Set(paintings.value.flatMap(p => (p.plots ?? []).map(pl => pl.name)))].sort()
);
const authors = computed(() =>
  [...new Set(paintings.value.flatMap(p => (p.authors ?? []).map(authorName)))].sort()
);
const techniques = computed(() =>
  [...new Set(paintings.value.filter(p => p.material).map(p => p.material.name))].sort()
);

const isSortOpen  = ref(false);
const currentSort = ref('default');
const sortRef     = ref(null);

const sortOptions = {
  default:  'По умолчанию',
  dateDesc: 'Сначала новые',
  dateAsc:  'Сначала старые',
  nameAsc:  'По названию (А-Я)',
  nameDesc: 'По названию (Я-А)',
};

const selectSort = (option) => {
  currentSort.value = option;
  isSortOpen.value  = false;
};

const handleClickOutside = (e) => {
  if (sortRef.value && !sortRef.value.contains(e.target)) {
    isSortOpen.value = false;
  }
};

onMounted(async () => {
  document.addEventListener('mousedown', handleClickOutside);
  document.addEventListener('keydown', handleKeydown);

  try {
    paintings.value = await api.getPaintings();
  } catch (e) {
    loadError.value = 'Не удалось загрузить картины.';
  } finally {
    isLoading.value = false;
  }
});

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside);
  document.removeEventListener('keydown', handleKeydown);
});

const filteredPaintings = computed(() => {
  let result = paintings.value;

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(p => p.title.toLowerCase().includes(q));
  }
  if (selectedStyle.value)
    result = result.filter(p => (p.styles ?? []).some(s => s.name === selectedStyle.value));
  if (selectedSubject.value)
    result = result.filter(p => (p.plots ?? []).some(pl => pl.name === selectedSubject.value));
  if (selectedAuthor.value)
    result = result.filter(p => (p.authors ?? []).some(a => authorName(a) === selectedAuthor.value));
  if (selectedTechnique.value)
    result = result.filter(p => p.material?.name === selectedTechnique.value);
  if (yearFrom.value)
    result = result.filter(p => p.year != null && p.year >= Number(yearFrom.value));
  if (yearTo.value)
    result = result.filter(p => p.year != null && p.year <= Number(yearTo.value));

  switch (currentSort.value) {
    case 'dateDesc': return [...result].sort((a, b) => (b.year ?? 0) - (a.year ?? 0));
    case 'dateAsc':  return [...result].sort((a, b) => (a.year ?? 0) - (b.year ?? 0));
    case 'nameAsc':  return [...result].sort((a, b) => a.title.localeCompare(b.title, 'ru'));
    case 'nameDesc': return [...result].sort((a, b) => b.title.localeCompare(a.title, 'ru'));
    default:         return result;
  }
});
</script>

<template>
  <div class="container">

    <div class="filter-controls">
      <button class="toggle-btn" @click="isFilterOpen = !isFilterOpen">
        {{ isFilterOpen ? 'Скрыть фильтры' : 'Показать фильтры' }}
        <img src="../assets/icons/arrow.svg" alt="" :class="['arrow-icon', { rotate: isFilterOpen }]">
      </button>

      <div class="sort-box" ref="sortRef">
        <div class="sort-dropdown">
          <button class="sort-trigger" @click="isSortOpen = !isSortOpen">
            {{ sortOptions[currentSort] }}
            <img src="../assets/icons/arrow.svg" alt="" :class="['sort-arrow', { rotate: isSortOpen }]">
          </button>

          <transition name="fade">
            <div v-if="isSortOpen" class="sort-menu">
              <button
                v-for="(label, key) in sortOptions"
                :key="key"
                @click="selectSort(key)"
                :class="['sort-option', { active: currentSort === key }]"
              >
                {{ label }}
              </button>
            </div>
          </transition>
        </div>
      </div>
    </div>

    <transition name="expand">
      <div v-if="isFilterOpen" class="filter-panel">
        <div class="filter-grid">
          <div class="filter-group">
            <label>Название</label>
            <input type="text" v-model="searchQuery" placeholder="Поиск...">
          </div>

          <div class="filter-group">
            <label>Стиль</label>
            <select v-model="selectedStyle">
              <option value="">Все стили</option>
              <option v-for="s in styles" :key="s" :value="s">{{ s }}</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Сюжет</label>
            <select v-model="selectedSubject">
              <option value="">Все сюжеты</option>
              <option v-for="sub in subjects" :key="sub" :value="sub">{{ sub }}</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Автор</label>
            <select v-model="selectedAuthor">
              <option value="">Все авторы</option>
              <option v-for="a in authors" :key="a" :value="a">{{ a }}</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Техника</label>
            <select v-model="selectedTechnique">
              <option value="">Все техники</option>
              <option v-for="t in techniques" :key="t" :value="t">{{ t }}</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Год создания</label>
            <div class="year-range">
              <input type="number" v-model="yearFrom" placeholder="От" min="1000" max="2099">
              <input type="number" v-model="yearTo"   placeholder="До" min="1000" max="2099">
            </div>
          </div>
        </div>

        <button v-if="hasActiveFilters" class="reset-btn" @click="resetFilters">
          Сбросить фильтры
        </button>
      </div>
    </transition>

    <section class="gallery">
      <p v-if="isLoading" class="no-results">Загрузка...</p>
      <p v-else-if="loadError" class="no-results">{{ loadError }}</p>
      <p v-else-if="filteredPaintings.length === 0" class="no-results">
        Ничего не найдено. Попробуйте изменить фильтры.
      </p>

      <div v-for="painting in filteredPaintings" :key="painting.id" class="art-card" @click="openModal(painting)">
        <div class="image-wrapper">
          <img v-if="painting.image_path" :src="BASE + painting.image_path" :alt="painting.title" class="art-image">
        </div>
        <div class="art-info">
          <h3 class="art-title">{{ painting.title }}</h3>
          <p class="art-artist">{{ primaryAuthor(painting) }}</p>
        </div>
      </div>
    </section>

  </div>

  <PaintingModal
    v-if="selectedPainting"
    :painting="selectedPainting"
    @close="closeModal"
  />
</template>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 50px;
}

.filter-controls {
  padding: 30px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toggle-btn {
  background: none;
  border: none;
  color: #000;
  font-family: 'Raleway', sans-serif;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 2px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
}

.arrow-icon {
  width: 16px;
  height: 16px;
  transition: transform 0.2s ease-in-out;
}
.arrow-icon.rotate { transform: rotate(180deg); }

.sort-box {
  display: flex;
  align-items: center;
}

.sort-dropdown {
  position: relative;
}

.sort-trigger {
  background: none;
  border: none;
  color: #000;
  font-family: 'Raleway', sans-serif;
  font-size: 0.8rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 0;
  text-transform: uppercase;
  font-weight: 500;
  letter-spacing: 1px;
}

.sort-arrow {
  width: 16px;
  height: 16px;
  transition: transform 0.2s ease;
}
.sort-arrow.rotate { transform: rotate(180deg); }

.sort-menu {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  z-index: 100;
  background: #ffffff;
  border: 1px solid #dedede;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.12);
  min-width: 180px;
  display: flex;
  flex-direction: column;
  padding: 8px 0;
}

.sort-option {
  background: #ffffff;
  border: none;
  padding: 10px 20px;
  text-align: left;
  font-family: 'Raleway', sans-serif;
  font-size: 0.85rem;
  cursor: pointer;
  color: #666;
  transition: background 0.2s, color 0.2s;
}
.sort-option:hover  { background: #ebebeb; color: #000; }
.sort-option.active { background: #f5f5f5; color: #000; font-weight: 600; }

.year-range {
  display: flex;
  gap: 8px;
}

.year-range input {
  width: 50%;
}

.reset-btn {
  margin-top: 20px;
  background: none;
  border: none;
  color: #666;
  font-family: 'Raleway', sans-serif;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
}
.reset-btn:hover { color: #000; }

.fade-enter-active,
.fade-leave-active  { transition: opacity 0.2s, transform 0.2s; }
.fade-enter-from,
.fade-leave-to      { opacity: 0; transform: translateY(-6px); }

.filter-panel {
  padding: 20px 0;
  margin-bottom: 20px;
  border-bottom: 1px solid #dedede;
  overflow: hidden;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 30px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  color: #000;
  font-family: 'Raleway', sans-serif;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.filter-group input,
.filter-group select {
  background: #ebebeb;
  border: 1px solid transparent;
  border-radius: 2px;
  color: #000;
  padding: 10px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s;
}
.filter-group input:focus,
.filter-group select:focus { border-color: #444; }

.expand-enter-active,
.expand-leave-active {
  transition: max-height 0.4s ease, opacity 0.3s ease, padding 0.3s ease;
  max-height: 400px;
}
.expand-enter-from,
.expand-leave-to {
  max-height: 0;
  opacity: 0;
  padding: 0;
}

.gallery {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 30px;
  padding: 0;
}

.no-results {
  grid-column: 1 / -1;
  text-align: center;
  color: #666;
  padding: 60px 0;
  font-size: 0.9rem;
  letter-spacing: 1px;
}

.art-card {
  cursor: pointer;
}

.image-wrapper {
  position: relative;
  aspect-ratio: 4 / 5;
  background: #1a1a1a;
  margin-bottom: 15px;
  overflow: hidden;
}

.image-wrapper::after {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.15);
  opacity: 0;
  transition: opacity 0.2s ease;
  pointer-events: none;
}

.art-card:hover .image-wrapper::after { opacity: 1; }

.art-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.art-title {
  color: #000;
  font-size: 1rem;
  margin: 0 0 5px;
}

.art-artist {
  color: #666;
  font-size: 0.85rem;
  margin: 0;
}

@media (max-width: 1100px) {
  .gallery { grid-template-columns: repeat(3, 1fr); }
}

@media (max-width: 750px) {
  .container { padding: 0 20px; }
  .gallery   { grid-template-columns: repeat(2, 1fr); gap: 16px; }
}

@media (max-width: 480px) {
  .gallery { grid-template-columns: 1fr; }

  .filter-controls {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
}

</style>