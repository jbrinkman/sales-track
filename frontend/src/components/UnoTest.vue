<template>
  <div class="p-6 space-y-4">
    <h3 class="text-xl font-semibold text-gray-800 mb-4">UnoCSS Feature Test</h3>
    
    <!-- Button Tests -->
    <div class="space-x-2">
      <button class="btn-primary">Primary Button</button>
      <button class="btn-secondary">Secondary Button</button>
      <button class="btn bg-success-600 text-white hover:bg-success-700">Success Button</button>
    </div>

    <!-- Card Test -->
    <div class="card p-4">
      <h4 class="font-medium text-gray-900 mb-2">Sample Card</h4>
      <p class="text-gray-600 text-sm">This card uses the custom 'card' shortcut class.</p>
    </div>

    <!-- Input Test -->
    <div class="space-y-2">
      <label class="block text-sm font-medium text-gray-700">Sample Input</label>
      <input type="text" class="input w-full max-w-xs" placeholder="Test input field">
    </div>

    <!-- Color Palette Test -->
    <div class="grid grid-cols-5 gap-2">
      <div class="h-12 bg-primary-500 rounded flex items-center justify-center text-white text-xs">Primary</div>
      <div class="h-12 bg-success-500 rounded flex items-center justify-center text-white text-xs">Success</div>
      <div class="h-12 bg-warning-500 rounded flex items-center justify-center text-white text-xs">Warning</div>
      <div class="h-12 bg-error-500 rounded flex items-center justify-center text-white text-xs">Error</div>
      <div class="h-12 bg-gray-500 rounded flex items-center justify-center text-white text-xs">Gray</div>
    </div>

    <!-- Responsive Test -->
    <div class="bg-gradient-to-r from-primary-400 to-primary-600 p-4 rounded-lg text-white">
      <p class="text-sm md:text-base lg:text-lg">Responsive text sizing test</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
// Component for testing UnoCSS functionality
</script>
