<template>
  <div class="app-container">
    <el-card>
      <div style="margin-bottom: 20px">
        <el-button type="success" @click="handleAdd"><el-icon><Plus /></el-icon> 新增</el-button>
      </div>

      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="角色编码" width="150" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="success" @click="handlePermissions(row)">分配权限</el-button>
            <el-button link type="warning" @click="handleMenus(row)">分配菜单</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色编码" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permDialogVisible" title="分配权限" width="600px">
      <el-tree
        ref="permTreeRef"
        :data="permissionTree"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedPermIds"
        :props="{ label: 'name', children: 'children' }"
      />
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPermissions">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="menuDialogVisible" title="分配菜单" width="600px">
      <el-tree
        ref="menuTreeRef"
        :data="menuTree"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedMenuIds"
        :props="{ label: 'name', children: 'children' }"
      />
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitMenus">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRoles, createRole, updateRole, deleteRole, assignPermissions, assignMenus, getPermissionsTree, getAllMenus } from '@/api/role'

const loading = ref(false)
const list = ref([])
const permissionTree = ref([])
const menuTree = ref([])

const dialogVisible = ref(false)
const dialogTitle = ref('新增角色')
const formRef = ref()
const form = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }]
}

const permDialogVisible = ref(false)
const permTreeRef = ref()
const currentRoleId = ref(null)
const checkedPermIds = ref([])

const menuDialogVisible = ref(false)
const menuTreeRef = ref()
const checkedMenuIds = ref([])

async function getList() {
  loading.value = true
  try {
    list.value = await getRoles()
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  dialogTitle.value = '新增角色'
  form.id = null
  form.name = ''
  form.code = ''
  form.description = ''
  form.status = 1
  dialogVisible.value = true
}

function handleEdit(row) {
  dialogTitle.value = '编辑角色'
  form.id = row.id
  form.name = row.name
  form.code = row.code
  form.description = row.description
  form.status = row.status
  dialogVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  if (form.id) {
    await updateRole(form.id, {
      name: form.name,
      description: form.description,
      status: form.status
    })
    ElMessage.success('更新成功')
  } else {
    await createRole(form)
    ElMessage.success('新增成功')
  }

  dialogVisible.value = false
  getList()
}

function handleDelete(row) {
  ElMessageBox.confirm(`确定要删除角色"${row.name}"吗？`, '提示', { type: 'warning' })
    .then(async () => {
      await deleteRole(row.id)
      ElMessage.success('删除成功')
      getList()
    }).catch(() => {})
}

async function handlePermissions(row) {
  currentRoleId.value = row.id
  const role = await getRoleDetail(row.id)
  checkedPermIds.value = role.permissions?.map(p => p.id) || []
  permDialogVisible.value = true
}

async function submitPermissions() {
  const keys = permTreeRef.value?.getCheckedKeys() || []
  await assignPermissions(currentRoleId.value, keys)
  ElMessage.success('分配成功')
  permDialogVisible.value = false
}

async function handleMenus(row) {
  currentRoleId.value = row.id
  const role = await getRoleDetail(row.id)
  checkedMenuIds.value = role.menus?.map(m => m.id) || []
  menuDialogVisible.value = true
}

async function submitMenus() {
  const keys = menuTreeRef.value?.getCheckedKeys() || []
  await assignMenus(currentRoleId.value, keys)
  ElMessage.success('分配成功')
  menuDialogVisible.value = false
}

async function getRoleDetail(id) {
  try {
    return await getRoles().then(() => getRoleById(id))
  } catch (e) {
    return { permissions: [], menus: [] }
  }
}

async function getRoleById(id) {
  const all = list.value
  return all.find(r => r.id === id) || { permissions: [], menus: [] }
}

onMounted(async () => {
  await Promise.all([
    getList(),
    loadPermissionTree(),
    loadMenuTree()
  ])
})

async function loadPermissionTree() {
  try {
    permissionTree.value = await getPermissionsTree()
  } catch (e) {
    permissionTree.value = []
  }
}

async function loadMenuTree() {
  try {
    menuTree.value = await getAllMenus()
  } catch (e) {
    menuTree.value = []
  }
}
</script>
