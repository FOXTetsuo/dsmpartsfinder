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
                <h1 class="text-4xl sm:text-6xl font-bold mb-6 text-shado w-lg">
                    DSM Parts Finder
                </h1>

                <div
                    class="space-y-4 sm:space-y-0 sm:space-x-4 sm:flex sm:justify-center"
                ></div>
            </div>
        </div>

        <!-- Features Section -->
        <div id="features" class="py-16 bg-white">
            <div class="max-w-6xl mx-auto px-4">
                <div class="text-center mb-12">
                    <h2
                        class="text-3xl sm:text-4xl font-bold text-gray-900 mb-4"
                    >
                        Why should you use DSM Parts Finder?
                    </h2>
                    <p
                        class="text-xl text-gray-600 max-w-2xl mb-4 mx-auto italic"
                    >
                        "Why haven't you hung up that light yet?!" "Are you
                        looking for car parts AGAIN?" "COME TO BED!"
                    </p>
                    <p class="text-xl text-gray-600 max-w-2xl mb-8 mx-auto">
                        Are these types of questions driving you crazy? Do you
                        just want to find new parts for your 1995-2000
                        Mitsubishi Eclipse?
                    </p>
                    <h2
                        class="text-3xl sm:text-4xl font-bold text-gray-900 mb-4"
                    >
                        We've got you covered!
                    </h2>
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
                            Automatic part fetching
                        </h3>
                        <p class="text-gray-600">
                            Parts from connected sites are fetched every hour,
                            ensuring you always have the latest greatest at your
                            fingertips.
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
                            Done with stupidly slow websites that don't work
                            properly? So am I! With DSM Parts Finder, you can
                            search for parts quickly and easily, without having
                            to wait for slow websites to load or deal with
                            broken links.
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
                            Reliability
                        </h3>
                        <p class="text-gray-600">
                            This gorgeous website is so reliable, even your
                            missus will be happy. Sure, you're still looking at
                            parts all day, but now it's 10x as fast, and it
                            never stops working!
                        </p>
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
                        to="/browse"
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
