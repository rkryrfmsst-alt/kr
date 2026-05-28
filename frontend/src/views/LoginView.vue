<script setup>
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuth } from '../composables/useAuth.js'

const router = useRouter()
const { login, isLoggedIn } = useAuth()

if (isLoggedIn.value) router.replace('/')

const email    = ref('')
const password = ref('')
const error    = ref('')
const loading  = ref(false)

async function submit() {
  error.value   = ''
  loading.value = true
  try {
    await login(email.value.trim(), password.value)
    router.replace('/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <div class="card">
      <h1 class="title">Войти</h1>
      <p class="subtitle">Добро пожаловать в галерею</p>

      <form class="form" novalidate @submit.prevent="submit">
        <div class="field">
          <label for="email" class="label">Электронная почта</label>
          <input
            id="email"
            v-model="email"
            type="email"
            class="input"
            placeholder="example@mail.com"
            autocomplete="email"
            required
          />
        </div>

        <div class="field">
          <label for="password" class="label">Пароль</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="input"
            placeholder="••••••••"
            autocomplete="current-password"
            required
          />
        </div>

        <p v-if="error" class="error">{{ error }}</p>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>
      </form>

      <p class="switch">
        Нет аккаунта?
        <RouterLink to="/register" class="switch-link">Регистрация</RouterLink>
      </p>
    </div>
  </div>
</template>

<style scoped>
.page {
  min-height: calc(100vh - 230px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.card {
  width: 100%;
  max-width: 420px;
}

.title {
  font-size: 1.75rem;
  font-weight: 600;
  color: #000;
  margin: 0 0 6px;
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 0.875rem;
  color: #aaa;
  margin: 0 0 40px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label {
  font-family: 'Raleway', sans-serif;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #000;
}

.input {
  background: #ebebeb;
  border: 1px solid transparent;
  border-radius: 2px;
  padding: 12px 14px;
  font-family: 'Inter', sans-serif;
  font-size: 0.9rem;
  color: #000;
  outline: none;
  transition: border-color 0.2s;
}
.input:focus { border-color: #444; }
.input::placeholder { color: #aaa; }

.error {
  font-size: 0.8rem;
  color: #c00;
  margin: 0;
}

.submit-btn {
  margin-top: 8px;
  background: #000;
  color: #fff;
  border: none;
  border-radius: 2px;
  padding: 14px;
  font-family: 'Raleway', sans-serif;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 2px;
  cursor: pointer;
  transition: background 0.2s;
}
.submit-btn:hover:not(:disabled) { background: #222; }
.submit-btn:disabled { opacity: 0.5; cursor: default; }

.switch {
  margin-top: 28px;
  font-size: 0.85rem;
  color: #888;
  text-align: center;
}

.switch-link {
  color: #000;
  font-weight: 500;
  text-decoration: none;
  border-bottom: 1px solid currentColor;
  transition: opacity 0.2s;
}
.switch-link:hover { opacity: 0.5; }
</style>
