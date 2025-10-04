<template>
  <div class="parts-page">
    <!-- Page Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">DSM Parts Catalog</h1>
      <p class="text-lg text-gray-600">Browse and manage your DSM vehicle parts</p>
    </div>

    <!-- Add New Part Section -->
    <div class="card mb-8">
      <div class="card-body">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-900">Manage Parts</h2>
          <button
            @click="showAddForm = !showAddForm"
            class="btn btn-primary"
          >
            {{ showAddForm ? 'Cancel' : 'Add New Part' }}
          </button>
        </div>

        <!-- Add Part Form -->
        <div v-if="showAddForm" class="border-t border-gray-200 pt-6 mt-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Add New Part</h3>
          <form @submit.prevent="addPart" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label for="name" class="form-label">Part Name</label>
                <input
                  type="text"
                  id="name"
                  v-model="newPart.name"
                  required
                  placeholder="Enter part name"
                  class="form-input"
                />
              </div>
              <div>
                <label for="price" class="form-label">Price</label>
                <input
                  type="text"
                  id="price"
                  v-model="newPart.price"
                  required
                  placeholder="$19.99"
                  class="form-input"
                />
              </div>
            </div>
            <div>
              <label for="description" class="form-label">Description</label>
              <textarea
                id="description"
                v-model="newPart.description"
                required
                placeholder="Enter part description"
                rows="3"
                class="form-textarea"
              ></textarea>
            </div>
            <div class="flex justify-end space-x-3">
              <button
                type="button"
                @click="showAddForm = false"
                class="btn btn-outline-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="btn btn-success"
              >
                <span v-if="submitting" class="spinner w-4 h-4 mr-2"></span>
                {{ submitting ? 'Adding...' : 'Add Part' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Parts List Section -->
    <div class="card">
      <div class="card-header">
        <h2 class="text-xl font-semibold text-gray-900">Parts Inventory</h2>
        <p class="text-sm text-gray-600 mt-1">
          {{ parts.length }} part{{ parts.length !== 1 ? 's' : '' }} found
        </p>
      </div>

      <div class="card-body">
        <!-- Loading State -->
        <div v-if="loading" class="text-center py-12">
          <div class="spinner w-8 h-8 mx-auto mb-4 border-primary-600"></div>
          <p class="text-gray-600">Loading parts...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-12">
          <div class="alert alert-error inline-flex items-center">
            <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
            </svg>
            Error loading parts: {{ error }}
          </div>
          <button @click="fetchParts" class="btn btn-secondary mt-4">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
            </svg>
            Retry
          </button>
        </div>

        <!-- Empty State -->
        <div v-else-if="parts.length === 0" class="text-center py-12">
          <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M9 9h.01M15 9h.01"></path>
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">No parts found</h3>
          <p class="text-gray-600 mb-4">Add your first part to get started!</p>
          <button @click="showAddForm = true" class="btn btn-primary">
            Add Your First Part
          </button>
        </div>

        <!-- Parts Grid -->
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="part in parts"
            :key="part.id"
            class="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden"
          >
            <!-- Part Header -->
            <div class="p-6 pb-4">
              <div class="flex justify-between items-start mb-3">
                <h3 class="text-lg font-semibold text-gray-900 truncate pr-2">
                  {{ part.name }}
                </h3>
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                  #{{ part.id }}
                </span>
              </div>
              <p class="text-sm text-gray-600 line-clamp-3 mb-4">
                {{ part.description }}
              </p>
            </div>

            <!-- Part Footer -->
            <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-between items-center">
              <div class="flex items-center">
                <span class="text-lg font-bold text-green-600">
                  {{ formatPrice(part.price) }}
                </span>
              </div>
              <button
                @click="viewPart(part.id)"
                class="btn btn-outline text-sm"
              >
                View Details
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Toast -->
    <Transition name="fade">
      <div
        v-if="successMessage"
        class="fixed top-4 right-4 z-50 alert alert-success shadow-lg animate-slide-in"
      >
        <div class="flex items-center">
          <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
          </svg>
          {{ successMessage }}
        </div>
      </div>
    </Transition>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Parts',
  data() {
    return {
      parts: [],
      loading: false,
      error: null,
      showAddForm: false,
      submitting: false,
      successMessage: '',
      newPart: {
        name: '',
        description: '',
        price: ''
      }
    }
  },
  mounted() {
    this.fetchParts()
  },
  methods: {
    async fetchParts() {
      this.loading = true
      this.error = null

      try {
        const response = await axios.get('/api/v1/parts')
        this.parts = response.data.data || []
      } catch (error) {
        this.error = error.response?.data?.message || error.message || 'Failed to fetch parts'
        console.error('Error fetching parts:', error)
      } finally {
        this.loading = false
      }
    },

    async addPart() {
      if (!this.newPart.name || !this.newPart.description || !this.newPart.price) {
        this.showNotification('Please fill in all fields', 'error')
        return
      }

      this.submitting = true

      try {
        const response = await axios.post('/api/v1/parts', this.newPart)

        // Add the new part to the list
        if (response.data.data) {
          this.parts.push(response.data.data)
        }

        // Reset form
        this.newPart = {
          name: '',
          description: '',
          price: ''
        }

        this.showAddForm = false
        this.showSuccessMessage('Part added successfully!')
      } catch (error) {
        this.error = error.response?.data?.message || error.message || 'Failed to add part'
        console.error('Error adding part:', error)
      } finally {
        this.submitting = false
      }
    },

    async viewPart(partId) {
      try {
        const response = await axios.get(`/api/v1/parts/${partId}`)
        const part = response.data.data

        // Simple alert for demo - in a real app, you'd use a modal or navigate to a detail page
        alert(`Part Details:\n\nName: ${part.name}\nDescription: ${part.description}\nPrice: ${this.formatPrice(part.price)}\nID: #${part.id}`)
      } catch (error) {
        console.error('Error fetching part details:', error)
        alert('Error fetching part details: ' + (error.response?.data?.message || error.message))
      }
    },

    formatPrice(price) {
      if (typeof price === 'string' && price.startsWith('$')) {
        return price
      }
      return typeof price === 'number' ? `$${price.toFixed(2)}` : `$${price}`
    },

    showSuccessMessage(message) {
      this.successMessage = message
      setTimeout(() => {
        this.successMessage = ''
      }, 3000)
    },

    showNotification(message, type) {
      // If parent provides notification system, use it
      if (this.$parent && this.$parent.showNotification) {
        this.$parent.showNotification(message, type)
      } else {
        // Fallback to simple alert
        alert(message)
      }
    }
  }
}
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
