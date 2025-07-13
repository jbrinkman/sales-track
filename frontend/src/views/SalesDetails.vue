<script lang="ts" setup>
import { ref, onMounted, computed, watch } from 'vue'
import { GetRecentImports, ImportHTMLData, ValidateHTMLData } from '../../wailsjs/go/main/App'
import type { models } from '../../wailsjs/go/models'
import ImportDialog from '../components/ImportDialog.vue'

// Reactive data
const loading = ref(true)
const error = ref<string | null>(null)
const salesRecords = ref<models.SalesRecord[]>([])
const filteredRecords = ref<models.SalesRecord[]>([])

// Filters and search
const searchQuery = ref('')
const selectedYear = ref<string>('all')
const selectedMonth = ref<string>('all')
const showImportDialog = ref(false)

// Pagination
const currentPage = ref(1)
const itemsPerPage = ref(25)
const totalPages = computed(() => Math.ceil(filteredRecords.value.length / itemsPerPage.value))

// Available years and months
const availableYears = computed(() => {
  const years = new Set<string>()
  salesRecords.value.forEach(record => {
    const year = new Date(record.date).getFullYear().toString()
    years.add(year)
  })
  return Array.from(years).sort().reverse()
})

const availableMonths = [
  { value: 'all', label: 'All Months' },
  { value: '1', label: 'January' },
  { value: '2', label: 'February' },
  { value: '3', label: 'March' },
  { value: '4', label: 'April' },
  { value: '5', label: 'May' },
  { value: '6', label: 'June' },
  { value: '7', label: 'July' },
  { value: '8', label: 'August' },
  { value: '9', label: 'September' },
  { value: '10', label: 'October' },
  { value: '11', label: 'November' },
  { value: '12', label: 'December' }
]

// Computed properties
const paginatedRecords = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredRecords.value.slice(start, end)
})

const totalSales = computed(() => {
  return filteredRecords.value.reduce((sum, record) => sum + record.sale_price, 0)
})

const totalCommission = computed(() => {
  return filteredRecords.value.reduce((sum, record) => sum + record.commission, 0)
})

const totalRemaining = computed(() => {
  return filteredRecords.value.reduce((sum, record) => sum + record.remaining, 0)
})

// Methods
const loadSalesData = async () => {
  try {
    loading.value = true
    error.value = null

    // Load recent imports (we'll expand this to load all records)
    const records = await GetRecentImports(1000) // Get up to 1000 records
    salesRecords.value = records || []
    
    applyFilters()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load sales data'
    console.error('Sales data loading error:', err)
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  let filtered = [...salesRecords.value]

  // Apply search filter (fuzzy matching)
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase().trim()
    filtered = filtered.filter(record => {
      const searchText = [
        record.description,
        record.store,
        record.vendor
      ].join(' ').toLowerCase()
      
      // Full match
      if (searchText.includes(query)) return true
      
      // Partial match (each word)
      const queryWords = query.split(' ')
      return queryWords.every(word => searchText.includes(word))
    })
  }

  // Apply year filter
  if (selectedYear.value !== 'all') {
    filtered = filtered.filter(record => {
      const recordYear = new Date(record.date).getFullYear().toString()
      return recordYear === selectedYear.value
    })
  }

  // Apply month filter
  if (selectedMonth.value !== 'all') {
    filtered = filtered.filter(record => {
      const recordMonth = (new Date(record.date).getMonth() + 1).toString()
      return recordMonth === selectedMonth.value
    })
  }

  filteredRecords.value = filtered
  currentPage.value = 1 // Reset to first page when filters change
}

const clearFilters = () => {
  searchQuery.value = ''
  selectedYear.value = 'all'
  selectedMonth.value = 'all'
  applyFilters()
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount)
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const editRecord = (record: models.SalesRecord) => {
  // TODO: Implement edit functionality
  console.log('Edit record:', record)
  alert('Edit functionality will be implemented in a future update')
}

const deleteRecord = (record: models.SalesRecord) => {
  // TODO: Implement delete functionality
  if (confirm(`Are you sure you want to delete the record for "${record.description}"?`)) {
    console.log('Delete record:', record)
    alert('Delete functionality will be implemented in a future update')
  }
}

const openImportDialog = () => {
  showImportDialog.value = true
}

const handleImportSuccess = () => {
  showImportDialog.value = false
  loadSalesData() // Refresh data after import
}

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

// Watchers
watch([searchQuery, selectedYear, selectedMonth], () => {
  applyFilters()
}, { deep: true })

