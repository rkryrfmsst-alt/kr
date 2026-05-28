<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import PaintingModal from '../components/PaintingModal.vue';
import { api, BASE } from '../api/index.js';
import { useAuth } from '../composables/useAuth.js';

const route = useRoute();
const { adminMode } = useAuth();

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
  yearFrom.value          = minYear.value;
  yearTo.value            = maxYear.value;
};

const hasActiveFilters = computed(() =>
  searchQuery.value || selectedStyle.value || selectedSubject.value ||
  selectedAuthor.value || selectedTechnique.value ||
  (yearFrom.value !== '' && Number(yearFrom.value) !== minYear.value) ||
  (yearTo.value   !== '' && Number(yearTo.value)   !== maxYear.value)
);

const styles = computed(() =>
  [...new Set(paintings.value.flatMap(p => (p.styles ?? []).map(s => s.name)))].sort()
);
const subjects = computed(() =>
  [...new Set(paintings.value.flatMap(p => (p.plots ?? []).map(pl => pl.name)))].sort()
);
const allAuthors   = ref([]);
const allMaterials = ref([]);
const allStyles    = ref([]);
const allPlots     = ref([]);
const authors = computed(() =>
  allAuthors.value.map(authorName).sort()
);
const techniques = computed(() =>
  [...new Set(paintings.value.filter(p => p.material).map(p => p.material.name))].sort()
);

const minYear = computed(() => {
  const years = paintings.value.map(p => p.year).filter(y => y != null);
  return years.length ? Math.min(...years) : '';
});
const maxYear = computed(() => {
  const years = paintings.value.map(p => p.year).filter(y => y != null);
  return years.length ? Math.max(...years) : '';
});

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
    [paintings.value, allAuthors.value, allMaterials.value, allStyles.value, allPlots.value] = await Promise.all([
      api.getPaintings(), api.getAuthors(), api.getMaterials(), api.getStyles(), api.getPlots(),
    ]);
    yearFrom.value  = minYear.value;
    yearTo.value    = maxYear.value;
    if (route.query.author) {
      selectedAuthor.value = route.query.author;
      isFilterOpen.value   = true;
    }
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

// --- Add painting modal ---
const showAddModal    = ref(false);
const addLoading      = ref(false);
const addError        = ref('');
const newTitle        = ref('');
const newYear         = ref('');
const newDescription  = ref('');
const newAuthorID     = ref('');
const newMaterialID   = ref('');
const newStyleID      = ref('');
const newPlotID       = ref('');
const newImageFile    = ref(null);
const imagePreview    = ref('');

function openAddModal() {
  newTitle.value       = '';
  newYear.value        = '';
  newDescription.value = '';
  newAuthorID.value    = '';
  newMaterialID.value  = '';
  newStyleID.value     = '';
  newPlotID.value      = '';
  newImageFile.value   = null;
  imagePreview.value   = '';
  addError.value       = '';
  showAddModal.value   = true;
}

function closeAddModal() {
  showAddModal.value = false;
}

function onImageChange(e) {
  const file = e.target.files[0];
  if (!file) return;
  newImageFile.value = file;
  imagePreview.value = URL.createObjectURL(file);
}

async function submitAddPainting() {
  addError.value = '';
  addLoading.value = true;
  try {
    const fd = new FormData();
    fd.append('title', newTitle.value.trim());
    if (newYear.value)        fd.append('year',        newYear.value);
    if (newDescription.value) fd.append('description', newDescription.value.trim());
    if (newMaterialID.value)  fd.append('material_id', newMaterialID.value);
    if (newAuthorID.value)    fd.append('author_ids[]', newAuthorID.value);
    if (newStyleID.value) fd.append('style_ids[]', newStyleID.value);
    if (newPlotID.value)  fd.append('plot_ids[]',  newPlotID.value);
    if (newImageFile.value)   fd.append('image',       newImageFile.value);

    const painting = await api.createPainting(fd);
    paintings.value.push(painting);
    closeAddModal();
  } catch (e) {
    addError.value = e.message;
  } finally {
    addLoading.value = false;
  }
}

