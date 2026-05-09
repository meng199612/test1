<template>
  <div class="app-container">
    <el-card>
      <el-form :inline="true" :model="queryForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker v-model="queryForm.start_time" type="date" placeholder="选择开始时间" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker v-model="queryForm.end_time" type="date" placeholder="选择结束时间" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList"><el-icon><Search /></el-icon> 查询</el-button>
          <el-button @click="resetQuery"><el-icon><Refresh /></el-icon> 重置</el-button>
          <el-button type="success" @click="handleExport"><el-icon><Download /></el-icon> 导出</el-button>
          <el-button type="danger" @click="handleClear"><el-icon><Delete /></el-icon> 清理</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="操作用户" width="120" />
        <el-table-column prop="module" label="操作模块" width="120" />
        <el-table-column prop="operation" label="操作类型" width="100" />
        <el-table-column prop="method" label="请求方法" width="80" />
        <el-table-column prop="path" label="请求路径" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.page_size"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination-container"
        @size-change="getList"
        @current-change="getList"
      />
    </el-card>

    <el-dialog v-model="detailVisible" title="日志详情" width="800px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item label="操作用户">{{ currentLog.username }}</el-descriptions-item>
        <el-descriptions-item label="操作模块">{{ currentLog.module }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ currentLog.operation }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ currentLog.method }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item label="请求路径" :span="2">{{ currentLog.path }}</el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <pre style="white-space: pre-wrap; word-break: break-all">{{ currentLog.params }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应结果" :span="2">
          <pre style="white-space: pre-wrap; word-break: break-all; max-height: 200px; overflow-y: auto">{{ currentLog.result }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentLog.status === 1 ? 'success' : 'danger'">
            {{ currentLog.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ currentLog.created_at }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as XLSX from 'xlsx'
import { getOperationLogs, getOperationLog, deleteOperationLog, clearOperationLogs } from '@/api/log'

const loading = ref(false)
const list = ref([])
const total = ref(0)

const queryForm = reactive({
  page: 1,
  page_size: 10,
  username: '',
  start_time: '',
  end_time: ''
})

const detailVisible = ref(false)
const currentLog = ref({})

async function getList() {
  loading.value = true
  try {
    const res = await getOperationLogs(queryForm)
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  queryForm.page = 1
  queryForm.username = ''
  queryForm.start_time = ''
  queryForm.end_time = ''
  getList()
}

async function handleDetail(row) {
  currentLog.value = await getOperationLog(row.id)
  detailVisible.value = true
}

function handleExport() {
  if (list.value.length === 0) {
    ElMessage.warning('没有数据可导出')
    return
  }

  const exportData = list.value.map(item => ({
    ID: item.id,
    用户名: item.username,
    模块: item.module,
    操作: item.operation,
    方法: item.method,
    路径: item.path,
    IP: item.ip,
    状态: item.status === 1 ? '成功' : '失败',
    时间: item.created_at
  }))

  const ws = XLSX.utils.json_to_sheet(exportData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '操作日志')
  XLSX.writeFile(wb, `操作日志_${new Date().toISOString().split('T')[0]}.xlsx`)
}

function handleClear() {
  ElMessageBox.confirm('确定要清理30天前的日志吗？此操作不可恢复！', '提示', {
    type: 'warning'
  }).then(async () => {
    await clearOperationLogs(30)
    ElMessage.success('清理成功')
    getList()
  }).catch(() => {})
}

onMounted(() => {
  getList()
})
</script>
