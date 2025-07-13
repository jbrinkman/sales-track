<script lang="ts" setup>
import { ref, computed } from 'vue'
import { ImportHTMLData, ValidateHTMLData, ImportHTMLDataWithOptions } from '../../wailsjs/go/main/App'
import type { main } from '../../wailsjs/go/models'

// Props and emits
const emit = defineEmits<{
  close: []
  success: []
}>()

// Reactive data
const htmlData = ref('')
const loading = ref(false)
const validating = ref(false)
const importResult = ref<main.ImportResult | null>(null)
const validationResult = ref<main.ValidationResult | null>(null)
const error = ref<string | null>(null)
const step = ref<'input' | 'preview' | 'result'>('input')

// Import options
const useConsignableFormat = ref(false)
const useBatchImport = ref(true)
const strictMode = ref(false)

// Computed properties
const isValidData = computed(() => {
  return htmlData.value.trim().length > 0
})

const canProceed = computed(() => {
  return isValidData.value && !loading.value && !validating.value
})

const hasValidationErrors = computed(() => {
  return validationResult.value && (!validationResult.value.valid || (validationResult.value.errors && validationResult.value.errors.length > 0))
})

// Methods
const validateData = async () => {
  if (!isValidData.value) return

  try {
    validating.value = true
    error.value = null
    
    const result = await ValidateHTMLData(htmlData.value)
    validationResult.value = result
    
    if (result.valid) {
      step.value = 'preview'
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Validation failed'
    console.error('Validation error:', err)
  } finally {
    validating.value = false
  }
}

const importData = async () => {
  if (!canProceed.value) return

  try {
    loading.value = true
    error.value = null

    let result: main.ImportResult

    if (useConsignableFormat.value || strictMode.value || useBatchImport.value) {
      // Use options-based import
      const options = {
        use_consignable_format: useConsignableFormat.value,
        strict_mode: strictMode.value,
        use_batch_import: useBatchImport.value,
        custom_column_mapping: []
      }
      result = await ImportHTMLDataWithOptions(htmlData.value, options)
    } else {
      // Use basic import
      result = await ImportHTMLData(htmlData.value)
    }

    importResult.value = result
    step.value = 'result'

    if (result.success) {
      // Auto-close after successful import
      setTimeout(() => {
        emit('success')
      }, 2000)
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Import failed'
    console.error('Import error:', err)
  } finally {
    loading.value = false
  }
}

const resetDialog = () => {
  htmlData.value = ''
  importResult.value = null
  validationResult.value = null
  error.value = null
  step.value = 'input'
  useConsignableFormat.value = false
  useBatchImport.value = true
  strictMode.value = false
}

const closeDialog = () => {
  emit('close')
}

const goBack = () => {
  if (step.value === 'preview') {
    step.value = 'input'
  } else if (step.value === 'result') {
    step.value = 'input'
    resetDialog()
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount)
}
</script>

<template>
  <!-- Modal Overlay -->
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200">
        <div>
          <h2 class="text-xl font-semibold text-gray-900">Import Sales Data</h2>
          <p class="text-sm text-gray-600 mt-1">
            <span v-if="step === 'input'">Paste HTML table data to import</span>
            <span v-else-if="step === 'preview'">Review data before importing</span>
            <span v-else>Import completed</span>
          </p>
        </div>
        <button
          @click="closeDialog"
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <span class="text-xl">×</span>
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 overflow-y-auto max-h-[calc(90vh-200px)]">
        <!-- Step 1: Input -->
        <div v-if="step === 'input'" class="space-y-6">
          <!-- Instructions -->
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h3 class="font-medium text-blue-800 mb-2">How to Import Data</h3>
            <ol class="text-sm text-blue-700 space-y-1 list-decimal list-inside">
              <li>Copy HTML table data from your sales website</li>
              <li>Paste the data into the text area below</li>
              <li>Configure import options if needed</li>
              <li>Click "Validate Data" to preview</li>
            </ol>
          </div>

          <!-- HTML Data Input -->
          <div>
            <label for="html-data" class="block text-sm font-medium text-gray-700 mb-2">
              HTML Table Data
            </label>
            <textarea
              id="html-data"
              v-model="htmlData"
              rows="12"
              placeholder="Paste your HTML table data here..."
              class="input w-full font-mono text-sm"
            ></textarea>
          </div>

          <!-- Import Options -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="font-medium text-gray-800 mb-3">Import Options</h3>
            <div class="space-y-3">
              <label class="flex items-center space-x-2">
                <input
                  v-model="useConsignableFormat"
                  type="checkbox"
                  class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span class="text-sm text-gray-700">Use Consignable format (headerless rows)</span>
              </label>
              <label class="flex items-center space-x-2">
                <input
                  v-model="useBatchImport"
                  type="checkbox"
                  class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span class="text-sm text-gray-700">Use batch import (faster for large datasets)</span>
              </label>
              <label class="flex items-center space-x-2">
                <input
                  v-model="strictMode"
                  type="checkbox"
                  class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span class="text-sm text-gray-700">Strict mode (fail on any parsing errors)</span>
              </label>
            </div>
          </div>
        </div>

        <!-- Step 2: Preview -->
        <div v-else-if="step === 'preview'" class="space-y-6">
          <!-- Validation Results -->
          <div v-if="validationResult" class="space-y-4">
            <!-- Success Summary -->
            <div v-if="validationResult.valid" class="bg-success-50 border border-success-200 rounded-lg p-4">
              <div class="flex items-center space-x-2">
                <span class="text-success-600">✅</span>
                <div>
                  <h3 class="font-medium text-success-800">Data Validation Successful</h3>
                  <p class="text-success-600 text-sm mt-1">
                    Found {{ validationResult.valid_rows }} valid records out of {{ validationResult.total_rows }} total rows
                  </p>
                </div>
              </div>
            </div>

            <!-- Error Summary -->
            <div v-if="hasValidationErrors" class="bg-warning-50 border border-warning-200 rounded-lg p-4">
              <div class="flex items-center space-x-2">
                <span class="text-warning-600">⚠️</span>
                <div>
                  <h3 class="font-medium text-warning-800">Validation Warnings</h3>
                  <p class="text-warning-600 text-sm mt-1">
                    {{ validationResult.invalid_rows }} rows have issues but import can continue
                  </p>
                </div>
              </div>
            </div>

            <!-- Statistics -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div class="bg-gray-50 rounded-lg p-3">
                <div class="text-sm text-gray-600">Total Rows</div>
                <div class="text-lg font-bold text-gray-900">{{ validationResult.total_rows }}</div>
              </div>
              <div class="bg-success-50 rounded-lg p-3">
                <div class="text-sm text-success-600">Valid Rows</div>
                <div class="text-lg font-bold text-success-900">{{ validationResult.valid_rows }}</div>
              </div>
              <div class="bg-warning-50 rounded-lg p-3">
                <div class="text-sm text-warning-600">Invalid Rows</div>
                <div class="text-lg font-bold text-warning-900">{{ validationResult.invalid_rows }}</div>
              </div>
              <div class="bg-primary-50 rounded-lg p-3">
                <div class="text-sm text-primary-600">Processing Time</div>
                <div class="text-lg font-bold text-primary-900">{{ validationResult.processing_time }}ns</div>
              </div>
            </div>

            <!-- Column Mapping -->
            <div v-if="validationResult.column_mapping" class="bg-gray-50 rounded-lg p-4">
              <h3 class="font-medium text-gray-800 mb-3">Detected Columns</h3>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-2 text-sm">
                <div v-for="(index, column) in validationResult.column_mapping" :key="column" class="flex justify-between">
                  <span class="text-gray-600">{{ column }}:</span>
                  <span class="font-medium">Column {{ index + 1 }}</span>
                </div>
              </div>
            </div>

            <!-- Errors -->
            <div v-if="validationResult.errors && validationResult.errors.length > 0" class="bg-error-50 border border-error-200 rounded-lg p-4">
              <h3 class="font-medium text-error-800 mb-3">Validation Errors</h3>
              <div class="space-y-2 max-h-40 overflow-y-auto">
                <div v-for="(error, index) in validationResult.errors" :key="index" class="text-sm">
                  <span class="font-medium text-error-700">Row {{ error.row }}:</span>
                  <span class="text-error-600 ml-2">{{ error.message }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 3: Result -->
        <div v-else-if="step === 'result'" class="space-y-6">
          <div v-if="importResult">
            <!-- Success Result -->
            <div v-if="importResult.success" class="bg-success-50 border border-success-200 rounded-lg p-6 text-center">
              <div class="text-success-600 text-4xl mb-4">🎉</div>
              <h3 class="text-lg font-semibold text-success-800 mb-2">Import Successful!</h3>
              <p class="text-success-600 mb-4">
                Successfully imported {{ importResult.imported_rows }} records
              </p>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                <div>
                  <div class="text-success-600">Total Rows</div>
                  <div class="font-bold text-success-900">{{ importResult.total_rows }}</div>
                </div>
                <div>
                  <div class="text-success-600">Parsed Rows</div>
                  <div class="font-bold text-success-900">{{ importResult.parsed_rows }}</div>
                </div>
                <div>
                  <div class="text-success-600">Imported Rows</div>
                  <div class="font-bold text-success-900">{{ importResult.imported_rows }}</div>
                </div>
              </div>
            </div>

            <!-- Partial Success Result -->
            <div v-else-if="importResult.imported_rows > 0" class="bg-warning-50 border border-warning-200 rounded-lg p-6">
              <div class="flex items-center space-x-2 mb-4">
                <span class="text-warning-600 text-2xl">⚠️</span>
                <div>
                  <h3 class="font-semibold text-warning-800">Partial Import</h3>
                  <p class="text-warning-600 text-sm">{{ importResult.error_message }}</p>
                </div>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                <div>
                  <div class="text-warning-600">Imported</div>
                  <div class="font-bold text-warning-900">{{ importResult.imported_rows }}</div>
                </div>
                <div>
                  <div class="text-warning-600">Failed</div>
                  <div class="font-bold text-warning-900">{{ importResult.parsed_rows - importResult.imported_rows }}</div>
                </div>
                <div>
                  <div class="text-warning-600">Success Rate</div>
                  <div class="font-bold text-warning-900">
                    {{ Math.round((importResult.imported_rows / importResult.parsed_rows) * 100) }}%
                  </div>
                </div>
              </div>
            </div>

            <!-- Failure Result -->
            <div v-else class="bg-error-50 border border-error-200 rounded-lg p-6 text-center">
              <div class="text-error-600 text-4xl mb-4">❌</div>
              <h3 class="text-lg font-semibold text-error-800 mb-2">Import Failed</h3>
              <p class="text-error-600">{{ importResult.error_message || 'Unknown error occurred' }}</p>
            </div>
          </div>
        </div>

        <!-- Error Display -->
        <div v-if="error" class="bg-error-50 border border-error-200 rounded-lg p-4">
          <div class="flex items-center space-x-2">
            <span class="text-error-600">⚠️</span>
            <div>
              <h3 class="font-medium text-error-800">Error</h3>
              <p class="text-error-600 text-sm mt-1">{{ error }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between p-6 border-t border-gray-200 bg-gray-50">
        <div class="flex items-center space-x-2">
          <button
            v-if="step !== 'input'"
            @click="goBack"
            class="btn btn-secondary"
          >
            ← Back
          </button>
        </div>

        <div class="flex items-center space-x-3">
          <button
            @click="closeDialog"
            class="btn btn-secondary"
          >
            {{ step === 'result' && importResult?.success ? 'Close' : 'Cancel' }}
          </button>

          <button
            v-if="step === 'input'"
            @click="validateData"
            :disabled="!canProceed"
            class="btn btn-primary flex items-center space-x-2"
          >
            <span v-if="validating" class="animate-spin">⏳</span>
            <span>{{ validating ? 'Validating...' : 'Validate Data' }}</span>
          </button>

          <button
            v-else-if="step === 'preview'"
            @click="importData"
            :disabled="!canProceed"
            class="btn btn-primary flex items-center space-x-2"
          >
            <span v-if="loading" class="animate-spin">⏳</span>
            <span>{{ loading ? 'Importing...' : 'Import Data' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */
</style>
