import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api/http'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('digcatalog_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('digcatalog_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const roleLabel = computed(() => (user.value?.role === 'admin' ? '管理员' : '记录员'))

  async function login(username, password) {
    const { data } = await api.post('/auth/login', { username, password })
    token.value = data.token
    user.value = data.user
    localStorage.setItem('digcatalog_token', data.token)
    localStorage.setItem('digcatalog_user', JSON.stringify(data.user))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('digcatalog_token')
    localStorage.removeItem('digcatalog_user')
  }

  return { token, user, isLoggedIn, roleLabel, login, logout }
})
