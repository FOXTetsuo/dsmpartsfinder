<template>
    <div class="parts-page">
        <n-space vertical :size="24">
            <!-- Page Header -->
            <n-card>
                <n-space vertical :size="12">
                    <div>
                        <h1
                            style="
                                font-size: 28px;
                                font-weight: bold;
                                margin: 0;
                            "
                        >
                            Parts Inventory
                        </h1>
                        <p style="color: #666; margin-top: 8px">
                            View and manage all parts from registered sites
                        </p>
                    </div>

                    <n-space align="center">
                        <n-button
                            type="primary"
                            :loading="fetchingAll"
                            @click="fetchFromAllSites"
                            size="large"
                        >
                            <template #icon>
                                <n-icon>
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <path
                                            d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"
                                        ></path>
                                        <polyline
                                            points="3.27 6.96 12 12.01 20.73 6.96"
                                        ></polyline>
                                        <line
                                            x1="12"
                                            y1="22.08"
                                            x2="12"
                                            y2="12"
                                        ></line>
                                    </svg>
                                </n-icon>
                            </template>
                            Fetch Parts from All Sites
                        </n-button>

                        <div
                            style="display: flex; align-items: center; gap: 8px"
                        >
                            <span style="color: #666; font-size: 14px"
                                >Limit:</span
                            >
                            <n-input-number
                                v-model:value="fetchLimit"
                                :min="10"
                                :max="10000"
                                :step="100"
                                style="width: 120px"
                                size="small"
                            />
                            <span style="color: #999; font-size: 12px"
                                >parts</span
                            >
                        </div>

                        <n-button
                            :loading="loading"
                            @click="loadParts"
                            size="large"
                        >
                            <template #icon>
                                <n-icon>
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <path
                                            d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"
                                        />
                                    </svg>
                                </n-icon>
                            </template>
                            Refresh
                        </n-button>
                    </n-space>

                    <n-alert
                        v-if="fetchAllSuccess"
                        type="success"
                        closable
                        @close="fetchAllSuccess = false"
                    >
                        {{ fetchAllSuccessMessage }}
                    </n-alert>

                    <n-alert
                        v-if="fetchAllError"
                        type="warning"
                        closable
                        @close="fetchAllError = false"
                    >
                        {{ fetchAllErrorMessage }}
                    </n-alert>
                </n-space>
            </n-card>

            <!-- Parts Statistics -->
            <n-card v-if="!loading && parts.length > 0" title="Statistics">
                <n-space :size="24">
                    <n-statistic label="Total Parts" :value="parts.length" />
                    <n-statistic
                        label="Unique Sites"
                        :value="uniqueSitesCount"
                    />
                </n-space>
            </n-card>

            <!-- Parts Table -->
            <n-card title="Parts List">
                <n-data-table
                    :columns="columns"
                    :data="parts"
                    :loading="loading"
                    :pagination="pagination"
                    :bordered="false"
                    :single-line="false"
                    striped
                />
            </n-card>
        </n-space>
    </div>
</template>

<script>
import { defineComponent, h, ref, computed, onMounted } from "vue";
import {
    NSpace,
    NCard,
    NButton,
    NDataTable,
    NIcon,
    NImage,
    NTag,
    NStatistic,
    NAlert,
    NInputNumber,
    useMessage,
} from "naive-ui";
import axios from "axios";

