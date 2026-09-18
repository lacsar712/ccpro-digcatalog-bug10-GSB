<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">概览</h2>
        <p class="page-sub">工地、探方与出土文物登记统计</p>
      </div>
      <button class="btn secondary" @click="load">刷新</button>
    </div>

    <div v-if="error" class="error card">{{ error }}</div>

    <div class="stats">
      <div class="stat card">
        <div class="label">发掘工地</div>
        <div class="value">{{ data.siteCount ?? '-' }}</div>
      </div>
      <div class="stat card">
        <div class="label">探方/发掘单位</div>
        <div class="value">{{ data.unitCount ?? '-' }}</div>
      </div>
      <div class="stat card">
        <div class="label">出土文物总数</div>
        <div class="value">{{ data.findCount ?? '-' }}</div>
      </div>
    </div>

    <div class="card">
      <h3>按器物类型统计</h3>
      <table class="table" v-if="data.byType?.length">
        <thead>
          <tr>
            <th>器物类型</th>
            <th>数量</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in data.byType" :key="row.artifactType">
            <td><span class="tag">{{ row.artifactType }}</span></td>
            <td>{{ row.count }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="page-sub">暂无统计数据</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import api from '../api/http'

const data = reactive({
  siteCount: 0,
  unitCount: 0,
  findCount: 0,
  byType: []
})
const error = ref('')

async function load() {
  error.value = ''
  try {
    const { data: res } = await api.get('/overview')
    Object.assign(data, res)
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

onMounted(load)
</script>

<style scoped>
.stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.stat .label {
  color: var(--muted);
  font-size: 0.9rem;
}

.stat .value {
  font-size: 2rem;
  font-weight: 700;
  margin-top: 0.35rem;
  color: var(--accent);
}

h3 {
  margin: 0 0 0.75rem;
}

@media (max-width: 800px) {
  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
