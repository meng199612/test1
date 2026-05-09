<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon user">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.users }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon role">
            <el-icon><UserFilled /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.roles }}</div>
            <div class="stat-label">角色数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon menu">
            <el-icon><Menu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.menus }}</div>
            <div class="stat-label">菜单数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon log">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.logs }}</div>
            <div class="stat-label">今日操作</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>欢迎使用</template>
          <div class="welcome">
            <h2>欢迎使用后台管理系统</h2>
            <p>当前用户：{{ userStore.userInfo.nickname || userStore.userInfo.username }}</p>
            <p>登录时间：{{ loginTime }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>快捷入口</template>
          <div class="shortcuts">
            <el-button type="primary" @click="$router.push('/system/users')">
              <el-icon><User /></el-icon> 用户管理
            </el-button>
            <el-button type="success" @click="$router.push('/system/roles')">
              <el-icon><UserFilled /></el-icon> 角色管理
            </el-button>
            <el-button type="warning" @click="$router.push('/system/logs')">
              <el-icon><Document /></el-icon> 操作日志
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { getUsers } from '@/api/user'
import { getRoles, getAllMenus } from '@/api/role'
import { getOperationLogs } from '@/api/log'

const userStore = useUserStore()

const stats = ref({ users: 0, roles: 0, menus: 0, logs: 0 })
const loginTime = ref(new Date().toLocaleString())

onMounted(async () => {
  try {
    const [users, roles, menus, logs] = await Promise.all([
      getUsers({ page: 1, page_size: 1 }),
      getRoles(),
      getAllMenus(),
      getOperationLogs({ page: 1, page_size: 1, start_time: new Date().toISOString().split('T')[0] })
    ])
    stats.value.users = users.total || 0
    stats.value.roles = roles?.length || 0
    stats.value.menus = menus?.length || 0
    stats.value.logs = logs.total || 0
  } catch (e) {}
})
</script>

<style scoped lang="scss">
.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
    color: #fff;

    &.user { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
    &.role { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
    &.menu { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
    &.log { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
  }

  .stat-info {
    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: #333;
    }
    .stat-label {
      font-size: 14px;
      color: #999;
    }
  }
}

.welcome {
  h2 { margin-bottom: 10px; color: #333; }
  p { margin: 5px 0; color: #666; }
}

.shortcuts {
  display: flex;
  flex-direction: column;
  gap: 10px;

  .el-button {
    width: 140px;
    justify-content: flex-start;
  }
}
</style>
