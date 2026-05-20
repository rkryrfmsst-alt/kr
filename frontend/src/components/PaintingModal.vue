<script setup>
import { computed } from 'vue';
import { BASE } from '../api/index.js';

const props = defineProps({
  painting: { type: Object, required: true },
});

const emit = defineEmits(['close']);

const authorName = (a) => [a.last_name, a.first_name, a.middle_name].filter(Boolean).join(' ');

const artists   = computed(() => (props.painting.authors ?? []).map(authorName).join(', ') || '—');
const styles    = computed(() => (props.painting.styles  ?? []).map(s => s.name).join(', ') || '—');
const plots     = computed(() => (props.painting.plots   ?? []).map(p => p.name).join(', ') || '—');
const technique = computed(() => props.painting.material?.name ?? '—');
</script>

<template>
  <teleport to="body">
    <transition name="modal-fade">
      <div class="overlay" @click.self="emit('close')">
        <div class="modal">

          <div class="image-pane">
            <div class="passepartout">
              <img
                v-if="painting.image_path"
                :src="BASE + painting.image_path"
                :alt="painting.title"
                class="painting-img"
              >
              <div v-else class="painting-placeholder"></div>
            </div>
          </div>

          <div class="info-pane">
            <button class="close-btn" @click="emit('close')" aria-label="Закрыть">&#x2715;</button>

            <div class="info-body">
              <h2 class="info-title">{{ painting.title }}</h2>
              <p class="info-artist">{{ artists }}</p>

              <dl class="meta">
                <div class="meta-row">
                  <dt>Стиль</dt>
                  <dd>{{ styles }}</dd>
                </div>
                <div class="meta-row">
                  <dt>Техника</dt>
                  <dd>{{ technique }}</dd>
                </div>
                <div class="meta-row">
                  <dt>Сюжет</dt>
                  <dd>{{ plots }}</dd>
                </div>
                <div class="meta-row">
                  <dt>Год</dt>
                  <dd>{{ painting.year ?? '—' }}</dd>
                </div>
              </dl>
            </div>
          </div>

        </div>
      </div>
    </transition>
  </teleport>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.72);
  z-index: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
}

.modal {
  display: grid;
  grid-template-columns: 62% 38%;
  width: 100%;
  max-width: 1060px;
  max-height: calc(100vh - 64px);
  background: #fff;
  overflow: hidden;
  border-radius: 6px;
}

.image-pane {
  background: #fff;
  padding: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.passepartout {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.painting-placeholder,
.painting-img {
  width: 100%;
  aspect-ratio: 4 / 5;
  max-height: calc(100vh - 200px);
  background: #1a1a1a;
  box-shadow:
    0 0 0 1px #d0d0d0,
    0 8px 40px rgba(0, 0, 0, 0.18);
}

.painting-img {
  object-fit: contain;
  display: block;
}

.info-pane {
  background: #fff;
  border-left: 1px solid #e8e8e8;
  position: relative;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  background: none;
  border: none;
  font-size: 1rem;
  color: #aaa;
  cursor: pointer;
  line-height: 1;
  padding: 4px;
  transition: color 0.15s;
}
.close-btn:hover { color: #000; }

.info-body {
  padding: 52px 40px 48px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.info-title {
  font-size: 1.55rem;
  font-weight: 600;
  color: #000;
  margin: 0 0 6px;
  line-height: 1.25;
}

.info-artist {
  font-size: 0.95rem;
  color: #666;
  margin: 0 0 44px;
}

.meta {
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.meta-row {
  display: grid;
  grid-template-columns: 90px 1fr;
  align-items: baseline;
  padding: 14px 0;
  border-top: 1px solid #ebebeb;
}
.meta-row:last-child { border-bottom: 1px solid #ebebeb; }

.meta-row dt {
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #666;
}

.meta-row dd {
  margin: 0;
  font-size: 0.9rem;
  color: #111;
}

.modal-fade-enter-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}
.modal-fade-leave-active {
  transition: opacity 0.18s ease;
}
.modal-fade-enter-from {
  opacity: 0;
  transform: scale(0.97);
}
.modal-fade-leave-to {
  opacity: 0;
}

@media (max-width: 750px) {
  .overlay { padding: 0; }

  .modal {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
    max-width: 100%;
    max-height: 100vh;
    height: 100vh;
  }

  .image-pane {
    padding: 24px 24px 16px;
  }

  .painting-placeholder {
    aspect-ratio: 3 / 2;
    max-height: 42vh;
  }

  .info-pane { border-left: none; border-top: 1px solid #e8e8e8; }

  .info-body { padding: 32px 24px 32px; }

  .info-title { font-size: 1.25rem; }
  .info-artist { margin-bottom: 28px; }
}
</style>
