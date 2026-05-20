<template>
  <div class="container">
    <div class="contacts">
      <h1 class="title">Контакты</h1>

      <ul class="list">
        <li class="item">
          <span class="label">Почта</span>
          <a href="mailto:gallery@example.com" class="link">gallery@example.com</a>
        </li>
        <li class="item">
          <span class="label">Telegram</span>
          <a href="https://t.me/gallery" target="_blank" rel="noopener" class="link">@gallery</a>
        </li>
        <li class="item">
          <span class="label">Режим работы</span>
          <span class="value">Пн — Пт, 10:00 — 19:00</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 50px;
}

.contacts {
  padding: 80px 0 100px;
}

.title {
  font-size: 2rem;
  font-weight: 600;
  color: #000;
  margin: 0 0 48px;
}

.list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.item {
  display: grid;
  grid-template-columns: 160px 1fr;
  align-items: baseline;
  padding: 20px 0;
  border-top: 1px solid #e8e8e8;
}
.item:last-child { border-bottom: 1px solid #e8e8e8; }

.label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #aaa;
}

.link {
  font-size: 0.95rem;
  color: #000;
  text-decoration: none;
  transition: opacity 0.15s;
}
.link:hover { opacity: 0.5; }

.value {
  font-size: 0.95rem;
  color: #000;
}

@media (max-width: 750px) {
  .container { padding: 0 20px; }
  .contacts { padding: 50px 0 70px; }
  .title { font-size: 1.5rem; margin-bottom: 32px; }
  .item { grid-template-columns: 130px 1fr; }
}
</style>
