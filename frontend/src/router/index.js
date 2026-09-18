import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Login from '../views/Login.vue'
import Layout from '../components/Layout.vue'
import Overview from '../views/Overview.vue'
import Sites from '../views/Sites.vue'
import Units from '../views/Units.vue'
import Finds from '../views/Finds.vue'
import Materials from '../views/Materials.vue'

const routes = [
  { path: '/login', name: 'login', component: Login, meta: { public: true } },
  {
    path: '/',
    component: Layout,
    children: [
      { path: '', name: 'overview', component: Overview },
      { path: 'sites', name: 'sites', component: Sites },
      { path: 'units', name: 'units', component: Units },
      { path: 'finds', name: 'finds', component: Finds },
      { path: 'materials', name: 'materials', component: Materials }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: 'login' }
  }
  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'overview' }
  }
})

export default router
