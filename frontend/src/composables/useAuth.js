import { ref, computed } from 'vue'
import { api } from '../api/index.js'

const user      = ref(null)
const checked   = ref(false)
const adminMode = ref(localStorage.getItem('adminMode') === 'true')

export function useAuth() {
  const isLoggedIn = computed(() => user.value !== null)
  const isAdmin    = computed(() => user.value?.role === 'admin')

  async function fetchMe() {
    try {
      user.value = await api.me()
    } catch {
      user.value = null
    } finally {
      checked.value = true
    }
  }

  async function login(email, password) {
    user.value = await api.login(email, password)
  }

  async function register(email, username, password) {
    user.value = await api.register(email, username, password)
  }

  async function logout() {
    await api.logout()
    user.value = null
    adminMode.value = false
    localStorage.removeItem('adminMode')
  }

  function toggleAdminMode() {
    adminMode.value = !adminMode.value
    localStorage.setItem('adminMode', adminMode.value)
  }

  return { user, isLoggedIn, isAdmin, checked, adminMode, fetchMe, login, register, logout, toggleAdminMode }
}
