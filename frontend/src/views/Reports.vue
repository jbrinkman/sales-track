<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { GetRecentImports } from '../../wailsjs/go/main/App'
import type { models } from '../../wailsjs/go/models'

// Reactive data
const loading = ref(true)
const error = ref<string | null>(null)
const salesRecords = ref<models.SalesRecord[]>([])
const selectedReport = ref('pivot-table')
const selectedYear = ref<string>('all')

// Report types
const reportTypes = [
  {
    value: 'pivot-table',
    label: 'Pivot Table Report',
    description: 'Hierarchical view by Year > Month > Date'
  },
  {
    value: 'summary',
    label: 'Summary Report',
    description: 'High-level overview and totals'
  },
  {
    value: 'trends',
    label: 'Trends Analysis',
    description: 'Sales trends over time (Coming Soon)'
  }
]

// Pivot table data structure
interface PivotData {
  year: string
  months: {
    month: string
    monthName: string
    days: {
      day: string
      records: models.SalesRecord[]
      totalSales: number
      totalCommission: number
      totalRemaining: number
      itemCount: number
    }[]
    totalSales: number
    totalCommission: number
    totalRemaining: number
    itemCount: number
    expanded: boolean
  }[]
  totalSales: number
  totalCommission: number
  totalRemaining: number
  itemCount: number
  expanded: boolean
}

const pivotData = ref<PivotData[]>([])
const expandedItems = ref<Set<string>>(new Set())

// Computed properties
const availableYears = computed(() => {
  const years = new Set<string>()
  salesRecords.value.forEach(record => {
    const year = new Date(record.date).getFullYear().toString()
    years.add(year)
  })
  return Array.from(years).sort().reverse()
})

const filteredRecords = computed(() => {
  if (selectedYear.value === 'all') {
    return salesRecords.value
  }
  return salesRecords.value.filter(record => {
    const recordYear = new Date(record.date).getFullYear().toString()
    return recordYear === selectedYear.value
  })
})

const summaryStats = computed(() => {
  const records = filteredRecords.value
  return {
    totalRecords: records.length,
    totalSales: records.reduce((sum, r) => sum + r.sale_price, 0),
    totalCommission: records.reduce((sum, r) => sum + r.commission, 0),
    totalRemaining: records.reduce((sum, r) => sum + r.remaining, 0),
    averageSale: records.length > 0 ? records.reduce((sum, r) => sum + r.sale_price, 0) / records.length : 0,
    uniqueStores: new Set(records.map(r => r.store)).size,
    uniqueVendors: new Set(records.map(r => r.vendor)).size
  }
})

// Methods
const loadReportsData = async () => {
  try {
    loading.value = true
    error.value = null

    // Load sales records (we'll expand this to use proper reporting API later)
    const records = await GetRecentImports(1000)
    salesRecords.value = records || []
    
    generatePivotData()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load reports data'
    console.error('Reports data loading error:', err)
  } finally {
    loading.value = false
  }
}

const generatePivotData = () => {
  const records = filteredRecords.value
  const yearGroups = new Map<string, models.SalesRecord[]>()

  // Group by year
  records.forEach(record => {
    const year = new Date(record.date).getFullYear().toString()
    if (!yearGroups.has(year)) {
      yearGroups.set(year, [])
    }
    yearGroups.get(year)!.push(record)
  })

  // Build pivot structure
  pivotData.value = Array.from(yearGroups.entries())
    .map(([year, yearRecords]) => {
      const monthGroups = new Map<string, models.SalesRecord[]>()
      
      // Group by month
      yearRecords.forEach(record => {
        const date = new Date(record.date)
        const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
        if (!monthGroups.has(monthKey)) {
          monthGroups.set(monthKey, [])
        }
        monthGroups.get(monthKey)!.push(record)
      })

      const months = Array.from(monthGroups.entries())
        .map(([monthKey, monthRecords]) => {
          const dayGroups = new Map<string, models.SalesRecord[]>()
          
          // Group by day
          monthRecords.forEach(record => {
            const dayKey = record.date.split('T')[0] // Get YYYY-MM-DD part
            if (!dayGroups.has(dayKey)) {
              dayGroups.set(dayKey, [])
            }
            dayGroups.get(dayKey)!.push(record)
          })

          const days = Array.from(dayGroups.entries())
            .map(([day, dayRecords]) => ({
              day,
              records: dayRecords,
              totalSales: dayRecords.reduce((sum, r) => sum + r.sale_price, 0),
              totalCommission: dayRecords.reduce((sum, r) => sum + r.commission, 0),
              totalRemaining: dayRecords.reduce((sum, r) => sum + r.remaining, 0),
              itemCount: dayRecords.length
            }))
            .sort((a, b) => a.day.localeCompare(b.day))

          const monthDate = new Date(monthKey + '-01')
          return {
            month: monthKey,
            monthName: monthDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }),
            days,
            totalSales: monthRecords.reduce((sum, r) => sum + r.sale_price, 0),
            totalCommission: monthRecords.reduce((sum, r) => sum + r.commission, 0),
            totalRemaining: monthRecords.reduce((sum, r) => sum + r.remaining, 0),
            itemCount: monthRecords.length,
            expanded: false
          }
        })
        .sort((a, b) => a.month.localeCompare(b.month))

      return {
        year,
        months,
        totalSales: yearRecords.reduce((sum, r) => sum + r.sale_price, 0),
        totalCommission: yearRecords.reduce((sum, r) => sum + r.commission, 0),
        totalRemaining: yearRecords.reduce((sum, r) => sum + r.remaining, 0),
        itemCount: yearRecords.length,
        expanded: false
      }
    })
    .sort((a, b) => b.year.localeCompare(a.year)) // Most recent first
}

