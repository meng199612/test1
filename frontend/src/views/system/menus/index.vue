<template>
  <div class="app-container">
    <el-card>
      <el-table :data="list" v-loading="loading" border stripe row-key="id" default-expand-all>
        <el-table-column prop="name" label="菜单名称" width="200" />
        <el-table-column prop="path" label="路由路径" />
        <el-table-column prop="component" label="组件路径" />
        <el-table-column prop="icon" label="图标" width="100" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="可见" width="80">
          <template #default="{ row }">
            <el-tag :type="row.visible === 1 ? 'success' : 'info'">
              {{ row.visible === 1 ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parent_id"
            :data="menuOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择上级菜单"
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="路由路径" prop="path">
          <el-input v-model="form.path" placeholder="请输入路由路径" />
        </el-form-item>
        <el-form-item label="组件路径">
          <el-input v-model="form.component" placeholder="如：system/users/index" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="请输入图标名称" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="可见">
          <el-switch v-model="form.visible" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAllMenus, getRole } from '@/api/role'
import request from '@/utils/request'

const loading = ref(false)
const list = ref([])
const menuOptions = ref([])

const dialogVisible = ref(false)
const dialogTitle = ref('编辑菜单')
const formRef = ref()
const form = reactive({
  id: null,
  parent_id: 0,
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  visible: 1
})

const rules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  path: [{ required: true, message: '请输入路由路径', trigger: 'blur' }]
}

async function getList() {
  loading.value = true
  try {
    list.value = await getAllMenus()
    menuOptions.value = [{ id: 0, name: '顶级菜单', children: list.value }]
  } finally {
    loading.value = false
  }
}

function handleEdit(row) {
  dialogTitle.value = '编辑菜单'
  Object.assign(form, {
    id: row.id,
    parent_id: row.parent_id,
    name: row.name,
    path: row.path,
    component: row.component,
    icon: row.icon,
    sort: row.sort,
    visible: row.visible
  })
  dialogVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  if (form.id) {
    await request.put(`/menus/${form.id}`, form)
    ElMessage.success('更新成功')
  } else {
    await request.post('/menus', form)
    ElMessage.success('新增成功')
  }

  dialogVisible.value = false
  getList()
}

function handleDelete(row) {
  ElMessageBox.confirm(`确定要删除菜单"${row.name}"吗？`, '提示', { type: 'warning' })
    .then(async () => {
      await request.delete(`/menus/${row.id}`)
      ElMessage.success('删除成功')
      getList()
    }).catch(() => {})
}

onMounted(() => {
  getList()
})
</script>
