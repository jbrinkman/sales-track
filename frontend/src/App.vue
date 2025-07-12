<script lang="ts" setup>
import UnoTest from './components/UnoTest.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Header -->
    <header class="bg-primary-600 text-white shadow-lg">
      <div class="container mx-auto px-4 py-6">
        <h1 class="text-3xl font-bold">Sales Track</h1>
        <p class="text-primary-100 mt-2">Modern sales data tracking and reporting</p>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 container mx-auto px-4 py-8">
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-2xl font-semibold text-gray-800 mb-4">Welcome to Sales Track</h2>
        <p class="text-gray-600 mb-6">
          Replace your Excel-based sales tracking with this modern desktop application.
          Import HTML table data, generate reports, and analyze your sales data with ease.
        </p>
        
        <!-- Feature Cards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
          <div class="bg-blue-50 p-4 rounded-lg border border-blue-200">
            <h3 class="font-semibold text-blue-800 mb-2">HTML Data Import</h3>
            <p class="text-blue-600 text-sm">Paste HTML table data directly from websites</p>
          </div>
          <div class="bg-green-50 p-4 rounded-lg border border-green-200">
            <h3 class="font-semibold text-green-800 mb-2">SQLite Storage</h3>
            <p class="text-green-600 text-sm">Reliable local database storage</p>
          </div>
          <div class="bg-purple-50 p-4 rounded-lg border border-purple-200">
            <h3 class="font-semibold text-purple-800 mb-2">Pivot Reports</h3>
            <p class="text-purple-600 text-sm">Interactive reports with drill-down</p>
          </div>
        </div>

        <!-- UnoCSS Status -->
        <div class="mt-8 p-4 bg-gradient-to-r from-primary-500 to-primary-600 text-white rounded-lg">
          <p class="font-medium">✅ UnoCSS Configuration Complete!</p>
          <p class="text-primary-100 text-sm mt-1">Custom theme colors, shortcuts, and utility classes are active.</p>
        </div>
      </div>

      <!-- UnoCSS Test Component -->
      <div class="bg-white rounded-lg shadow-md">
        <UnoTest />
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-gray-800 text-gray-300 py-4">
      <div class="container mx-auto px-4 text-center">
        <p class="text-sm">Sales Track - Built with Go, Wails, Vue, and UnoCSS</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Custom styles if needed - UnoCSS handles most styling */
</style>