// Lifecycle
onMounted(() => {
  loadSalesData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Sales Details</h1>
        <p class="text-gray-600 mt-1">Manage and analyze your sales records</p>
      </div>
      <div class="flex items-center space-x-3">
        <button
          @click="loadSalesData"
          :disabled="loading"
          class="btn btn-secondary flex items-center space-x-2"
        >
          <span :class="loading ? 'animate-spin' : ''">🔄</span>
          <span>Refresh</span>
        </button>
        <button
          @click="openImportDialog"
          class="btn btn-primary flex items-center space-x-2"
        >
          <span>📥</span>
          <span>Import Data</span>
        </button>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0 lg:space-x-4">
        <!-- Search -->
        <div class="flex-1 max-w-md">
          <label for="search" class="block text-sm font-medium text-gray-700 mb-2">
            Search Products
          </label>
          <div class="relative">
            <input
              id="search"
              v-model="searchQuery"
              type="text"
              placeholder="Search by product name, store, or vendor..."
              class="input w-full pl-10"
            />
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <span class="text-gray-400">🔍</span>
            </div>
          </div>
        </div>

        <!-- Filters -->
        <div class="flex items-end space-x-4">
          <div>
            <label for="year-filter" class="block text-sm font-medium text-gray-700 mb-2">
              Year
            </label>
            <select id="year-filter" v-model="selectedYear" class="input">
              <option value="all">All Years</option>
              <option v-for="year in availableYears" :key="year" :value="year">
                {{ year }}
              </option>
            </select>
          </div>

          <div>
            <label for="month-filter" class="block text-sm font-medium text-gray-700 mb-2">
              Month
            </label>
            <select id="month-filter" v-model="selectedMonth" class="input">
              <option v-for="month in availableMonths" :key="month.value" :value="month.value">
                {{ month.label }}
              </option>
            </select>
          </div>

          <button
            @click="clearFilters"
            class="btn btn-secondary"
          >
            Clear Filters
          </button>
        </div>
      </div>

      <!-- Summary Stats -->
      <div class="mt-6 grid grid-cols-1 sm:grid-cols-4 gap-4">
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="text-sm text-gray-600">Total Records</div>
          <div class="text-xl font-bold text-gray-900">{{ filteredRecords.length.toLocaleString() }}</div>
        </div>
        <div class="bg-primary-50 rounded-lg p-4">
          <div class="text-sm text-primary-600">Total Sales</div>
          <div class="text-xl font-bold text-primary-900">{{ formatCurrency(totalSales) }}</div>
        </div>
        <div class="bg-success-50 rounded-lg p-4">
          <div class="text-sm text-success-600">Total Commission</div>
          <div class="text-xl font-bold text-success-900">{{ formatCurrency(totalCommission) }}</div>
        </div>
        <div class="bg-warning-50 rounded-lg p-4">
          <div class="text-sm text-warning-600">Total Remaining</div>
          <div class="text-xl font-bold text-warning-900">{{ formatCurrency(totalRemaining) }}</div>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-error-50 border border-error-200 rounded-lg p-4">
      <div class="flex items-center space-x-2">
        <span class="text-error-600">⚠️</span>
        <div>
          <h3 class="font-medium text-error-800">Error Loading Sales Data</h3>
          <p class="text-error-600 text-sm mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <div class="animate-pulse space-y-4">
        <div class="h-4 bg-gray-200 rounded w-1/4"></div>
        <div class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-12 bg-gray-200 rounded"></div>
        </div>
      </div>
    </div>

    <!-- Sales Data Table -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Store
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Vendor
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Description
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Sale Price
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Commission
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Remaining
              </th>
              <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="paginatedRecords.length === 0">
              <td colspan="8" class="px-6 py-12 text-center text-gray-500">
                <div class="flex flex-col items-center space-y-2">
                  <span class="text-4xl">📭</span>
                  <p class="text-lg font-medium">No sales records found</p>
                  <p class="text-sm">Try adjusting your filters or import some data</p>
                </div>
              </td>
            </tr>
            <tr v-for="record in paginatedRecords" :key="record.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDate(record.date) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ record.store }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ record.vendor }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-900 max-w-xs truncate" :title="record.description">
                {{ record.description }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right font-medium">
                {{ formatCurrency(record.sale_price) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right">
                {{ formatCurrency(record.commission) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right">
                {{ formatCurrency(record.remaining) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-center">
                <div class="flex items-center justify-center space-x-2">
                  <button
                    @click="editRecord(record)"
                    class="text-primary-600 hover:text-primary-800 transition-colors"
                    title="Edit record"
                  >
                    ✏️
                  </button>
                  <button
                    @click="deleteRecord(record)"
                    class="text-error-600 hover:text-error-800 transition-colors"
                    title="Delete record"
                  >
                    🗑️
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="bg-gray-50 px-6 py-3 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Showing {{ ((currentPage - 1) * itemsPerPage) + 1 }} to 
            {{ Math.min(currentPage * itemsPerPage, filteredRecords.length) }} of 
            {{ filteredRecords.length }} results
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage === 1"
              class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            <div class="flex items-center space-x-1">
              <button
                v-for="page in Math.min(5, totalPages)"
                :key="page"
                @click="goToPage(page)"
                :class="[
                  'px-3 py-1 text-sm border rounded-md',
                  page === currentPage
                    ? 'bg-primary-600 text-white border-primary-600'
                    : 'border-gray-300 hover:bg-gray-100'
                ]"
              >
                {{ page }}
              </button>
            </div>
            <button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Import Dialog -->
    <ImportDialog
      v-if="showImportDialog"
      @close="showImportDialog = false"
      @success="handleImportSuccess"
    />
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */
</style>
