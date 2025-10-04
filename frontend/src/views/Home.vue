<script>
import axios from "axios";

export default {
    name: "Home",
    data() {
        return {
            apiStatus: {
                healthy: false,
                message: "Checking...",
            },
        };
    },
    mounted() {
        this.checkApiStatus();
    },
    methods: {
        scrollToFeatures() {
            document.getElementById("features").scrollIntoView({
                behavior: "smooth",
            });
        },

        async checkApiStatus() {
            try {
                const response = await axios.get("/api/health");
                this.apiStatus.healthy = true;
                this.apiStatus.message =
                    response.data.message || "API is running";
            } catch (error) {
                this.apiStatus.healthy = false;
                this.apiStatus.message = "Unable to connect to API";
                console.error("API health check failed:", error);
            }
        },

        showContactInfo() {
            alert(
                "Contact Information:\n\nThis is a demo application.\nIn a real app, this would show contact details or open a contact form.",
            );
        },
    },
};
</script>

<template>
    <div class="home">
        <!-- Hero Section -->
        <div
            class="bg-gradient-to-br from-primary-600 to-primary-800 text-white"
        >
            <div class="max-w-4xl mx-auto px-4 py-16 sm:py-24 text-center">
                <h1 class="text-4xl sm:text-6xl font-bold mb-6 text-shadow-lg">
                    DSM Parts Finder
                </h1>
                <p class="text-xl sm:text-2xl mb-8 text-primary-100 font-light">
                    Your ultimate destination for DSM vehicle parts and
                    accessories
                </p>
                <div
                    class="space-y-4 sm:space-y-0 sm:space-x-4 sm:flex sm:justify-center"
                >
                    <router-link
                        to="/parts"
                        class="btn bg-white text-primary-600 hover:bg-primary-50 font-semibold text-lg px-8 py-3 rounded-lg shadow-lg transition-all duration-300 hover:shadow-xl transform hover:-translate-y-1"
                    >
                        Browse Parts Catalog
                    </router-link>
                    <button
                        @click="scrollToFeatures"
                        class="btn bg-transparent border-2 border-white text-white hover:bg-white hover:text-primary-600 font-semibold text-lg px-8 py-3 rounded-lg transition-all duration-300"
                    >
                        Learn More
                    </button>
                </div>
            </div>
        </div>

        <!-- Features Section -->
        <div id="features" class="py-16 bg-white">
            <div class="max-w-6xl mx-auto px-4">
                <div class="text-center mb-12">
                    <h2
                        class="text-3xl sm:text-4xl font-bold text-gray-900 mb-4"
                    >
                        Why Choose DSM Parts Finder?
                    </h2>
                    <p class="text-xl text-gray-600 max-w-2xl mx-auto">
                        We specialize in Diamond Star Motors vehicles with
                        authentic parts and expert knowledge
                    </p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <!-- Feature 1 -->
                    <div
                        class="text-center p-6 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors duration-300"
                    >
                        <div
                            class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                            <svg
                                class="w-8 h-8 text-primary-600"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                ></path>
                            </svg>
                        </div>
                        <h3 class="text-xl font-semibold text-gray-900 mb-3">
                            Authentic Parts
                        </h3>
                        <p class="text-gray-600">
                            Genuine OEM and high-quality aftermarket parts
                            specifically for DSM vehicles including Eclipse,
                            Talon, and Laser.
                        </p>
                    </div>

                    <!-- Feature 2 -->
                    <div
                        class="text-center p-6 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors duration-300"
                    >
                        <div
                            class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                            <svg
                                class="w-8 h-8 text-primary-600"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M13 10V3L4 14h7v7l9-11h-7z"
                                ></path>
                            </svg>
                        </div>
                        <h3 class="text-xl font-semibold text-gray-900 mb-3">
                            Fast Search
                        </h3>
                        <p class="text-gray-600">
                            Quickly find the exact parts you need with our
                            advanced search and filtering system designed for
                            DSM enthusiasts.
                        </p>
                    </div>

                    <!-- Feature 3 -->
                    <div
                        class="text-center p-6 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors duration-300"
                    >
                        <div
                            class="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                            <svg
                                class="w-8 h-8 text-primary-600"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M12 2.25a9.75 9.75 0 1 0 0 19.5 9.75 9.75 0 0 0 0-19.5Z"
                                ></path>
                            </svg>
                        </div>
                        <h3 class="text-xl font-semibold text-gray-900 mb-3">
                            Expert Support
                        </h3>
                        <p class="text-gray-600">
                            Get help from DSM experts who understand these
                            unique vehicles and can guide you to the right
                            parts.
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Vehicle Models Section -->
        <div class="py-16 bg-gray-50">
            <div class="max-w-6xl mx-auto px-4">
                <div class="text-center mb-12">
                    <h2
                        class="text-3xl sm:text-4xl font-bold text-gray-900 mb-4"
                    >
                        Supported DSM Models
                    </h2>
                    <p class="text-xl text-gray-600">
                        We carry parts for all Diamond Star Motors vehicles
                    </p>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-3 gap-8">
                    <!-- Eclipse -->
                    <div
                        class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300"
                    >
                        <div
                            class="h-48 bg-gradient-to-br from-red-400 to-red-600 flex items-center justify-center"
                        >
                            <div class="text-center text-white">
                                <h3 class="text-2xl font-bold">
                                    Mitsubishi Eclipse
                                </h3>
                                <p class="text-red-100">1989-1999</p>
                            </div>
                        </div>
                        <div class="p-6">
                            <p class="text-gray-600 mb-4">
                                The iconic sports car that started it all. Find
                                parts for 1G and 2G Eclipse models.
                            </p>
                            <div class="flex flex-wrap gap-2">
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >1G (1989-1994)</span
                                >
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >2G (1995-1999)</span
                                >
                            </div>
                        </div>
                    </div>

                    <!-- Talon -->
                    <div
                        class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300"
                    >
                        <div
                            class="h-48 bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center"
                        >
                            <div class="text-center text-white">
                                <h3 class="text-2xl font-bold">Eagle Talon</h3>
                                <p class="text-blue-100">1989-1998</p>
                            </div>
                        </div>
                        <div class="p-6">
                            <p class="text-gray-600 mb-4">
                                Eagle's sports car variant of the DSM platform.
                                Performance parts and OEM replacements.
                            </p>
                            <div class="flex flex-wrap gap-2">
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >TSi</span
                                >
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >TSi AWD</span
                                >
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >ESi</span
                                >
                            </div>
                        </div>
                    </div>

                    <!-- Laser -->
                    <div
                        class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300"
                    >
                        <div
                            class="h-48 bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center"
                        >
                            <div class="text-center text-white">
                                <h3 class="text-2xl font-bold">
                                    Plymouth Laser
                                </h3>
                                <p class="text-green-100">1989-1994</p>
                            </div>
                        </div>
                        <div class="p-6">
                            <p class="text-gray-600 mb-4">
                                Plymouth's take on the DSM platform. Share parts
                                with Eclipse and Talon models.
                            </p>
                            <div class="flex flex-wrap gap-2">
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >Base</span
                                >
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >RS</span
                                >
                                <span
                                    class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                                    >RS Turbo</span
                                >
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- CTA Section -->
        <div class="bg-primary-600 text-white py-16">
            <div class="max-w-4xl mx-auto px-4 text-center">
                <h2 class="text-3xl sm:text-4xl font-bold mb-6">
                    Ready to Find Your Parts?
                </h2>
                <p class="text-xl text-primary-100 mb-8 max-w-2xl mx-auto">
                    Start browsing our extensive catalog of DSM parts or contact
                    our experts for personalized assistance.
                </p>
                <div
                    class="space-y-4 sm:space-y-0 sm:space-x-4 sm:flex sm:justify-center"
                >
                    <router-link
                        to="/parts"
                        class="btn bg-white text-primary-600 hover:bg-primary-50 font-semibold text-lg px-8 py-3 rounded-lg shadow-lg transition-all duration-300 hover:shadow-xl transform hover:-translate-y-1"
                    >
                        Browse Parts Now
                    </router-link>
                    <a
                        href="#"
                        @click.prevent="showContactInfo"
                        class="btn bg-transparent border-2 border-white text-white hover:bg-white hover:text-primary-600 font-semibold text-lg px-8 py-3 rounded-lg transition-all duration-300"
                    >
                        Contact Us
                    </a>
                </div>
            </div>
        </div>

        <!-- API Status -->
        <div class="py-8 bg-gray-100">
            <div class="max-w-4xl mx-auto px-4">
                <div class="bg-white rounded-lg shadow p-6">
                    <h3 class="text-lg font-semibold text-gray-900 mb-4">
                        System Status
                    </h3>
                    <div class="flex items-center space-x-4">
                        <div class="flex items-center">
                            <div
                                class="w-3 h-3 rounded-full mr-2"
                                :class="
                                    apiStatus.healthy
                                        ? 'bg-green-500'
                                        : 'bg-red-500'
                                "
                            ></div>
                            <span class="text-sm text-gray-600">
                                API:
                                {{ apiStatus.healthy ? "Online" : "Offline" }}
                            </span>
                        </div>
                        <div
                            v-if="apiStatus.message"
                            class="text-sm text-gray-500"
                        >
                            {{ apiStatus.message }}
                        </div>
                        <button
                            @click="checkApiStatus"
                            class="text-primary-600 hover:text-primary-700 text-sm font-medium"
                        >
                            Refresh Status
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.bg-clip-text {
    -webkit-background-clip: text;
    background-clip: text;
}

.gradient-primary {
    background: linear-gradient(
        135deg,
        theme("colors.primary.600") 0%,
        theme("colors.primary.700") 100%
    );
}

.text-shadow-lg {
    text-shadow:
        0 4px 8px rgba(0, 0, 0, 0.12),
        0 2px 4px rgba(0, 0, 0, 0.08);
}
</style>