// --- Edit painting modal ---
const showEditModal      = ref(false);
const editPaintingID     = ref(null);
const editLoading        = ref(false);
const editError          = ref('');
const editTitle          = ref('');
const editYear           = ref('');
const editDescription    = ref('');
const editAuthorID       = ref('');
const editMaterialID     = ref('');
const editStyleID        = ref('');
const editPlotID         = ref('');
const editImageFile      = ref(null);
const editImagePreview   = ref('');

function openEditModal(painting) {
  editPaintingID.value   = painting.id;
  editTitle.value        = painting.title;
  editYear.value         = painting.year ?? '';
  editDescription.value  = painting.description ?? '';
  editAuthorID.value     = painting.authors?.[0]?.id ?? '';
  editMaterialID.value   = painting.material?.id ?? '';
  editStyleID.value      = painting.styles?.[0]?.id ?? '';
  editPlotID.value       = painting.plots?.[0]?.id ?? '';
  editImageFile.value    = null;
  editImagePreview.value = painting.image_path ? BASE + painting.image_path : '';
  editError.value        = '';
  showEditModal.value    = true;
}

function closeEditModal() {
  showEditModal.value = false;
}

function onEditImageChange(e) {
  const file = e.target.files[0];
  if (!file) return;
  editImageFile.value    = file;
  editImagePreview.value = URL.createObjectURL(file);
}

async function submitEditPainting() {
  editError.value   = '';
  editLoading.value = true;
  try {
    const fd = new FormData();
    fd.append('title', editTitle.value.trim());
    if (editYear.value)        fd.append('year',         editYear.value);
    if (editDescription.value) fd.append('description',  editDescription.value.trim());
    if (editMaterialID.value)  fd.append('material_id',  editMaterialID.value);
    if (editAuthorID.value)    fd.append('author_ids[]', editAuthorID.value);
    if (editStyleID.value)     fd.append('style_ids[]',  editStyleID.value);
    if (editPlotID.value)      fd.append('plot_ids[]',   editPlotID.value);
    if (editImageFile.value)   fd.append('image',        editImageFile.value);

    const updated = await api.updatePainting(editPaintingID.value, fd);
    const idx = paintings.value.findIndex(p => p.id === editPaintingID.value);
    if (idx !== -1) paintings.value[idx] = updated;
    closeEditModal();
  } catch (e) {
    editError.value = e.message;
  } finally {
    editLoading.value = false;
  }
}

