<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">Dig</div>
        <div>
          <div class="brand-title">DigCatalog</div>
          <div class="brand-sub">出土文物编目</div>
        </div>
      </div>
      <nav>
        <router-link v-for="item in menus" :key="item.to" :to="item.to" class="nav-item">
          {{ item.label }}
        </router-link>
      </nav>
      <div class="sidebar-foot">
        <div>{{ auth.user?.username }} · {{ auth.roleLabel }}</div>
        <button class="btn secondary small" @click="onLogout">退出登录</button>
      </div>
    </aside>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const menus = [
  { to: '/', label: '概览' },
  { to: '/sites', label: '发掘工地' },
  { to: '/units', label: '探方单位' },
  { to: '/finds', label: '出土文物' },
  { to: '/materials', label: '材质分类' }
]

function onLogout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  min-height: 100vh;
}

.sidebar {
  background: linear-gradient(180deg, #34291d 0%, var(--sidebar) 100%);
  color: var(--sidebar-text);
  padding: 1.25rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.brand {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  padding: 0.4rem 0.5rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  background: #c4a574;
  color: #2c2419;
  font-weight: 700;
}

.brand-title {
  font-weight: 700;
  letter-spacing: 0.02em;
}

.brand-sub {
  font-size: 0.78rem;
  opacity: 0.75;
}

nav {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  flex: 1;
}

.nav-item {
  padding: 0.7rem 0.85rem;
  border-radius: 10px;
  color: var(--sidebar-text);
  opacity: 0.85;
}

.nav-item:hover,
.nav-item.router-link-exact-active {
  background: rgba(196, 165, 116, 0.18);
  opacity: 1;
}

.sidebar-foot {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  font-size: 0.85rem;
  padding-top: 0.75rem;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.main {
  padding: 1.5rem;
}

@media (max-width: 900px) {
  .layout {
    grid-template-columns: 1fr;
  }
  .sidebar {
    position: sticky;
    top: 0;
    z-index: 10;
  }
  nav {
    flex-direction: row;
    flex-wrap: wrap;
  }
}
</style>
