<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

// Sidebar state
const sidebarOpen = ref(true)
const sidebarCollapsed = ref(false)

// Navigation items
const navigationItems = [
  {
    name: 'Dashboard',
    path: '/',
    icon: '🏠',
    description: 'Overview and key metrics'
  },
  {
    name: 'Sales Details',
    path: '/details',
    icon: '📊',
    description: 'Raw data with filters and search'
  },
  {
    name: 'Reports',
    path: '/reports',
    icon: '📈',
    description: 'Pivot tables and analytics'
  }
]

// Computed properties
const currentPath = computed(() => route.path)
const sidebarWidth = computed(() => {
  if (!sidebarOpen.value) return 'w-0'
  return sidebarCollapsed.value ? 'w-16' : 'w-64'
})

// Methods
const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const toggleCollapse = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const navigateTo = (path: string) => {
  router.push(path)
}

const isActive = (path: string) => {
  return currentPath.value === path
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex">
    <!-- Sidebar -->
    <aside 
      :class="[
        'bg-white shadow-lg transition-all duration-300 ease-in-out flex-shrink-0',
        sidebarWidth,
        sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
      ]"
      class="fixed lg:relative h-full z-30"
    >
      <!-- Sidebar Header -->
      <div class="p-4 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <div v-if="!sidebarCollapsed" class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-sm">ST</span>
            </div>
            <div>
              <h1 class="font-bold text-gray-800">Sales Track</h1>
              <p class="text-xs text-gray-500">Data Management</p>
            </div>
          </div>
          <button
            @click="toggleCollapse"
            class="p-1 rounded-md hover:bg-gray-100 transition-colors"
            :class="sidebarCollapsed ? 'mx-auto' : ''"
          >
            <span class="text-gray-500">{{ sidebarCollapsed ? '→' : '←' }}</span>
          </button>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="p-4 space-y-2">
        <button
          v-for="item in navigationItems"
          :key="item.path"
          @click="navigateTo(item.path)"
          :class="[
            'w-full flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors text-left',
            isActive(item.path)
              ? 'bg-primary-100 text-primary-700 border border-primary-200'
              : 'text-gray-600 hover:bg-gray-100 hover:text-gray-800'
          ]"
        >
          <span class="text-lg flex-shrink-0">{{ item.icon }}</span>
          <div v-if="!sidebarCollapsed" class="flex-1 min-w-0">
            <div class="font-medium">{{ item.name }}</div>
            <div class="text-xs text-gray-500 truncate">{{ item.description }}</div>
          </div>
        </button>
      </nav>

      <!-- Sidebar Footer -->
      <div class="absolute bottom-4 left-4 right-4" v-if="!sidebarCollapsed">
        <div class="text-xs text-gray-400 text-center">
          <p>Version 1.0.0</p>
          <p>Built with Wails</p>
        </div>
      </div>
    </aside>

    <!-- Mobile Sidebar Overlay -->
    <div
      v-if="sidebarOpen"
      @click="toggleSidebar"
      class="fixed inset-0 bg-black bg-opacity-50 z-20 lg:hidden"
    ></div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top Header -->
      <header class="bg-white shadow-sm border-b border-gray-200 px-4 py-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <!-- Mobile Menu Button -->
            <button
              @click="toggleSidebar"
              class="lg:hidden p-2 rounded-md hover:bg-gray-100 transition-colors"
            >
              <span class="text-gray-600">☰</span>
            </button>
            
            <!-- Page Title -->
            <div>
              <h2 class="text-xl font-semibold text-gray-800">
                {{ navigationItems.find(item => item.path === currentPath)?.name || 'Sales Track' }}
              </h2>
              <p class="text-sm text-gray-500">
                {{ navigationItems.find(item => item.path === currentPath)?.description || 'Modern sales data management' }}
              </p>
            </div>
          </div>

          <!-- Header Actions -->
          <div class="flex items-center space-x-3">
            <!-- Database Status Indicator -->
            <div class="flex items-center space-x-2 text-sm">
              <div class="w-2 h-2 bg-success-500 rounded-full"></div>
              <span class="text-gray-600 hidden sm:inline">Connected</span>
            </div>
            
            <!-- Settings Button -->
            <button class="p-2 rounded-md hover:bg-gray-100 transition-colors">
              <span class="text-gray-600">⚙️</span>
            </button>
          </div>
        </div>
      </header>

      <!-- Main Content -->
      <main class="flex-1 p-6 overflow-auto">
        <slot />
      </main>
    </div>
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */
</style>
