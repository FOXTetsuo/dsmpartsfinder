<template>
    <div class="sites-page">
        <!-- Page Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-bold text-gray-900 mb-2">
                DSM Parts Sites
            </h1>
            <p class="text-lg text-gray-600">
                Manage sites where you source DSM parts
            </p>
        </div>

        <!-- Add New Site Section -->
        <div class="card mb-8">
            <div class="card-body">
                <div class="flex justify-between items-center mb-4">
                    <h2 class="text-xl font-semibold text-gray-900">
                        Manage Sites
                    </h2>
                    <button
                        @click="showAddForm = !showAddForm"
                        class="btn btn-primary"
                    >
                        {{ showAddForm ? "Cancel" : "Add New Site" }}
                    </button>
                </div>

                <!-- Add Site Form -->
                <div
                    v-if="showAddForm"
                    class="border-t border-gray-200 pt-6 mt-6"
                >
                    <h3 class="text-lg font-medium text-gray-900 mb-4">
                        {{ editingSite ? "Edit Site" : "Add New Site" }}
                    </h3>
                    <form @submit.prevent="saveSite" class="space-y-4">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label for="name" class="form-label"
                                    >Site Name</label
                                >
                                <input
                                    type="text"
                                    id="name"
                                    v-model="newSite.name"
                                    required
                                    placeholder="Enter site name"
                                    class="form-input"
                                />
                            </div>
                            <div>
                                <label for="url" class="form-label"
                                    >Site URL</label
                                >
                                <input
                                    type="url"
                                    id="url"
                                    v-model="newSite.url"
                                    required
                                    placeholder="https://example.com"
                                    class="form-input"
                                />
                            </div>
                        </div>
                        <div class="flex justify-end space-x-3">
                            <button
                                type="button"
                                @click="cancelEdit"
                                class="btn btn-outline-secondary"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="submitting"
                                class="btn btn-success"
                            >
                                <span
                                    v-if="submitting"
                                    class="spinner w-4 h-4 mr-2"
                                ></span>
                                {{
                                    submitting
                                        ? editingSite
                                            ? "Updating..."
                                            : "Adding..."
                                        : editingSite
                                          ? "Update Site"
                                          : "Add Site"
                                }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>

        <!-- Sites List Section -->
        <div class="card">
            <div class="card-header">
                <h2 class="text-xl font-semibold text-gray-900">
                    Sites Inventory
                </h2>
                <p class="text-sm text-gray-600 mt-1">
                    {{ sites.length }} site{{ sites.length !== 1 ? "s" : "" }}
                    found
                </p>
            </div>

            <div class="card-body">
                <!-- Loading State -->
                <div v-if="loading" class="text-center py-12">
                    <div
                        class="spinner w-8 h-8 mx-auto mb-4 border-primary-600"
                    ></div>
                    <p class="text-gray-600">Loading sites...</p>
                </div>

                <!-- Error State -->
                <div v-else-if="error" class="text-center py-12">
                    <div class="alert alert-error inline-flex items-center">
                        <svg
                            class="w-5 h-5 mr-2"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                clip-rule="evenodd"
                            ></path>
                        </svg>
                        Error loading sites: {{ error }}
                    </div>
                    <button @click="fetchSites" class="btn btn-secondary mt-4">
                        <svg
                            class="w-4 h-4 mr-2"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                            ></path>
                        </svg>
                        Retry
                    </button>
                </div>

                <!-- Empty State -->
                <div v-else-if="sites.length === 0" class="text-center py-12">
                    <svg
                        class="w-16 h-16 text-gray-300 mx-auto mb-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="1"
                            d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
                        ></path>
                    </svg>
                    <h3 class="text-lg font-medium text-gray-900 mb-2">
                        No sites found
                    </h3>
                    <p class="text-gray-600 mb-4">
                        Add your first site to get started!
                    </p>
                    <button @click="showAddForm = true" class="btn btn-primary">
                        Add Your First Site
                    </button>
                </div>

                <!-- Sites Grid -->
                <div
                    v-else
                    class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6"
                >
                    <div
                        v-for="site in sites"
                        :key="site.id"
                        class="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden"
                    >
                        <!-- Site Header -->
                        <div class="p-6 pb-4">
                            <div class="flex justify-between items-start mb-3">
                                <h3
                                    class="text-lg font-semibold text-gray-900 truncate pr-2"
                                >
                                    {{ site.name }}
                                </h3>
                                <span
                                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
                                >
                                    #{{ site.id }}
                                </span>
                            </div>
                            <a
                                :href="site.url"
                                target="_blank"
                                rel="noopener noreferrer"
                                class="text-sm text-blue-600 hover:text-blue-800 hover:underline break-all"
                            >
                                {{ site.url }}
                            </a>
                        </div>

                        <!-- Site Footer -->
                        <div
                            class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-between items-center"
                        >
                            <button
                                @click="editSite(site)"
                                class="btn btn-outline text-sm"
                            >
                                Edit
                            </button>
                            <button
                                @click="deleteSite(site.id)"
                                class="btn btn-outline text-sm text-red-600 hover:bg-red-50"
                            >
                                Delete
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
                    <svg
                        class="w-5 h-5 mr-2"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                            clip-rule="evenodd"
                        ></path>
                    </svg>
                    {{ successMessage }}
                </div>
            </div>
        </Transition>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "Sites",
    data() {
        return {
            sites: [],
            loading: false,
            error: null,
            showAddForm: false,
            submitting: false,
            successMessage: "",
            editingSite: null,
            newSite: {
                name: "",
                url: "",
            },
        };
    },
    mounted() {
        this.fetchSites();
    },
    methods: {
        async fetchSites() {
            this.loading = true;
            this.error = null;

            try {
                const response = await axios.get("/api/sites");
                this.sites = response.data.data || [];
            } catch (error) {
                this.error =
                    error.response?.data?.error ||
                    error.response?.data?.message ||
                    error.message ||
                    "Failed to fetch sites";
                console.error("Error fetching sites:", error);
            } finally {
                this.loading = false;
            }
        },

        async saveSite() {
            this.submitting = true;
            this.error = null;

            try {
                if (this.editingSite) {
                    // Update existing site
                    await axios.put(`/api/sites/${this.editingSite.id}`, {
                        name: this.newSite.name,
                        url: this.newSite.url,
                    });
                    this.showSuccessMessage("Site updated successfully!");
                } else {
                    // Create new site
                    await axios.post("/api/sites", {
                        name: this.newSite.name,
                        url: this.newSite.url,
                    });
                    this.showSuccessMessage("Site added successfully!");
                }

                // Refresh the sites list
                await this.fetchSites();

                // Reset form
                this.cancelEdit();
            } catch (error) {
                this.error =
                    error.response?.data?.error ||
                    error.response?.data?.message ||
                    error.message ||
                    "Failed to save site";
                console.error("Error saving site:", error);
                alert("Error: " + this.error);
            } finally {
                this.submitting = false;
            }
        },

        editSite(site) {
            this.editingSite = site;
            this.newSite = {
                name: site.name,
                url: site.url,
            };
            this.showAddForm = true;
        },

        cancelEdit() {
            this.showAddForm = false;
            this.editingSite = null;
            this.newSite = {
                name: "",
                url: "",
            };
        },

        async deleteSite(siteId) {
            if (!confirm("Are you sure you want to delete this site?")) {
                return;
            }

            try {
                await axios.delete(`/api/sites/${siteId}`);
                this.showSuccessMessage("Site deleted successfully!");
                await this.fetchSites();
            } catch (error) {
                this.error =
                    error.response?.data?.error ||
                    error.response?.data?.message ||
                    error.message ||
                    "Failed to delete site";
                console.error("Error deleting site:", error);
                alert("Error: " + this.error);
            }
        },

        showSuccessMessage(message) {
            this.successMessage = message;
            setTimeout(() => {
                this.successMessage = "";
            }, 3000);
        },
    },
};
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
