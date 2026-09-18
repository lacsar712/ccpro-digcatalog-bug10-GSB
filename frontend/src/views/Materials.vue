<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">材质分类</h2>
        <p class="page-sub">维护文物材质字典</p>
      </div>
      <button class="btn" @click="openCreate">新增材质</button>
    </div>

    <div class="card">
      <table class="table">
        <thead>
          <tr>
            <th>名称</th>
            <th>描述</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td>{{ item.name }}</td>
            <td>{{ item.description || '-' }}</td>
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
        <h3>{{ form.id ? '编辑材质' : '新增材质' }}</h3>
        <div class="form-grid">
          <label class="full">
            名称
            <input v-model="form.name" />
          </label>
          <label class="full">
            描述
            <textarea v-model="form.description" />
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
const form = reactive({ id: null, name: '', description: '' })

async function load() {
  error.value = ''
  try {
    const { data } = await api.get('/materials')
    list.value = data
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function openCreate() {
  Object.assign(form, { id: null, name: '', description: '' })
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
    const payload = { name: form.name, description: form.description }
    if (form.id) {
      await api.put(`/materials/${form.id}`, payload)
    } else {
      await api.post('/materials', payload)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || '保存失败'
  }
}

async function remove(item) {
  if (!confirm(`确认删除材质「${item.name}」？`)) return
  try {
    await api.delete(`/materials/${item.id}`)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(load)
</script>
