<template>
    <n-message-provider>
        <div
            id="app"
            class="min-h-screen bg-gray-50"
            @click="handleGlobalClick"
        >
            <!-- Navigation Header -->
            <nav class="bg-white shadow-sm border-b border-gray-200">
                <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div class="flex justify-between h-16">
                        <div class="flex">
                            <!-- Logo/Brand -->
                            <div class="flex-shrink-0 flex items-center">
                                <router-link
                                    to="/"
                                    class="text-2xl font-bold bg-gradient-to-r from-red-600 to-red-800 bg-clip-text text-transparent hover:from-red-700 hover:to-red-900 transition-colors duration-200"
                                >
                                    DSM Parts Finder
                                </router-link>
                            </div>

                            <!-- Navigation Links -->
                            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
                                <router-link
                                    to="/"
                                    class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200"
                                    :class="{
                                        'border-primary-500 text-primary-600':
                                            $route.name === 'Home',
                                    }"
                                >
                                    Home
                                </router-link>

                                <router-link
                                    to="/parts"
                                    class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200"
                                    :class="{
                                        'border-primary-500 text-primary-600':
                                            $route.name === 'Parts',
                                    }"
                                >
                                    Parts
                                </router-link>
                                <router-link
                                    to="/browse"
                                    class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200"
                                    :class="{
                                        'border-primary-500 text-primary-600':
                                            $route.name === 'Browse',
                                    }"
                                >
                                    Browse
                                </router-link>
                            </div>
                        </div>

                        <!-- Mobile menu button -->
                        <div class="sm:hidden flex items-center">
                            <button
                                @click="mobileMenuOpen = !mobileMenuOpen"
                                type="button"
                                class="bg-white inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
                                aria-controls="mobile-menu"
                                :aria-expanded="mobileMenuOpen"
                            >
                                <span class="sr-only">Open main menu</span>
                                <!-- Menu icon -->
                                <svg
                                    class="h-6 w-6"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                    v-if="!mobileMenuOpen"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M4 6h16M4 12h16M4 18h16"
                                    />
                                </svg>
                                <!-- Close icon -->
                                <svg
                                    class="h-6 w-6"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                    v-else
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Mobile menu -->
                <div class="sm:hidden" id="mobile-menu" v-show="mobileMenuOpen">
                    <div class="pt-2 pb-3 space-y-1">
                        <router-link
                            to="/"
                            @click="mobileMenuOpen = false"
                            class="bg-primary-50 border-primary-500 text-primary-700 block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                            :class="
                                $route.name === 'Home'
                                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                                    : 'border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300'
                            "
                        >
                            Home
                        </router-link>

                        <router-link
                            to="/parts"
                            @click="mobileMenuOpen = false"
                            class="border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300 block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                            :class="
                                $route.name === 'Parts'
                                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                                    : 'border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300'
                            "
                        >
                            Parts
                        </router-link>
                        <router-link
                            to="/browse"
                            @click="mobileMenuOpen = false"
                            class="border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300 block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                            :class="
                                $route.name === 'Browse'
                                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                                    : 'border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300'
                            "
                        >
                            Browse
                        </router-link>
                    </div>
                </div>
            </nav>

            <!-- Main Content -->
            <main class="flex-1">
                <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                    <router-view />
                </div>
            </main>

            <!-- Footer -->
            <footer class="bg-white border-t border-gray-200 mt-12">
                <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                    <div
                        class="flex flex-col items-center justify-center space-y-2"
                    >
                        <p class="text-sm text-gray-500">
                            © 2024 DSM Parts Finder. Built by a guy who was too
                            fucking lazy to keep searching for new parts in
                            junkyards.
                        </p>
                        <p class="text-xs text-gray-400">
                            Find and manage parts for your DSM vehicle
                        </p>
                    </div>
                </div>
            </footer>

            <!-- Global Loading Overlay -->
            <div
                v-if="loading"
                class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center"
            >
                <div
                    class="bg-white p-6 rounded-lg shadow-xl flex flex-col items-center space-y-4"
                >
                    <div
                        class="spinner w-8 h-8 border-4 border-primary-600"
                    ></div>
                    <p class="text-gray-700">Loading...</p>
                </div>
            </div>

            <!-- Toast Notifications -->
            <div class="fixed top-4 right-4 z-50 space-y-2">
                <div
                    v-for="notification in notifications"
                    :key="notification.id"
                    class="alert max-w-sm animate-slide-in"
                    :class="[
                        notification.type === 'success' ? 'alert-success' : '',
                        notification.type === 'error' ? 'alert-error' : '',
                        notification.type === 'warning' ? 'alert-warning' : '',
                        notification.type === 'info' ? 'alert-info' : '',
                    ]"
                >
                    <div class="flex justify-between items-start">
                        <p class="text-sm">{{ notification.message }}</p>
                        <button
                            @click="removeNotification(notification.id)"
                            class="ml-2 flex-shrink-0 text-current opacity-50 hover:opacity-75"
                        >
                            <svg
                                class="w-4 h-4"
                                fill="currentColor"
                                viewBox="0 0 20 20"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                    clip-rule="evenodd"
                                ></path>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </n-message-provider>
</template>

<script>
import { NMessageProvider } from "naive-ui";

export default {
    name: "App",
    components: {
        NMessageProvider,
    },
    data() {
        return {
            mobileMenuOpen: false,
            loading: false,
            notifications: [],
        };
    },
    methods: {
        showNotification(message, type = "info", duration = 5000) {
            const id = Date.now();
            this.notifications.push({
                id,
                message,
                type,
            });

            if (duration > 0) {
                setTimeout(() => {
                    this.removeNotification(id);
                }, duration);
            }
        },

        removeNotification(id) {
            this.notifications = this.notifications.filter((n) => n.id !== id);
        },

        setLoading(loading) {
            this.loading = loading;
        },
    },
    provide() {
        return {
            showNotification: this.showNotification,
            setLoading: this.setLoading,
        };
    },
    mounted() {
        // Close mobile menu when clicking outside
        document.addEventListener("click", (event) => {
            const nav = document.querySelector("nav");
            if (nav && !nav.contains(event.target)) {
                this.mobileMenuOpen = false;
            }
        });
    },
};
</script>

<style scoped>
.router-link-active {
    @apply border-primary-500 text-primary-600;
}
</style>
