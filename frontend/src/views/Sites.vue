<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">发掘工地</h2>
        <p class="page-sub">管理发掘工地基本信息与坐标</p>
      </div>
      <button class="btn" @click="openCreate">新增工地</button>
    </div>

    <div class="card">
      <table class="table">
        <thead>
          <tr>
            <th>名称</th>
            <th>时代</th>
            <th>经纬度</th>
            <th>负责人</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td>{{ item.name }}</td>
            <td><span class="tag">{{ item.period }}</span></td>
            <td>{{ item.latitude }}, {{ item.longitude }}</td>
            <td>{{ item.manager || '-' }}</td>
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
        <h3>{{ form.id ? '编辑工地' : '新增工地' }}</h3>
        <div class="form-grid">
          <label class="full">
            名称
            <input v-model="form.name" />
          </label>
          <label>
            时代
            <select v-model="form.period">
              <option value="新石器">新石器</option>
              <option value="夏商周">夏商周</option>
              <option value="商周">商周</option>
              <option value="秦汉">秦汉</option>
              <option value="其他">其他</option>
            </select>
          </label>
          <label>
            负责人
            <input v-model="form.manager" />
          </label>
          <label>
            纬度
            <input v-model.number="form.latitude" type="number" step="0.0001" />
          </label>
          <label>
            经度
            <input v-model.number="form.longitude" type="number" step="0.0001" />
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
const error = ref('')
const formError = ref('')
const showModal = ref(false)
const form = reactive({
  id: null,
  name: '',
  period: '新石器',
  latitude: 0,
  longitude: 0,
  manager: ''
})

async function load() {
  error.value = ''
  try {
    const { data } = await api.get('/sites')
    list.value = data
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function openCreate() {
  Object.assign(form, { id: null, name: '', period: '新石器', latitude: 0, longitude: 0, manager: '' })
  formError.value = ''
  showModal.value = true
}

function openEdit(item) {
  Object.assign(form, item)
  formError.value = ''
  showModal.value = true
}

async function save() {
  formError.value = ''
  try {
    const payload = {
      name: form.name,
      period: form.period,
      latitude: Number(form.latitude) || 0,
      longitude: Number(form.longitude) || 0,
      manager: form.manager
    }
    if (form.id) {
      await api.put(`/sites/${form.id}`, payload)
    } else {
      await api.post('/sites', payload)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || '保存失败'
  }
}

async function remove(item) {
  if (!confirm(`确认删除工地「${item.name}」？`)) return
  try {
    await api.delete(`/sites/${item.id}`)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(load)
</script>