export default defineComponent({
    name: "Parts",
    components: {
        NSpace,
        NCard,
        NButton,
        NDataTable,
        NIcon,
        NStatistic,
        NAlert,
        NInputNumber,
    },
    setup() {
        console.log("[Parts.vue] setup() called");
        const message = useMessage();
        const parts = ref([]);
        const loading = ref(false);
        const fetchingAll = ref(false);
        const fetchAllSuccess = ref(false);
        const fetchAllSuccessMessage = ref("");
        const fetchAllError = ref(false);
        const fetchAllErrorMessage = ref("");
        const fetchLimit = ref(3000); // Configurable limit for fetching parts

        console.log("[Parts.vue] Initial parts.value:", parts.value);

        const pagination = {
            pageSize: 20,
            showSizePicker: true,
            pageSizes: [10, 20, 50, 100],
            showQuickJumper: true,
        };

        const uniqueSitesCount = computed(() => {
            const siteIds = new Set(parts.value.map((part) => part.site_id));
            console.log(
                "[Parts.vue] uniqueSitesCount computed. Unique sites:",
                siteIds.size,
                "Site IDs:",
                Array.from(siteIds),
            );
            return siteIds.size;
        });

        const columns = [
            {
                title: "ID",
                key: "id",
                width: 80,
                sorter: (a, b) => a.id - b.id,
            },
            {
                title: "Image",
                key: "image_base64",
                width: 100,
                render(row) {
                    if (row.image_base64) {
                        return h(NImage, {
                            width: 60,
                            height: 60,
                            src: `data:image/jpeg;base64,${row.image_base64}`,
                            objectFit: "cover",
                            style: { borderRadius: "4px" },
                        });
                    }
                    return h(
                        "div",
                        {
                            style: {
                                width: "60px",
                                height: "60px",
                                backgroundColor: "#f0f0f0",
                                borderRadius: "4px",
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                color: "#999",
                            },
                        },
                        "No image",
                    );
                },
            },
            {
                title: "Part ID",
                key: "part_id",
                width: 120,
                ellipsis: {
                    tooltip: true,
                },
            },
            {
                title: "Name",
                key: "name",
                width: 150,
                ellipsis: {
                    tooltip: true,
                },
            },
            {
                title: "Description",
                key: "description",
                minWidth: 200,
                ellipsis: {
                    tooltip: true,
                },
            },
            {
                title: "Price",
                key: "price",
                width: 120,
                render(row) {
                    if (row.price) {
                        return h(
                            NTag,
                            {
                                type: "success",
                                size: "small",
                                strong: true,
                            },
                            { default: () => row.price },
                        );
                    }
                    return h("span", { style: { color: "#999" } }, "-");
                },
            },
            {
                title: "Type",
                key: "type_name",
                width: 120,
                ellipsis: {
                    tooltip: true,
                },
            },
            {
                title: "Actions",
                key: "actions",
                width: 120,
                render(row) {
                    return h(
                        NButton,
                        {
                            size: "small",
                            type: "primary",
                            text: true,
                            tag: "a",
                            href: row.url,
                            target: "_blank",
                            rel: "noopener noreferrer",
                        },
                        { default: () => "View Source" },
                    );
                },
            },
        ];

        const loadParts = async () => {
            console.log("[Parts.vue] loadParts() called");
            loading.value = true;
            try {
                console.log(
                    "[Parts.vue] Making GET request to /api/parts with limit=10000, offset=0",
                );
                const response = await axios.get("/api/parts", {
                    params: {
                        limit: 10000,
                        offset: 0,
                    },
                });
                console.log("[Parts.vue] Received response:", response);
                console.log("[Parts.vue] Response data:", response.data);
                console.log(
                    "[Parts.vue] Response data.data:",
                    response.data.data,
                );

                parts.value = response.data.data || [];
                console.log("[Parts.vue] Set parts.value to:", parts.value);
                console.log(
                    "[Parts.vue] parts.value.length:",
                    parts.value.length,
                );

                if (parts.value.length > 0) {
                    console.log("[Parts.vue] First part:", parts.value[0]);
                } else {
                    console.warn("[Parts.vue] No parts received from API");
                }

                message.success(`Loaded ${parts.value.length} parts`);
            } catch (error) {
                console.error("[Parts.vue] Error loading parts:", error);
                console.error("[Parts.vue] Error response:", error.response);
                console.error(
                    "[Parts.vue] Error response data:",
                    error.response?.data,
                );
                message.error(
                    error.response?.data?.error ||
                        error.message ||
                        "Failed to load parts",
                );
            } finally {
                loading.value = false;
                console.log(
                    "[Parts.vue] loadParts() finished. Final parts count:",
                    parts.value.length,
                );
            }
        };

        const fetchFromAllSites = async () => {
            console.log("[Parts.vue] fetchFromAllSites() called");
            if (
                !confirm(
                    "This will fetch parts from all registered sites. This may take a while. Continue?",
                )
            ) {
                console.log("[Parts.vue] User cancelled fetch operation");
                return;
            }

            console.log("[Parts.vue] Starting fetch from all sites");
            fetchingAll.value = true;
            fetchAllSuccess.value = false;
            fetchAllError.value = false;

            try {
                const requestBody = {
                    year_from: 1960,
                    year_to: 2025,
                    limit: fetchLimit.value,
                };
                console.log(
                    "[Parts.vue] Making POST request to /api/parts/fetch-all with body:",
                    requestBody,
                );

                const response = await axios.post(
                    "/api/parts/fetch-all",
                    requestBody,
                );

                console.log(
                    "[Parts.vue] Received response from fetch-all:",
                    response,
                );
                console.log("[Parts.vue] Response data:", response.data);

                const data = response.data;
                console.log("[Parts.vue] Extracted data.total:", data.total);
                console.log("[Parts.vue] Extracted data.sites:", data.sites);
                console.log(
                    "[Parts.vue] Extracted data.data length:",
                    data.data?.length,
                );

                fetchAllSuccessMessage.value = `Successfully fetched ${data.total} parts from ${data.sites} site(s)`;
                fetchAllSuccess.value = true;

                if (data.errors && Object.keys(data.errors).length > 0) {
                    const errorSites = Object.keys(data.errors).join(", ");
                    fetchAllErrorMessage.value = `Some sites had errors (Site IDs: ${errorSites}). Check console for details.`;
                    fetchAllError.value = true;
                    console.error(
                        "[Parts.vue] Fetch errors from sites:",
                        data.errors,
                    );
                }

                message.success(`Fetched ${data.total} parts successfully!`);

                // Reload parts after fetching
                console.log("[Parts.vue] Reloading parts after fetch...");
                await loadParts();
            } catch (error) {
                console.error(
                    "[Parts.vue] Error fetching from all sites:",
                    error,
                );
                console.error("[Parts.vue] Error response:", error.response);
                console.error(
                    "[Parts.vue] Error response data:",
                    error.response?.data,
                );
                fetchAllErrorMessage.value =
                    error.response?.data?.error ||
                    error.message ||
                    "Failed to fetch parts from sites";
                fetchAllError.value = true;
                message.error("Failed to fetch parts from all sites");
            } finally {
                fetchingAll.value = false;
                console.log("[Parts.vue] fetchFromAllSites() finished");
            }
        };

        onMounted(() => {
            console.log("[Parts.vue] Component mounted, loading parts...");
            loadParts();
        });

        return {
            parts,
            loading,
            fetchingAll,
            columns,
            pagination,
            uniqueSitesCount,
            loadParts,
            fetchFromAllSites,
            fetchAllSuccess,
            fetchAllSuccessMessage,
            fetchAllError,
            fetchAllErrorMessage,
            fetchLimit,
        };
    },
});
</script>

<style scoped>
.parts-page {
    padding: 24px;
    max-width: 1600px;
    margin: 0 auto;
}

h1 {
    margin: 0;
}
</style>
