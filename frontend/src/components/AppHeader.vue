<script setup>
import { RouterLink } from 'vue-router';
</script>

<template>
  <header class="header">
    <div class="container">
      <RouterLink to="/" class="logo">
        <img src="../assets/icons/logo.svg" alt="logo" class="logo-icon">
      </RouterLink>

      <nav class="nav">
        <RouterLink to="/" class="nav-link">Галерея</RouterLink>
        <RouterLink to="/artists" class="nav-link">Художники</RouterLink>
        <RouterLink to="/about" class="nav-link">О проекте</RouterLink>
        <RouterLink to="/contacts" class="nav-link">Контакты</RouterLink>
      </nav>

      <div class="auth">
        <RouterLink to="/login" class="auth-link">Войти</RouterLink>
        <span class="sep">/</span>
        <RouterLink to="/register" class="auth-link">Регистрация</RouterLink>
      </div>
    </div>
  </header>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap');

.header {
  background: #fff;
  box-shadow: 0 0 32px rgba(0, 0, 0, 0.2);
  color: #000;
  font-family: 'Inter', sans-serif;
  width: 100%;
  justify-content: center;
  padding: 15px 0;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 50px;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
}

.logo-icon {
  width: 128px;
  height: 64px;
}

.nav {
  display: flex;
  gap: 50px;
  font-size: 15px;
}

.nav-link {
  color: #000;
  font-weight: 400;
  white-space: nowrap;
}

.auth {
  justify-self: end;
  display: flex;
  align-items: center;
  gap: 8px;
}

.auth-link {
  color: #000;
  font-weight: 400;
  text-decoration: none;
  font-size: 0.9rem;
  transition: opacity 0.2s;
}
.auth-link:hover { opacity: 0.7; }

.sep {
  color: rgba(0, 0, 0, 0.3);
  font-size: 20px;
}
</style>
