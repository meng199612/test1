<template>
  <div class="app-container">
    <el-card>
      <el-form :inline="true" :model="queryForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList"><el-icon><Search /></el-icon> 查询</el-button>
          <el-button @click="resetQuery"><el-icon><Refresh /></el-icon> 重置</el-button>
          <el-button type="success" @click="handleAdd"><el-icon><Plus /></el-icon> 新增</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="handleResetPwd(row)">重置密码</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="用户名" v-if="!form.id">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" v-if="!form.id" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple placeholder="请选择角色" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pwdDialogVisible" title="重置密码" width="400px">
      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="80px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="pwdForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pwdDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPwd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUsers, createUser, updateUser, deleteUser, resetPassword } from '@/api/user'
import { getRoles } from '@/api/role'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const roles = ref([])

const queryForm = reactive({
  page: 1,
  page_size: 10,
  username: '',
  status: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const formRef = ref()
const form = reactive({
  id: null,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  status: 1,
  role_ids: []
})

const rules = {
  password: [{ required: true, message: '请输入密码', min: 6, trigger: 'blur' }]
}

const pwdDialogVisible = ref(false)
const pwdFormRef = ref()
const pwdForm = reactive({ id: null, password: '' })
const pwdRules = {
  password: [{ required: true, message: '请输入新密码', min: 6, trigger: 'blur' }]
}

async function getList() {
  loading.value = true
  try {
    const res = await getUsers(queryForm)
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  queryForm.page = 1
  queryForm.username = ''
  queryForm.status = ''
  getList()
}

function handleAdd() {
  dialogTitle.value = '新增用户'
  form.id = null
  form.username = ''
  form.password = ''
  form.nickname = ''
  form.email = ''
  form.phone = ''
  form.status = 1
  form.role_ids = []
  dialogVisible.value = true
}

function handleEdit(row) {
  dialogTitle.value = '编辑用户'
  form.id = row.id
  form.username = row.username
  form.password = ''
  form.nickname = row.nickname
  form.email = row.email
  form.phone = row.phone
  form.status = row.status
  form.role_ids = row.roles?.map(r => r.id) || []
  dialogVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid && !form.id) return

  if (form.id) {
    await updateUser(form.id, {
      nickname: form.nickname,
      email: form.email,
      phone: form.phone,
      status: form.status,
      role_ids: form.role_ids
    })
    ElMessage.success('更新成功')
  } else {
    await createUser(form)
    ElMessage.success('新增成功')
  }

  dialogVisible.value = false
  getList()
}

function handleDelete(row) {
  ElMessageBox.confirm(`确定要删除用户"${row.username}"吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    await deleteUser(row.id)
    ElMessage.success('删除成功')
    getList()
  }).catch(() => {})
}

function handleResetPwd(row) {
  pwdForm.id = row.id
  pwdForm.password = ''
  pwdDialogVisible.value = true
}

async function submitPwd() {
  const valid = await pwdFormRef.value?.validate().catch(() => false)
  if (!valid) return

  await resetPassword(pwdForm.id, { password: pwdForm.password })
  ElMessage.success('重置密码成功')
  pwdDialogVisible.value = false
}

onMounted(async () => {
  await Promise.all([getList(), loadRoles()])
})

async function loadRoles() {
  try {
    roles.value = await getRoles()
  } catch (e) {}
}
</script>
