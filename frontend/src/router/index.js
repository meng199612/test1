import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const constantRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '404' }
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '403' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes
})

const whiteList = ['/login', '/404', '/403']

let dynamicRoutesAdded = false

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()

  if (whiteList.includes(to.path)) {
    next()
    return
  }

  if (userStore.token) {
    if (to.path === '/login') {
      next('/')
    } else {
      if (!dynamicRoutesAdded && userStore.menus.length === 0) {
        try {
          await userStore.fetchUserInfo()
          const menus = await userStore.fetchMenus()
          addDynamicRoutes(menus)
          dynamicRoutesAdded = true
          next({ ...to, replace: true })
        } catch (error) {
          userStore.logout()
          next(`/login?redirect=${to.path}`)
        }
      } else {
        if (to.matched.length === 0) {
          next('/404')
        } else {
          next()
        }
      }
    }
  } else {
    next(`/login?redirect=${to.path}`)
  }
})

function addDynamicRoutes(menus) {
  const layoutRoute = {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
      }
    ]
  }

  const routes = generateRoutes(menus)
  layoutRoute.children.push(...routes)

  layoutRoute.children.push({
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  })

  router.addRoute(layoutRoute)
}

function generateRoutes(menus, parentPath = '') {
  const routes = []
  menus.forEach(menu => {
    const fullPath = menu.path
    const route = {
      path: fullPath,
      name: menu.name,
      meta: { title: menu.name, icon: menu.icon || 'Menu' }
    }

    if (menu.component) {
      try {
        route.component = () => import(`@/views/${menu.component}.vue`)
      } catch (e) {
        route.component = () => import('@/views/error/404.vue')
      }
    }

    if (menu.children && menu.children.length > 0) {
      route.children = generateRoutes(menu.children, fullPath)
    }

    routes.push(route)
  })
  return routes
}

export function resetRouter() {
  const newRouter = createRouter({
    history: createWebHistory(),
    routes: constantRoutes
  })
  router.matcher = newRouter.matcher
  dynamicRoutesAdded = false
}

export default router
