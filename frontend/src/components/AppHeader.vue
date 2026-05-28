<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth.js'

const router = useRouter()
const { user, isLoggedIn, isAdmin, adminMode, logout, toggleAdminMode } = useAuth()

const isMenuOpen = ref(false)
const menuRef    = ref(null)

function handleClickOutside(e) {
  if (menuRef.value && !menuRef.value.contains(e.target)) {
    isMenuOpen.value = false
  }
}

onMounted(()  => document.addEventListener('mousedown', handleClickOutside))
onUnmounted(() => document.removeEventListener('mousedown', handleClickOutside))

async function handleLogout() {
  isMenuOpen.value = false
  await logout()
  router.push('/')
}
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
        <template v-if="isLoggedIn">
          <div class="profile-dropdown" ref="menuRef">
            <button class="profile-btn" @click="isMenuOpen = !isMenuOpen">
              <img src="../assets/icons/profile.svg" alt="profile" class="profile-icon">
            </button>

            <transition name="fade">
              <div v-if="isMenuOpen" class="profile-menu">
                <span class="menu-username">{{ user.username }}</span>
                <div class="menu-divider"></div>
                <template v-if="isAdmin">
                  <label class="menu-toggle-item">
                    <span>Режим редактора</span>
                    <div class="toggle" :class="{ active: adminMode }" @click.prevent="toggleAdminMode"></div>
                  </label>
                  <div class="menu-divider"></div>
                </template>
                <button class="menu-item" @click="handleLogout">Выйти</button>
              </div>
            </transition>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/login" class="auth-link">Войти</RouterLink>
          <span class="sep">/</span>
          <RouterLink to="/register" class="auth-link">Регистрация</RouterLink>
        </template>
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

.profile-dropdown {
  position: relative;
}

.profile-btn {
  background: none;
  border: none;
  padding: 7px;
  cursor: pointer;
  display: flex;
  align-items: center;
  border-radius: 50%;
  transition: background 0.2s;
}
.profile-btn:hover { background: #f0f0f0; }

.profile-icon {
  width: 26px;
  height: 26px;
  display: block;
}

.profile-menu {
  position: absolute;
  top: calc(100% + 14px);
  right: 0;
  background: #fff;
  border: 1px solid #dedede;
  border-radius: 6px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  min-width: 210px;
  padding: 8px 0;
  z-index: 100;
}

.menu-username {
  display: block;
  padding: 8px 18px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #000;
  white-space: nowrap;
}

.menu-divider {
  height: 1px;
  background: #ebebeb;
  margin: 6px 0;
}

.menu-item {
  display: block;
  width: 100%;
  background: none;
  border: none;
  padding: 8px 18px;
  text-align: left;
  font-family: 'Inter', sans-serif;
  font-size: 0.85rem;
  color: #666;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
.menu-item:hover { background: #f5f5f5; color: #000; }

.menu-toggle-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 18px;
  cursor: pointer;
  font-size: 0.85rem;
  color: #666;
  user-select: none;
  gap: 12px;
}
.menu-toggle-item:hover { color: #000; }

.toggle {
  width: 34px;
  height: 18px;
  border-radius: 9px;
  background: #ccc;
  position: relative;
  transition: background 0.2s;
  flex-shrink: 0;
}
.toggle::after {
  content: '';
  position: absolute;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #fff;
  top: 3px;
  left: 3px;
  transition: transform 0.2s;
}
.toggle.active { background: #000; }
.toggle.active::after { transform: translateX(16px); }

.fade-enter-active,
.fade-leave-active  { transition: opacity 0.15s, transform 0.15s; }
.fade-enter-from,
.fade-leave-to      { opacity: 0; transform: translateY(-4px); }
</style>
