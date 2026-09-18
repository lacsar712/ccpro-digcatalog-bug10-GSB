<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">探方 / 发掘单位</h2>
        <p class="page-sub">登记探方编号、深度与地层简述</p>
      </div>
      <button class="btn" @click="openCreate">新增探方</button>
    </div>

    <div class="card">
      <label style="max-width: 260px; margin-bottom: 1rem;">
        按工地筛选
        <select v-model="filterSiteId" @change="load">
          <option value="">全部工地</option>
          <option v-for="s in sites" :key="s.id" :value="String(s.id)">{{ s.name }}</option>
        </select>
      </label>
      <table class="table">
        <thead>
          <tr>
            <th>编号</th>
            <th>所属工地</th>
            <th>深度区间(m)</th>
            <th>地层简述</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td>{{ item.code }}</td>
            <td>{{ item.site?.name || '-' }}</td>
            <td>{{ item.depthMin }} ~ {{ item.depthMax }}</td>
            <td>{{ item.stratumDesc || '-' }}</td>
            <td>
              <button class="btn secondary small" @click="openEdit(item)">编辑</button>
              <button class="btn danger small" @click="remove(item)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-if="!list.length" class="page-sub">暂无数据</p>
      <p v-if="error" class="error">{{ error }}</p>
    </div>

    <div v-if="showModal" class="modal-mask" @click.self="showModal = false">
      <div class="modal">
        <h3>{{ form.id ? '编辑探方' : '新增探方' }}</h3>
        <div class="form-grid">
          <label>
            所属工地
            <select v-model.number="form.siteId">
              <option :value="0" disabled>请选择</option>
              <option v-for="s in sites" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
          </label>
          <label>
            编号
            <input v-model="form.code" placeholder="如 T1" />
          </label>
          <label>
            深度下限(m)
            <input v-model.number="form.depthMin" type="number" step="0.1" />
          </label>
          <label>
            深度上限(m)
            <input v-model.number="form.depthMax" type="number" step="0.1" />
          </label>
          <label class="full">
            地层简述
            <textarea v-model="form.stratumDesc" />
          </label>
        </div>
        <p v-if="formError" class="error">{{ formError }}</p>
        <div class="modal-actions">
          <button class="btn secondary" @click="showModal = false">取消</button>
          <button class="btn" @click="save">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import api from '../api/http'

const list = ref([])
const sites = ref([])
const filterSiteId = ref('')
const error = ref('')
const formError = ref('')
const showModal = ref(false)
const form = reactive({
  id: null,
  siteId: 0,
  code: '',
  depthMin: 0,
  depthMax: 0,
  stratumDesc: ''
})

async function loadSites() {
  const { data } = await api.get('/sites')
  sites.value = data
}

async function load() {
  error.value = ''
  try {
    const params = {}
    if (filterSiteId.value) params.siteId = filterSiteId.value
    const { data } = await api.get('/units', { params })
    list.value = data
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function openCreate() {
  Object.assign(form, {
    id: null,
    siteId: sites.value[0]?.id || 0,
    code: '',
    depthMin: 0,
    depthMax: 1,
    stratumDesc: ''
  })
  formError.value = ''
  showModal.value = true
}

function openEdit(item) {
  Object.assign(form, {
    id: item.id,
    siteId: item.siteId,
    code: item.code,
    depthMin: item.depthMin,
    depthMax: item.depthMax,
    stratumDesc: item.stratumDesc
  })
  formError.value = ''
  showModal.value = true
}

async function save() {
  formError.value = ''
  try {
    const payload = {
      siteId: form.siteId,
      code: form.code,
      depthMin: Number(form.depthMin) || 0,
      depthMax: Number(form.depthMax) || 0,
      stratumDesc: form.stratumDesc
    }
    if (form.id) {
      await api.put(`/units/${form.id}`, payload)
    } else {
      await api.post('/units', payload)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || '保存失败'
  }
}

async function remove(item) {
  if (!confirm(`确认删除探方「${item.code}」？`)) return
  try {
    await api.delete(`/units/${item.id}`)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(async () => {
  await loadSites()
  await load()
})
</script>
