<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { GetImportStatistics, GetDatabaseHealth } from '../../wailsjs/go/main/App'
import type { main } from '../../wailsjs/go/models'

// Reactive data
const loading = ref(true)
const error = ref<string | null>(null)
const statistics = ref<main.ImportStatistics | null>(null)
const dbHealth = ref<main.DatabaseHealth | null>(null)

// Mock data for additional metrics (will be replaced with real API calls)
const mockMetrics = ref({
  mtdSales: 15420.50,
  ytdSales: 187650.25,
  bestSellingProduct: "Samsung 55\" Smart TV",
  monthlyItemsSold: 142,
  yearlyItemsSold: 1847
})

// Computed properties
const formattedMTDSales = computed(() => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(mockMetrics.value.mtdSales)
})

const formattedYTDSales = computed(() => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(mockMetrics.value.ytdSales)
})

const formattedTotalSales = computed(() => {
  if (!statistics.value) return '$0.00'
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(statistics.value.total_sales)
})

const formattedAveragePrice = computed(() => {
  if (!statistics.value) return '$0.00'
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(statistics.value.average_price)
})

// Methods
const loadDashboardData = async () => {
  try {
    loading.value = true
    error.value = null

    // Load statistics and database health
    const [statsResult, healthResult] = await Promise.all([
      GetImportStatistics(),
      GetDatabaseHealth()
    ])

    statistics.value = statsResult
    dbHealth.value = healthResult

  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load dashboard data'
    console.error('Dashboard data loading error:', err)
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  loadDashboardData()
}

// Lifecycle
onMounted(() => {
  loadDashboardData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-gray-600 mt-1">Overview of your sales performance and key metrics</p>
      </div>
      <button
        @click="refreshData"
        :disabled="loading"
        class="btn btn-primary flex items-center space-x-2"
      >
        <span :class="loading ? 'animate-spin' : ''">🔄</span>
        <span>Refresh</span>
      </button>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-error-50 border border-error-200 rounded-lg p-4">
      <div class="flex items-center space-x-2">
        <span class="text-error-600">⚠️</span>
        <div>
          <h3 class="font-medium text-error-800">Error Loading Dashboard</h3>
          <p class="text-error-600 text-sm mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="i in 6" :key="i" class="card animate-pulse">
        <div class="h-4 bg-gray-200 rounded w-3/4 mb-3"></div>
        <div class="h-8 bg-gray-200 rounded w-1/2 mb-2"></div>
        <div class="h-3 bg-gray-200 rounded w-full"></div>
      </div>
    </div>

    <!-- Dashboard Cards -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- MTD Sales -->
      <div class="card bg-gradient-to-br from-primary-50 to-primary-100 border-primary-200">
        <div class="flex items-center justify-between mb-4">
          <div class="p-2 bg-primary-600 rounded-lg">
            <span class="text-white text-xl">📈</span>
          </div>
          <span class="text-xs text-primary-600 font-medium">MTD</span>
        </div>
        <h3 class="text-sm font-medium text-primary-700 mb-1">Month-to-Date Sales</h3>
        <p class="text-2xl font-bold text-primary-900">{{ formattedMTDSales }}</p>
        <p class="text-xs text-primary-600 mt-2">Current month performance</p>
      </div>

      <!-- YTD Sales -->
      <div class="card bg-gradient-to-br from-success-50 to-success-100 border-success-200">
        <div class="flex items-center justify-between mb-4">
          <div class="p-2 bg-success-600 rounded-lg">
            <span class="text-white text-xl">💰</span>
          </div>
          <span class="text-xs text-success-600 font-medium">YTD</span>
        </div>
        <h3 class="text-sm font-medium text-success-700 mb-1">Year-to-Date Sales</h3>
        <p class="text-2xl font-bold text-success-900">{{ formattedYTDSales }}</p>
        <p class="text-xs text-success-600 mt-2">Annual performance</p>
      </div>

      <!-- Best Selling Product -->
      <div class="card bg-gradient-to-br from-warning-50 to-warning-100 border-warning-200">
        <div class="flex items-center justify-between mb-4">
          <div class="p-2 bg-warning-600 rounded-lg">
            <span class="text-white text-xl">🏆</span>
          </div>
          <span class="text-xs text-warning-600 font-medium">TOP</span>
        </div>
        <h3 class="text-sm font-medium text-warning-700 mb-1">Best Selling Product</h3>
        <p class="text-lg font-bold text-warning-900 truncate">{{ mockMetrics.bestSellingProduct }}</p>
        <p class="text-xs text-warning-600 mt-2">Most popular item</p>
      </div>

      <!-- Monthly Items Sold -->
      <div class="card bg-gradient-to-br from-purple-50 to-purple-100 border-purple-200">
        <div class="flex items-center justify-between mb-4">
          <div class="p-2 bg-purple-600 rounded-lg">
            <span class="text-white text-xl">📦</span>
          </div>
          <span class="text-xs text-purple-600 font-medium">MONTH</span>
        </div>
        <h3 class="text-sm font-medium text-purple-700 mb-1">Items Sold This Month</h3>
        <p class="text-2xl font-bold text-purple-900">{{ mockMetrics.monthlyItemsSold.toLocaleString() }}</p>
        <p class="text-xs text-purple-600 mt-2">Units moved in {{ new Date().toLocaleDateString('en-US', { month: 'long' }) }}</p>
      </div>

      <!-- Yearly Items Sold -->
      <div class="card bg-gradient-to-br from-indigo-50 to-indigo-100 border-indigo-200">
        <div class="flex items-center justify-between mb-4">
          <div class="p-2 bg-indigo-600 rounded-lg">
            <span class="text-white text-xl">📊</span>
          </div>
          <span class="text-xs text-indigo-600 font-medium">YEAR</span>
        </div>
        <h3 class="text-sm font-medium text-indigo-700 mb-1">Items Sold This Year</h3>
        <p class="text-2xl font-bold text-indigo-900">{{ mockMetrics.yearlyItemsSold.toLocaleString() }}</p>
        <p class="text-xs text-indigo-600 mt-2">Total units in {{ new Date().getFullYear() }}</p>
      </div>

      <!-- Database Statistics -->
      <div class="card bg-gradient-to-br from-gray-50 to-gray-100 border-gray-200">
        <div class="flex items-center justify-between mb-4">
          <div class="p-2 bg-gray-600 rounded-lg">
            <span class="text-white text-xl">🗄️</span>
          </div>
          <div class="flex items-center space-x-1">
            <div :class="[
              'w-2 h-2 rounded-full',
              dbHealth?.connected ? 'bg-success-500' : 'bg-error-500'
            ]"></div>
            <span class="text-xs text-gray-600 font-medium">
              {{ dbHealth?.connected ? 'ONLINE' : 'OFFLINE' }}
            </span>
          </div>
        </div>
        <h3 class="text-sm font-medium text-gray-700 mb-1">Database Statistics</h3>
        <div class="space-y-1">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Total Records:</span>
            <span class="font-medium">{{ statistics?.total_records?.toLocaleString() || '0' }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Recent Records:</span>
            <span class="font-medium">{{ statistics?.recent_records?.toLocaleString() || '0' }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Average Price:</span>
            <span class="font-medium">{{ formattedAveragePrice }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <router-link
          to="/details"
          class="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors group"
        >
          <div class="p-2 bg-primary-100 rounded-lg group-hover:bg-primary-200 transition-colors">
            <span class="text-primary-600">📊</span>
          </div>
          <div>
            <h3 class="font-medium text-gray-900">View Sales Details</h3>
            <p class="text-sm text-gray-500">Browse and manage records</p>
          </div>
        </router-link>

        <router-link
          to="/reports"
          class="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:border-success-300 hover:bg-success-50 transition-colors group"
        >
          <div class="p-2 bg-success-100 rounded-lg group-hover:bg-success-200 transition-colors">
            <span class="text-success-600">📈</span>
          </div>
          <div>
            <h3 class="font-medium text-gray-900">Generate Reports</h3>
            <p class="text-sm text-gray-500">Analyze sales data</p>
          </div>
        </router-link>

        <button
          @click="refreshData"
          class="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:border-warning-300 hover:bg-warning-50 transition-colors group"
        >
          <div class="p-2 bg-warning-100 rounded-lg group-hover:bg-warning-200 transition-colors">
            <span class="text-warning-600">🔄</span>
          </div>
          <div>
            <h3 class="font-medium text-gray-900">Refresh Data</h3>
            <p class="text-sm text-gray-500">Update dashboard metrics</p>
          </div>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */
</style>
