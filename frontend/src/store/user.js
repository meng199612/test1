import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, getCurrentUser, getUserMenus } from '@/api/auth'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref({})
  const menus = ref([])
  const permissions = ref([])

  const isLoggedIn = computed(() => !!token.value)

  async function doLogin(form) {
    const data = await login(form)
    token.value = data.token
    userInfo.value = data.user_info
    localStorage.setItem('token', data.token)
    return data
  }

  async function fetchUserInfo() {
    const data = await getCurrentUser()
    userInfo.value = data
    permissions.value = data.permissions || []
    return data
  }

  async function fetchMenus() {
    const data = await getUserMenus()
    menus.value = data || []
    return data
  }

  function logout() {
    token.value = ''
    userInfo.value = {}
    menus.value = []
    permissions.value = []
    localStorage.removeItem('token')
  }

  function hasPermission(code) {
    return permissions.value.includes(code)
  }

  return {
    token,
    userInfo,
    menus,
    permissions,
    isLoggedIn,
    doLogin,
    fetchUserInfo,
    fetchMenus,
    logout,
    hasPermission
  }
})