const toggleExpanded = (type: 'year' | 'month', yearIndex: number, monthIndex?: number) => {
  if (type === 'year') {
    pivotData.value[yearIndex].expanded = !pivotData.value[yearIndex].expanded
  } else if (type === 'month' && monthIndex !== undefined) {
    pivotData.value[yearIndex].months[monthIndex].expanded = !pivotData.value[yearIndex].months[monthIndex].expanded
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount)
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric'
  })
}

// Watchers
const refreshData = () => {
  generatePivotData()
}

// Lifecycle
onMounted(() => {
  loadReportsData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Reports</h1>
        <p class="text-gray-600 mt-1">Analyze your sales data with interactive reports</p>
      </div>
      <button
        @click="loadReportsData"
        :disabled="loading"
        class="btn btn-primary flex items-center space-x-2"
      >
        <span :class="loading ? 'animate-spin' : ''">🔄</span>
        <span>Refresh</span>
      </button>
    </div>

    <!-- Report Controls -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0 lg:space-x-4">
        <!-- Report Type Selection -->
        <div class="flex-1">
          <label for="report-type" class="block text-sm font-medium text-gray-700 mb-2">
            Report Type
          </label>
          <select id="report-type" v-model="selectedReport" class="input max-w-md">
            <option v-for="report in reportTypes" :key="report.value" :value="report.value">
              {{ report.label }}
            </option>
          </select>
          <p class="text-sm text-gray-500 mt-1">
            {{ reportTypes.find(r => r.value === selectedReport)?.description }}
          </p>
        </div>

        <!-- Year Filter -->
        <div>
          <label for="year-filter" class="block text-sm font-medium text-gray-700 mb-2">
            Year Filter
          </label>
          <select id="year-filter" v-model="selectedYear" @change="refreshData" class="input">
            <option value="all">All Years</option>
            <option v-for="year in availableYears" :key="year" :value="year">
              {{ year }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-error-50 border border-error-200 rounded-lg p-4">
      <div class="flex items-center space-x-2">
        <span class="text-error-600">⚠️</span>
        <div>
          <h3 class="font-medium text-error-800">Error Loading Reports</h3>
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

    <!-- Summary Report -->
    <div v-else-if="selectedReport === 'summary'" class="space-y-6">
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="card bg-gradient-to-br from-primary-50 to-primary-100 border-primary-200">
          <h3 class="text-sm font-medium text-primary-700 mb-1">Total Sales</h3>
          <p class="text-2xl font-bold text-primary-900">{{ formatCurrency(summaryStats.totalSales) }}</p>
          <p class="text-xs text-primary-600 mt-2">{{ summaryStats.totalRecords }} transactions</p>
        </div>
        <div class="card bg-gradient-to-br from-success-50 to-success-100 border-success-200">
          <h3 class="text-sm font-medium text-success-700 mb-1">Total Commission</h3>
          <p class="text-2xl font-bold text-success-900">{{ formatCurrency(summaryStats.totalCommission) }}</p>
          <p class="text-xs text-success-600 mt-2">Average: {{ formatCurrency(summaryStats.totalCommission / Math.max(summaryStats.totalRecords, 1)) }}</p>
        </div>
        <div class="card bg-gradient-to-br from-warning-50 to-warning-100 border-warning-200">
          <h3 class="text-sm font-medium text-warning-700 mb-1">Total Remaining</h3>
          <p class="text-2xl font-bold text-warning-900">{{ formatCurrency(summaryStats.totalRemaining) }}</p>
          <p class="text-xs text-warning-600 mt-2">Outstanding balance</p>
        </div>
        <div class="card bg-gradient-to-br from-purple-50 to-purple-100 border-purple-200">
          <h3 class="text-sm font-medium text-purple-700 mb-1">Average Sale</h3>
          <p class="text-2xl font-bold text-purple-900">{{ formatCurrency(summaryStats.averageSale) }}</p>
          <p class="text-xs text-purple-600 mt-2">Per transaction</p>
        </div>
      </div>

      <!-- Additional Stats -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">Business Insights</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="text-center">
            <div class="text-3xl font-bold text-gray-900">{{ summaryStats.uniqueStores }}</div>
            <div class="text-sm text-gray-600">Active Stores</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-gray-900">{{ summaryStats.uniqueVendors }}</div>
            <div class="text-sm text-gray-600">Vendor Partners</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-gray-900">
              {{ Math.round((summaryStats.totalCommission / Math.max(summaryStats.totalSales, 1)) * 100) }}%
            </div>
            <div class="text-sm text-gray-600">Commission Rate</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pivot Table Report -->
    <div v-else-if="selectedReport === 'pivot-table'" class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="p-6 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900">Pivot Table Report</h2>
        <p class="text-sm text-gray-600 mt-1">Hierarchical view of sales data by Year → Month → Date</p>
      </div>

      <div v-if="pivotData.length === 0" class="p-12 text-center text-gray-500">
        <div class="flex flex-col items-center space-y-2">
          <span class="text-4xl">📊</span>
          <p class="text-lg font-medium">No data available</p>
          <p class="text-sm">Import some sales data to see reports</p>
        </div>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Period
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Items Sold
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Total Sales
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Commission
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Remaining
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <!-- Year Level -->
            <template v-for="(yearData, yearIndex) in pivotData" :key="yearData.year">
              <tr class="bg-gray-50 hover:bg-gray-100 cursor-pointer" @click="toggleExpanded('year', yearIndex)">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center space-x-2">
                    <span class="text-gray-400">{{ yearData.expanded ? '▼' : '▶' }}</span>
                    <span class="font-semibold text-gray-900">📅 {{ yearData.year }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right font-medium text-gray-900">
                  {{ yearData.itemCount.toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right font-medium text-gray-900">
                  {{ formatCurrency(yearData.totalSales) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-gray-900">
                  {{ formatCurrency(yearData.totalCommission) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-gray-900">
                  {{ formatCurrency(yearData.totalRemaining) }}
                </td>
              </tr>

              <!-- Month Level -->
              <template v-if="yearData.expanded">
                <template v-for="(monthData, monthIndex) in yearData.months" :key="monthData.month">
                  <tr class="bg-blue-50 hover:bg-blue-100 cursor-pointer" @click="toggleExpanded('month', yearIndex, monthIndex)">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center space-x-2 ml-6">
                        <span class="text-blue-400">{{ monthData.expanded ? '▼' : '▶' }}</span>
                        <span class="font-medium text-blue-900">📊 {{ monthData.monthName }}</span>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right font-medium text-blue-900">
                      {{ monthData.itemCount.toLocaleString() }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right font-medium text-blue-900">
                      {{ formatCurrency(monthData.totalSales) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-blue-900">
                      {{ formatCurrency(monthData.totalCommission) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-blue-900">
                      {{ formatCurrency(monthData.totalRemaining) }}
                    </td>
                  </tr>

                  <!-- Day Level -->
                  <template v-if="monthData.expanded">
                    <tr v-for="dayData in monthData.days" :key="dayData.day" class="hover:bg-gray-50">
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="flex items-center space-x-2 ml-12">
                          <span class="text-green-600">📈</span>
                          <span class="text-gray-700">{{ formatDate(dayData.day) }}</span>
                        </div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-gray-900">
                        {{ dayData.itemCount.toLocaleString() }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-gray-900">
                        {{ formatCurrency(dayData.totalSales) }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-gray-900">
                        {{ formatCurrency(dayData.totalCommission) }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-gray-900">
                        {{ formatCurrency(dayData.totalRemaining) }}
                      </td>
                    </tr>
                  </template>
                </template>
              </template>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Coming Soon Reports -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
      <div class="text-6xl mb-4">🚧</div>
      <h2 class="text-xl font-semibold text-gray-900 mb-2">Coming Soon</h2>
      <p class="text-gray-600">This report type is under development and will be available in a future update.</p>
    </div>
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */
</style>
