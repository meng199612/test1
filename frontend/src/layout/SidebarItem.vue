<template>
  <template v-if="item.children && item.children.length > 0">
    <el-sub-menu :index="item.path">
      <template #title>
        <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
        <span>{{ item.name }}</span>
      </template>
      <sidebar-item
        v-for="child in item.children"
        :key="child.id"
        :item="child"
        :base-path="item.path"
      />
    </el-sub-menu>
  </template>
  <template v-else>
    <el-menu-item :index="item.path">
      <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
      <template #title>{{ item.name }}</template>
    </el-menu-item>
  </template>
</template>

<script setup>
defineProps({
  item: {
    type: Object,
    required: true
  },
  basePath: {
    type: String,
    default: ''
  }
})
</script>
