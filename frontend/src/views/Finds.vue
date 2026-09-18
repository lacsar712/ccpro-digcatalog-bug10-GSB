<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">出土文物</h2>
        <p class="page-sub">登记器物类型、材质、完整度与存放位置</p>
      </div>
      <button class="btn" @click="openCreate">新增文物</button>
    </div>

    <div class="card">
      <div class="filters">
        <label>
          探方筛选
          <select v-model="filterUnitId" @change="load">
            <option value="">全部探方</option>
            <option v-for="u in units" :key="u.id" :value="String(u.id)">
              {{ u.site?.name || '' }} / {{ u.code }}
            </option>
          </select>
        </label>
        <label>
          器物类型
          <select v-model="filterType" @change="load">
            <option value="">全部类型</option>
            <option v-for="t in artifactTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </label>
      </div>

      <table class="table">
        <thead>
          <tr>
            <th>登记号</th>
            <th>探方</th>
            <th>器物类型</th>
            <th>材质</th>
            <th>完整度</th>
            <th>出土日期</th>
            <th>存放位置</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td>{{ item.registerNo }}</td>
            <td>{{ item.unit?.code || '-' }}</td>
            <td><span class="tag">{{ item.artifactType }}</span></td>
            <td>{{ item.materialName || '-' }}</td>
            <td>{{ item.completeness || '-' }}</td>
            <td>{{ formatDate(item.findDate) }}</td>
            <td>{{ item.storageLoc || '-' }}</td>
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
        <h3>{{ form.id ? '编辑文物' : '新增文物' }}</h3>
        <div class="form-grid">
          <label>
            所属探方
            <select v-model.number="form.unitId">
              <option :value="0" disabled>请选择</option>
              <option v-for="u in units" :key="u.id" :value="u.id">
                {{ u.site?.name || '' }} / {{ u.code }}
              </option>
            </select>
          </label>
          <label>
            登记号
            <input v-model="form.registerNo" />
          </label>
          <label>
            器物类型
            <select v-model="form.artifactType">
              <option v-for="t in artifactTypes" :key="t" :value="t">{{ t }}</option>
            </select>
          </label>
          <label>
            材质
            <select v-model="form.materialId">
              <option :value="null">未指定</option>
              <option v-for="m in materials" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </label>
          <label>
            完整度
            <select v-model="form.completeness">
              <option>完整</option>
              <option>残缺</option>
              <option>碎片</option>
            </select>
          </label>
          <label>
            出土日期
            <input v-model="form.findDate" type="date" />
          </label>
          <label class="full">
            存放位置
            <input v-model="form.storageLoc" />
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

const artifactTypes = ['陶片', '青铜器', '骨器', '玉器', '石器', '铁器', '其他']
const list = ref([])
const units = ref([])
const materials = ref([])
const filterUnitId = ref('')
const filterType = ref('')
const error = ref('')
const formError = ref('')
const showModal = ref(false)

const form = reactive({
  id: null,
  unitId: 0,
  materialId: null,
  registerNo: '',
  artifactType: '陶片',
  completeness: '完整',
  findDate: '',
  description: '',
  storageLoc: ''
})

function formatDate(v) {
  if (!v) return '-'
  return String(v).slice(0, 10)
}

async function loadMeta() {
  const [u, m] = await Promise.all([api.get('/units'), api.get('/materials')])
  units.value = u.data
  materials.value = m.data
}

async function load() {
  error.value = ''
  try {
    const params = {}
    if (filterUnitId.value) params.unitId = filterUnitId.value
    if (filterType.value) params.artifactType = filterType.value
    const { data } = await api.get('/finds', { params })
    list.value = data
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function openCreate() {
  Object.assign(form, {
    id: null,
    unitId: units.value[0]?.id || 0,
    materialId: materials.value[0]?.id ?? null,
    registerNo: '',
    artifactType: '陶片',
    completeness: '完整',
    findDate: '',
    description: '',
    storageLoc: ''
  })
  formError.value = ''
  showModal.value = true
}

function openEdit(item) {
  Object.assign(form, {
    id: item.id,
    unitId: item.unitId,
    materialId: item.materialId,
    registerNo: item.registerNo,
    artifactType: item.artifactType,
    completeness: item.completeness || '完整',
    findDate: formatDate(item.findDate) === '-' ? '' : formatDate(item.findDate),
    description: item.description || '',
    storageLoc: item.storageLoc || ''
  })
  formError.value = ''
  showModal.value = true
}

async function save() {
  formError.value = ''
  try {
    const payload = {
      unitId: form.unitId,
      materialId: form.materialId || null,
      registerNo: form.registerNo,
      artifactType: form.artifactType,
      completeness: form.completeness,
      findDate: form.findDate || null,
      description: form.description,
      storageLoc: form.storageLoc
    }
    if (form.id) {
      await api.put(`/finds/${form.id}`, payload)
    } else {
      await api.post('/finds', payload)
    }
    showModal.value = false
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || '保存失败'
  }
}

async function remove(item) {
  if (!confirm(`确认删除文物「${item.registerNo}」？`)) return
  try {
    await api.delete(`/finds/${item.id}`)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(async () => {
  await loadMeta()
  await load()
})
</script>

<style scoped>
.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.filters label {
  min-width: 200px;
}
</style>
