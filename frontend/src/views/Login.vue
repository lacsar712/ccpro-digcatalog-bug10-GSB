<template>
  <div class="login-page">
    <div class="login-card card">
      <h1>考古发掘出土文物编目系统</h1>
      <p>登录后管理工地、探方与出土文物登记信息</p>
      <form @submit.prevent="onSubmit">
        <label>
          用户名
          <input v-model="username" autocomplete="username" placeholder="admin / recorder" />
        </label>
        <label>
          密码
          <input v-model="password" type="password" autocomplete="current-password" placeholder="123456" />
        </label>
        <p v-if="error" class="error">{{ error }}</p>
        <button class="btn" type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const username = ref('admin')
const password = ref('123456')
const error = ref('')
const loading = ref(false)

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(username.value.trim(), password.value)
    router.push({ name: 'overview' })
  } catch (e) {
    error.value = e.response?.data?.error || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 1rem;
  background:
    linear-gradient(135deg, rgba(44, 36, 25, 0.88), rgba(44, 36, 25, 0.55)),
    radial-gradient(circle at 20% 20%, #c4a574, transparent 45%),
    #2c2419;
}

.login-card {
  width: min(420px, 100%);
}

h1 {
  margin: 0 0 0.4rem;
  font-size: 1.35rem;
}

p {
  margin: 0 0 1.2rem;
  color: var(--muted);
}

form {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

button {
  width: 100%;
}
</style>