async function deletePainting(painting) {
  try {
    await api.deletePainting(painting.id);
    paintings.value = paintings.value.filter(p => p.id !== painting.id);
  } catch (e) {
    loadError.value = e.message;
  }
}

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

      <div class="filter-controls-right">
      <button v-if="adminMode" class="add-painting-btn" @click="openAddModal">+ Добавить картину</button>
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
          <div v-if="adminMode" class="card-admin-overlay">
            <button class="card-admin-btn" @click.stop="openEditModal(painting)">Изменить</button>
            <button class="card-admin-btn card-admin-btn--del" @click.stop="deletePainting(painting)">Удалить</button>
          </div>
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

  <teleport to="body">
    <div v-if="showEditModal" class="ap-backdrop" @click.self="closeEditModal">
      <div class="ap-modal">
        <h2 class="ap-title">Редактировать картину</h2>
        <form class="ap-form" novalidate @submit.prevent="submitEditPainting">
          <div class="ap-field">
            <label class="ap-label">Название <span style="color:#c00;font-family:Arial,sans-serif;font-size:0.9rem;">*</span></label>
            <input v-model="editTitle" class="ap-input" type="text" placeholder="Название работы" />
          </div>
          <div class="ap-row">
            <div class="ap-field">
              <label class="ap-label">Автор</label>
              <select v-model="editAuthorID" class="ap-input ap-select">
                <option value="">— не указан —</option>
                <option v-for="a in allAuthors" :key="a.id" :value="a.id">{{ authorName(a) }}</option>
              </select>
            </div>
            <div class="ap-field">
              <label class="ap-label">Год</label>
              <input v-model="editYear" class="ap-input" type="number" placeholder="2024" min="1000" max="2099" />
            </div>
          </div>
          <div class="ap-field">
            <label class="ap-label">Техника</label>
            <select v-model="editMaterialID" class="ap-input ap-select">
              <option value="">— не указана —</option>
              <option v-for="m in allMaterials" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>
          <div class="ap-row">
            <div class="ap-field">
              <label class="ap-label">Стиль</label>
              <select v-model="editStyleID" class="ap-input ap-select">
                <option value="">— не указан —</option>
                <option v-for="s in allStyles" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
            <div class="ap-field">
              <label class="ap-label">Сюжет</label>
              <select v-model="editPlotID" class="ap-input ap-select">
                <option value="">— не указан —</option>
                <option v-for="p in allPlots" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
          </div>
          <div class="ap-field">
            <label class="ap-label">Описание</label>
            <textarea v-model="editDescription" class="ap-input ap-textarea" placeholder="Краткое описание..." rows="3"></textarea>
          </div>
          <div class="ap-field">
            <label class="ap-label">Изображение (оставьте пустым, чтобы не менять)</label>
            <label class="ap-file-label">
              <input type="file" accept="image/*" class="ap-file-input" @change="onEditImageChange" />
              <span>{{ editImageFile ? editImageFile.name : 'Выбрать новый файл...' }}</span>
            </label>
            <img v-if="editImagePreview" :src="editImagePreview" class="ap-preview" alt="preview" />
          </div>
          <p v-if="editError" class="ap-error">{{ editError }}</p>
          <div class="ap-actions">
            <button type="button" class="ap-cancel" @click="closeEditModal">Отмена</button>
            <button type="submit" class="ap-submit" :disabled="editLoading">
              {{ editLoading ? 'Сохранение...' : 'Сохранить' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </teleport>

  <teleport to="body">
    <div v-if="showAddModal" class="ap-backdrop" @click.self="closeAddModal">
      <div class="ap-modal">
        <h2 class="ap-title">Новая картина</h2>

        <form class="ap-form" novalidate @submit.prevent="submitAddPainting">

          <div class="ap-field">
            <label class="ap-label">Название <span style="color:#c00;font-family:Arial,sans-serif;font-size:0.9rem;">*</span></label>
            <input v-model="newTitle" class="ap-input" type="text" placeholder="Название работы" />
          </div>

          <div class="ap-row">
            <div class="ap-field">
              <label class="ap-label">Автор</label>
              <select v-model="newAuthorID" class="ap-input ap-select">
                <option value="">— не указан —</option>
                <option v-for="a in allAuthors" :key="a.id" :value="a.id">{{ authorName(a) }}</option>
              </select>
            </div>
            <div class="ap-field">
              <label class="ap-label">Год</label>
              <input v-model="newYear" class="ap-input" type="number" placeholder="2024" min="1000" max="2099" />
            </div>
          </div>

          <div class="ap-field">
            <label class="ap-label">Техника</label>
            <select v-model="newMaterialID" class="ap-input ap-select">
              <option value="">— не указана —</option>
              <option v-for="m in allMaterials" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>

          <div class="ap-row">
            <div class="ap-field">
              <label class="ap-label">Стиль</label>
              <select v-model="newStyleID" class="ap-input ap-select">
                <option value="">— не указан —</option>
                <option v-for="s in allStyles" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
            <div class="ap-field">
              <label class="ap-label">Сюжет</label>
              <select v-model="newPlotID" class="ap-input ap-select">
                <option value="">— не указан —</option>
                <option v-for="p in allPlots" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
          </div>

          <div class="ap-field">
            <label class="ap-label">Описание</label>
            <textarea v-model="newDescription" class="ap-input ap-textarea" placeholder="Краткое описание..." rows="3"></textarea>
          </div>

          <div class="ap-field">
            <label class="ap-label">Изображение</label>
            <label class="ap-file-label">
              <input type="file" accept="image/*" class="ap-file-input" @change="onImageChange" />
              <span>{{ newImageFile ? newImageFile.name : 'Выбрать файл...' }}</span>
            </label>
            <img v-if="imagePreview" :src="imagePreview" class="ap-preview" alt="preview" />
          </div>

          <p v-if="addError" class="ap-error">{{ addError }}</p>

          <div class="ap-actions">
            <button type="button" class="ap-cancel" @click="closeAddModal">Отмена</button>
            <button type="submit" class="ap-submit" :disabled="addLoading">
              {{ addLoading ? 'Сохранение...' : 'Добавить' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </teleport>
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
  gap: 16px;
}

.filter-controls-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.add-painting-btn {
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
.add-painting-btn:hover { background: #222; }

.toggle-btn {
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
  box-sizing: border-box;
  width: 100%;
  height: 42px;
}
.filter-group input:focus,
.filter-group select:focus { border-color: #444; }

.filter-group select {
  -webkit-appearance: none;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' viewBox='0 0 12 8'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%23666' stroke-width='1.5' fill='none' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 32px;
  cursor: pointer;
}

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

.card-admin-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: flex-end;
  gap: 6px;
  padding: 10px;
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 2;
}
.image-wrapper:hover .card-admin-overlay { opacity: 1; }

.card-admin-btn {
  flex: 1;
  background: rgba(255,255,255,0.92);
  border: none;
  border-radius: 2px;
  padding: 7px 0;
  font-family: 'Raleway', sans-serif;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #000;
  cursor: pointer;
  transition: background 0.15s;
}
.card-admin-btn:hover { background: #fff; }
.card-admin-btn--del { color: #c00; }
.card-admin-btn--del:hover { background: #fff; }

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

/* ---- Add painting modal ---- */
.ap-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 20px;
}
.ap-modal {
  background: #fff;
  width: 100%;
  max-width: 700px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 40px 48px;
  border-radius: 4px;
}
.ap-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #000;
  margin: 0 0 28px;
}
.ap-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.ap-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.ap-field {
  display: flex;
  flex-direction: column;
  gap: 7px;
}
.ap-label {
  font-family: 'Raleway', sans-serif;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #000;
}
.ap-input {
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
.ap-input:focus { border-color: #444; }
.ap-input::placeholder { color: #aaa; }
.ap-select {
  -webkit-appearance: none;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' viewBox='0 0 12 8'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%23666' stroke-width='1.5' fill='none' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 32px;
  cursor: pointer;
}
.ap-textarea {
  resize: none;
  min-height: 80px;
}
.ap-file-label {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #ebebeb;
  border: 1px solid transparent;
  border-radius: 2px;
  padding: 11px 14px;
  cursor: pointer;
  font-family: 'Inter', sans-serif;
  font-size: 0.85rem;
  color: #666;
  transition: border-color 0.2s;
}
.ap-file-label:hover { border-color: #bbb; }
.ap-file-input {
  display: none;
}
.ap-preview {
  width: 100%;
  max-height: 180px;
  object-fit: contain;
  border: 1px solid #e8e8e8;
  border-radius: 2px;
  margin-top: 4px;
}
.ap-error {
  font-size: 0.8rem;
  color: #c00;
  margin: 0;
}
.ap-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}
.ap-cancel {
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
.ap-cancel:hover { border-color: #999; color: #000; }
.ap-submit {
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
.ap-submit:hover:not(:disabled) { background: #222; }
.ap-submit:disabled { opacity: 0.5; cursor: default; }

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