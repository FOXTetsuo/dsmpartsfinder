<template>
    <n-config-provider :theme-overrides="themeOverrides">
        <div class="browse-page">
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
                                Browse Parts
                            </h1>
                            <p style="color: #666; margin-top: 8px">
                                Search and filter through all available parts
                            </p>
                        </div>
                    </n-space>
                </n-card>

                <!-- Search and Filters -->
                <n-card title="Search & Filters">
                    <n-space vertical :size="16">
                        <!-- Search Input -->
                        <n-input
                            v-model:value="searchQuery"
                            placeholder="Search by name, description, type, or part ID..."
                            size="large"
                            clearable
                            @input="debouncedSearch"
                        >
                            <template #prefix>
                                <n-icon>
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <circle cx="11" cy="11" r="8"></circle>
                                        <path d="m21 21-4.35-4.35"></path>
                                    </svg>
                                </n-icon>
                            </template>
                        </n-input>

                        <!-- Filter Row -->
                        <n-space :size="12" wrap>
                            <!-- Site Filter -->
                            <n-select
                                v-model:value="filters.siteId"
                                :options="siteOptions"
                                placeholder="All Sites"
                                clearable
                                style="width: 180px"
                                @update:value="applyFilters"
                            />

                            <!-- Type filter -->
                            <!-- <n-select
                                v-model:value="filters.typeName"
                                :options="typeOptions"
                                placeholder="All Types"
                                clearable
                                filterable
                                style="width: 200px"
                                @update:value="applyFilters"
                            /> -->

                            <!-- Sort By -->
                            <n-select
                                v-model:value="sortBy"
                                :options="sortOptions"
                                placeholder="Sort By"
                                style="width: 180px"
                                @update:value="applyFilters"
                            />

                            <!-- Only New Toggle -->
                            <n-checkbox
                                v-model:checked="filters.showOnlyNew"
                                @update:checked="applyFilters"
                                class="only-new-checkbox"
                            >
                                Only New items
                            </n-checkbox>

                            <!-- View Mode Toggle -->
                            <n-button-group>
                                <n-button
                                    :type="
                                        viewMode === 'grid'
                                            ? 'primary'
                                            : 'default'
                                    "
                                    @click="viewMode = 'grid'"
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
                                                <rect
                                                    x="3"
                                                    y="3"
                                                    width="7"
                                                    height="7"
                                                ></rect>
                                                <rect
                                                    x="14"
                                                    y="3"
                                                    width="7"
                                                    height="7"
                                                ></rect>
                                                <rect
                                                    x="14"
                                                    y="14"
                                                    width="7"
                                                    height="7"
                                                ></rect>
                                                <rect
                                                    x="3"
                                                    y="14"
                                                    width="7"
                                                    height="7"
                                                ></rect>
                                            </svg>
                                        </n-icon>
                                    </template>
                                </n-button>
                                <n-button
                                    :type="
                                        viewMode === 'list'
                                            ? 'primary'
                                            : 'default'
                                    "
                                    @click="viewMode = 'list'"
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
                                                <line
                                                    x1="8"
                                                    y1="6"
                                                    x2="21"
                                                    y2="6"
                                                ></line>
                                                <line
                                                    x1="8"
                                                    y1="12"
                                                    x2="21"
                                                    y2="12"
                                                ></line>
                                                <line
                                                    x1="8"
                                                    y1="18"
                                                    x2="21"
                                                    y2="18"
                                                ></line>
                                                <line
                                                    x1="3"
                                                    y1="6"
                                                    x2="3.01"
                                                    y2="6"
                                                ></line>
                                                <line
                                                    x1="3"
                                                    y1="12"
                                                    x2="3.01"
                                                    y2="12"
                                                ></line>
                                                <line
                                                    x1="3"
                                                    y1="18"
                                                    x2="3.01"
                                                    y2="18"
                                                ></line>
                                            </svg>
                                        </n-icon>
                                    </template>
                                </n-button>
                            </n-button-group>

                            <!-- Reset Filters -->
                            <n-button @click="resetFilters" quaternary>
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
                                                d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"
                                            ></path>
                                            <path d="M21 3v5h-5"></path>
                                            <path
                                                d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"
                                            ></path>
                                            <path d="M8 16H3v5"></path>
                                        </svg>
                                    </n-icon>
                                </template>
                                Reset Filters
                            </n-button>
                        </n-space>

                        <!-- Active Filters Display -->
                        <n-space v-if="hasActiveFilters" :size="8">
                            <n-text depth="3">Active filters:</n-text>
                            <n-tag
                                v-if="filters.siteId"
                                closable
                                @close="
                                    filters.siteId = null;
                                    applyFilters();
                                "
                            >
                                Site: {{ getSiteName(filters.siteId) }}
                            </n-tag>
                            <!-- <n-tag
                                v-if="filters.typeName"
                                closable
                                @close="
                                    filters.typeName = null;
                                    applyFilters();
                                "
                            >
                                Type: {{ filters.typeName }}
                            </n-tag> -->
                            <n-tag
                                v-if="searchQuery"
                                closable
                                @close="
                                    searchQuery = '';
                                    applyFilters();
                                "
                            >
                                Search: "{{ searchQuery }}"
                            </n-tag>
                        </n-space>
                    </n-space>
                </n-card>

                <!-- Results Summary -->
                <n-card v-if="!loading">
                    <n-space align="center" justify="space-between">
                        <n-statistic
                            label="Total Results"
                            :value="totalItems"
                        />
                        <n-text depth="3">
                            Showing {{ parts.length }} of {{ totalItems }} parts
                        </n-text>
                    </n-space>
                </n-card>

                <!-- Loading State -->
                <n-card v-if="loading">
                    <n-space
                        vertical
                        align="center"
                        justify="center"
                        :size="16"
                        style="padding: 40px"
                    >
                        <n-spin size="large" />
                        <n-text>Loading parts...</n-text>
                    </n-space>
                </n-card>

                <!-- Empty State -->
                <n-empty
                    v-else-if="totalItems === 0"
                    description="No parts found matching your criteria"
                    size="large"
                    style="padding: 60px 0"
                >
                    <template #extra>
                        <n-button @click="resetFilters">
                            Clear Filters
                        </n-button>
                    </template>
                </n-empty>

                <!-- Grid View -->
                <div v-else-if="viewMode === 'grid'" class="parts-grid">
                    <n-card
                        v-for="part in parts"
                        :key="part.id"
                        hoverable
                        class="part-card"
                        @click="selectPart(part)"
                    >
                        <!-- NEW Badge -->
                        <div v-if="isNewPart(part)" class="new-badge">
                            <span class="new-text">NEW!</span>
                        </div>

                        <div class="part-card-content">
                            <!-- Image -->
                            <div class="part-image">
                                <div
                                    v-if="part.image_base64"
                                    class="part-image-blur-bg"
                                    :style="{
                                        backgroundImage: `url('data:image/jpeg;base64,${part.image_base64}')`,
                                    }"
                                >
                                    <img
                                        :src="`data:image/jpeg;base64,${part.image_base64}`"
                                        class="part-image-centered"
                                        alt="Part Image"
                                    />
                                </div>
                                <div v-else class="no-image">
                                    <n-icon size="48" :depth="3">
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        >
                                            <rect
                                                x="3"
                                                y="3"
                                                width="18"
                                                height="18"
                                                rx="2"
                                                ry="2"
                                            ></rect>
                                            <circle
                                                cx="8.5"
                                                cy="8.5"
                                                r="1.5"
                                            ></circle>
                                            <polyline
                                                points="21 15 16 10 5 21"
                                            ></polyline>
                                        </svg>
                                    </n-icon>
                                    <n-text depth="3">No Image</n-text>
                                </div>
                            </div>

                            <!-- Info -->
                            <div class="part-info">
                                <n-ellipsis
                                    :line-clamp="2"
                                    :tooltip="{ placement: 'top' }"
                                >
                                    <n-text strong>{{ part.name }}</n-text>
                                </n-ellipsis>
                                <n-ellipsis
                                    :line-clamp="2"
                                    :tooltip="{ placement: 'top' }"
                                    style="margin-top: 4px"
                                >
                                    <n-text depth="3" style="font-size: 13px">{{
                                        part.description
                                    }}</n-text>
                                </n-ellipsis>
                                <n-space :size="8" style="margin-top: 8px" wrap>
                                    <!-- Site Name Tag -->
                                    <n-tag size="small" type="default">
                                        {{ getSiteName(part.site_id) }}
                                    </n-tag>
                                    <!-- <n-tag size="small" type="info">
                                        {{ part.type_name }}
                                    </n-tag> -->
                                    <!-- (Site tag replaced above with site name) -->
                                </n-space>
                            </div>
                            <!-- Price in bottom right -->
                            <div class="part-card-price" v-if="part.price">
                                <span>{{ part.price }}</span>
                            </div>
                        </div>
                    </n-card>
                </div>

                <!-- List View -->
                <n-list v-else-if="viewMode === 'list'" bordered>
                    <n-list-item
                        v-for="part in parts"
                        :key="part.id"
                        class="part-list-item"
                        @click="selectPart(part)"
                    >
                        <!-- NEW Badge for List View -->
                        <div v-if="isNewPart(part)" class="new-badge-list">
                            <span class="new-text">NEW!</span>
                        </div>

                        <template #prefix>
                            <n-avatar
                                v-if="part.image_base64"
                                :src="`data:image/jpeg;base64,${part.image_base64}`"
                                :size="80"
                                object-fit="cover"
                            />
                            <n-avatar v-else :size="80">
                                <n-icon size="32">
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <rect
                                            x="3"
                                            y="3"
                                            width="18"
                                            height="18"
                                            rx="2"
                                            ry="2"
                                        ></rect>
                                        <circle
                                            cx="8.5"
                                            cy="8.5"
                                            r="1.5"
                                        ></circle>
                                        <polyline
                                            points="21 15 16 10 5 21"
                                        ></polyline>
                                    </svg>
                                </n-icon>
                            </n-avatar>
                        </template>
                        <n-thing>
                            <template #header>
                                <n-text strong>{{ part.name }}</n-text>
                            </template>
                            <template #header-extra>
                                <n-space :size="8">
                                    <!-- Site Name Tag -->
                                    <n-tag size="small" type="default">
                                        {{ getSiteName(part.site_id) }}
                                    </n-tag>
                                    <!-- <n-tag size="small" type="info">
                                        {{ part.type_name }}
                                    </n-tag> -->
                                    <n-tag
                                        v-if="part.price"
                                        size="small"
                                        type="success"
                                        strong
                                    >
                                        {{ part.price }}
                                    </n-tag>
                                </n-space>
                            </template>
                            <template #description>
                                <n-ellipsis :line-clamp="2">
                                    {{ part.description }}
                                </n-ellipsis>
                            </template>
                            <n-space :size="8" style="margin-top: 8px">
                                <n-text depth="3" style="font-size: 12px">
                                    Part ID: {{ part.part_id }}
                                </n-text>
                                <n-divider vertical />
                                <n-text depth="3" style="font-size: 12px">
                                    Last seen:
                                    {{ formatDate(part.last_seen) }}
                                </n-text>
                            </n-space>
                        </n-thing>
                        <template #suffix>
                            <n-button
                                text
                                tag="a"
                                :href="part.url"
                                target="_blank"
                                @click.stop
                            >
                                View Source
                            </n-button>
                        </template>
                    </n-list-item>
                </n-list>

                <!-- Pagination -->
                <n-card v-if="totalItems > 0">
                    <n-space align="center" justify="space-between">
                        <n-select
                            v-model:value="pageSize"
                            :options="pageSizeOptions"
                            @update:value="handlePageSizeChange"
                            class="page-size-select"
                        />
                        <n-pagination
                            v-model:page="currentPage"
                            :page-count="totalPages"
                            :page-size="pageSize"
                            :item-count="totalItems"
                            @update:page="handlePageChange"
                        />
                        <n-text>
                            {{ (currentPage - 1) * pageSize + 1 }} -
                            {{ Math.min(currentPage * pageSize, totalItems) }}
                            of {{ totalItems }}
                        </n-text>
                    </n-space>
                </n-card>
            </n-space>

            <!-- Part Details Drawer -->
            <n-drawer
                v-model:show="showDetailsDrawer"
                :width="600"
                placement="right"
            >
                <n-drawer-content
                    v-if="selectedPart"
                    :title="selectedPart.name"
                >
                    <n-space vertical :size="20">
                        <!-- Image -->
                        <div v-if="selectedPart.image_base64">
                            <n-image
                                :src="`data:image/jpeg;base64,${selectedPart.image_base64}`"
                                object-fit="contain"
                                style="width: 100%"
                            />
                        </div>

                        <!-- Description -->
                        <n-card title="Description" size="small">
                            <n-text>{{ selectedPart.description }}</n-text>
                        </n-card>

                        <!-- Details -->
                        <n-card title="Details" size="small">
                            <n-descriptions :column="1" bordered>
                                <n-descriptions-item
                                    v-if="selectedPart.price"
                                    label="Price"
                                >
                                    <n-tag type="success" strong>
                                        {{ selectedPart.price }}
                                    </n-tag>
                                </n-descriptions-item>
                                <n-descriptions-item label="Part ID">
                                    {{ selectedPart.part_id }}
                                </n-descriptions-item>
                                <!-- <n-descriptions-item label="Type">
                                    {{ selectedPart.type_name }}
                                </n-descriptions-item> -->
                                <n-descriptions-item label="Site">
                                    Site {{ selectedPart.site_id }}
                                </n-descriptions-item>
                                <n-descriptions-item label="Added">
                                    {{ formatDate(selectedPart.creation_date) }}
                                </n-descriptions-item>
                            </n-descriptions>
                        </n-card>

                        <!-- Actions -->
                        <n-space>
                            <n-button
                                type="primary"
                                tag="a"
                                :href="selectedPart.url"
                                target="_blank"
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
                                                d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"
                                            ></path>
                                            <polyline
                                                points="15 3 21 3 21 9"
                                            ></polyline>
                                            <line
                                                x1="10"
                                                y1="14"
                                                x2="21"
                                                y2="3"
                                            ></line>
                                        </svg>
                                    </n-icon>
                                </template>
                                View on Source Site
                            </n-button>
                        </n-space>
                    </n-space>
                </n-drawer-content>
            </n-drawer>
        </div>
    </n-config-provider>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from "vue";
import {
    NConfigProvider,
    NSpace,
    NCard,
    NInput,
    NSelect,
    NButton,
    NButtonGroup,
    NIcon,
    NTag,
    NText,
    NStatistic,
    NSpin,
    NEmpty,
    NImage,
    NEllipsis,
    NList,
    NListItem,
    NAvatar,
    NThing,
    NDivider,
    NPagination,
    NDrawer,
    NDrawerContent,
    NDescriptions,
    NDescriptionsItem,
    NCheckbox,
    useMessage,
} from "naive-ui";
import axios from "axios";

export default defineComponent({
    name: "Browse",
    components: {
        NConfigProvider,
        NSpace,
        NCard,
        NInput,
        NSelect,
        NButton,
        NButtonGroup,
        NIcon,
        NTag,
        NText,
        NStatistic,
        NSpin,
        NEmpty,
        NImage,
        NEllipsis,
        NList,
        NListItem,
        NAvatar,
        NThing,
        NDivider,
        NPagination,
        NDrawer,
        NDrawerContent,
        NDescriptions,
        NDescriptionsItem,
        NCheckbox,
    },
    setup() {
        const message = useMessage();
        const parts = ref([]);
        const sites = ref([]);
        const loading = ref(false);
        const searchQuery = ref("");
        const viewMode = ref("grid");
        const currentPage = ref(1);
        const pageSize = ref(24);
        const totalItems = ref(0);
        const showDetailsDrawer = ref(false);
        const selectedPart = ref(null);

        // Page size options
        const pageSizeOptions = [
            { label: "24 per page", value: 24 },
            { label: "48 per page", value: 48 },
            { label: "96 per page", value: 96 },
        ];

        const filters = ref({
            siteId: null,
            // typeName: null,
            showOnlyNew: false,
        });

        const sortBy = ref("creation_date_asc");

        const themeOverrides = {
            common: {
                primaryColor: "#18a058",
            },
        };

        // Computed: Site options for filter
        const siteOptions = computed(() => {
            return [
                { label: "All Sites", value: null },
                ...sites.value.map((site) => ({
                    label: site.name,
                    value: site.id,
                })),
            ];
        });

        // Computed: Type options for filter
        const typeOptions = computed(() => {
            const types = new Set();
            parts.value.forEach((part) => {
                // if (part.type_name) {
                //     types.add(part.type_name);
                // }
            });
            return [
                { label: "All Types", value: null },
                ...Array.from(types)
                    .sort()
                    .map((type) => ({
                        label: type,
                        value: type,
                    })),
            ];
        });

        // Sort options
        const sortOptions = [
            { label: "Creation date newest", value: "creation_date_desc" },
            { label: "Creation date oldest", value: "creation_date_asc" },
            { label: "Name (A-Z)", value: "name_asc" },
            { label: "Name (Z-A)", value: "name_desc" },
            { label: "Newest First", value: "newest" },
            { label: "Oldest First", value: "oldest" },
            { label: "Recently Seen", value: "recent_seen" },
        ];

        // Check if there are active filters
        const hasActiveFilters = computed(() => {
            return (
                filters.value.siteId !== null ||
                // filters.value.typeName !== null ||
                filters.value.showOnlyNew ||
                searchQuery.value.length > 0
            );
        });

        // Total pages
        const totalPages = computed(() => {
            return Math.ceil(totalItems.value / pageSize.value);
        });

        // Handle page change
        const handlePageChange = (page) => {
            currentPage.value = page;
            loadParts();
        };

        // Load parts from API
        const loadParts = async () => {
            loading.value = true;
            try {
                const offset = (currentPage.value - 1) * pageSize.value;
                console.log("Making API request with params:", {
                    currentPage: currentPage.value,
                    pageSize: pageSize.value,
                    calculatedOffset: offset,
                });
                const params = {
                    limit: pageSize.value,
                    offset: offset,
                    site_id: filters.value.siteId,
                    // type_name: filters.value.typeName,
                    search: searchQuery.value || undefined,
                    sort: sortBy.value,
                    sort_desc: sortBy.value.endsWith("_desc"),
                    newer_than_hours: filters.value.showOnlyNew
                        ? 72
                        : undefined,
                };
                console.log("Request params:", params);
                const response = await axios.get("/api/parts", { params });
                console.log("API Response:", {
                    receivedItems: response.data.data.length,
                    firstItemId: response.data.data[0]?.id,
                    lastItemId:
                        response.data.data[response.data.data.length - 1]?.id,
                });
                parts.value = response.data.data || [];
                totalItems.value = response.data.total || 0;
                console.log("Pagination Debug:", {
                    totalItems: totalItems.value,
                    pageSize: pageSize.value,
                    totalPages: Math.ceil(totalItems.value / pageSize.value),
                    currentPage: currentPage.value,
                    responseTotal: response.data.total,
                });
            } catch (error) {
                console.error("Error loading parts:", error);
                message.error("Failed to load parts");
            } finally {
                loading.value = false;
            }
        };

        // Load sites
        const loadSites = async () => {
            try {
                const response = await axios.get("/api/sites");
                sites.value = response.data.data || [];
            } catch (error) {
                console.error("Error loading sites:", error);
            }
        };

        // Get site name by ID
        const getSiteName = (siteId) => {
            const site = sites.value.find((s) => s.id === siteId);
            return site ? site.name : `Site ${siteId}`;
        };

        // Format date
        const formatDate = (dateString) => {
            if (!dateString) return "N/A";
            const date = new Date(dateString);
            return date.toLocaleString();
        };

        // Check if part is new (within 72 hours)
        const isNewPart = (part) => {
            if (!part || !part.creation_date) return false;
            const created = new Date(part.creation_date);
            const now = new Date();
            const hoursDiff = (now - created) / (1000 * 60 * 60);
            return hoursDiff <= 72;
        };

        // Apply filters
        const applyFilters = () => {
            currentPage.value = 1; // Reset to first page
            loadParts(); // Reload with new filters
        };

        // Handle page size change
        const handlePageSizeChange = (newSize) => {
            pageSize.value = newSize;
            currentPage.value = 1; // Reset to first page
            loadParts();
        };

        // Reset filters
        const resetFilters = () => {
            searchQuery.value = "";
            filters.value = {
                siteId: null,
                // typeName: null,
                showOnlyNew: false,
            };
            sortBy.value = "newest";
            currentPage.value = 1;
        };

        // Debounced search
        let searchTimeout = null;
        const debouncedSearch = () => {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                applyFilters();
            }, 300);
        };

        // Select part to view details
        const selectPart = (part) => {
            selectedPart.value = part;
            showDetailsDrawer.value = true;
        };

        // Load data on mount
        onMounted(() => {
            loadParts();
            loadSites();
        });

        return {
            parts,
            sites,
            loading,
            searchQuery,
            viewMode,
            currentPage,
            pageSize,
            filters,
            sortBy,
            siteOptions,
            typeOptions,
            sortOptions,
            hasActiveFilters,
            pageSizeOptions,
            totalItems,
            handlePageSizeChange,
            showDetailsDrawer,
            selectedPart,
            themeOverrides,
            loadParts,
            getSiteName,
            formatDate,
            applyFilters,
            resetFilters,
            debouncedSearch,
            selectPart,
            isNewPart,
            handlePageChange,
        };
    },
});
</script>

<style scoped>
.browse-page {
    padding: 24px;
    max-width: 1600px;
    margin: 0 auto;
}

.parts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
}

.part-card-price {
    position: absolute;
    right: 16px;
    bottom: 16px;
    background: #fff;
    color: #18a058;
    font-size: 1.2rem;
    font-weight: bold;
}

.part-card {
    position: relative;
    cursor: pointer;
    transition: all 0.3s ease;
    overflow: visible;
}

.part-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.part-card-content {
    display: flex;
    flex-direction: column;
}

.part-image {
    width: 100%;
    height: 200px;
    margin-bottom: 12px;
    border-radius: 4px;
    overflow: hidden;
    background-color: #f5f5f5;
}

/* Blurred background for image area */
.part-image-blur-bg {
    width: 100%;
    height: 200px;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    overflow: hidden;
    background-size: cover;
    background-position: center;
}
.part-image-blur-bg::before {
    content: "";
    position: absolute;
    inset: 0;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    filter: blur(18px) brightness(1.15) saturate(1.2);
    z-index: 1;
    background-image: inherit;
    /* fallback for browsers that don't support inherit, will be overridden inline */
}
.part-image-centered {
    position: relative;
    z-index: 2;
    height: 100%;
    width: auto;
    max-width: 100%;
    object-fit: contain;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    background: transparent;
}

.no-image {
    width: 100%;
    height: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    background-color: #f5f5f5;
}

.part-info {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.part-list-item {
    cursor: pointer;
    transition: background-color 0.3s ease;
}

.part-list-item {
    position: relative;
}

.part-list-item:hover {
    background-color: #f5f5f5;
}

/* 90s Style NEW Badge */
.new-badge {
    position: absolute;
    top: -8px;
    right: -8px;
    z-index: 10;
    background: linear-gradient(135deg, #ff0000 0%, #ff6b6b 50%, #ff0000 100%);
    padding: 4px 12px;
    transform: rotate(12deg);
    box-shadow:
        0 4px 8px rgba(0, 0, 0, 0.3),
        0 0 0 3px #ffeb3b,
        0 0 0 6px #ff0000,
        inset 0 2px 0 rgba(255, 255, 255, 0.3),
        inset 0 -2px 0 rgba(0, 0, 0, 0.3);
    border-radius: 4px;
    animation:
        pulse 1.5s ease-in-out infinite,
        wiggle 0.5s ease-in-out;
}

.new-badge::before {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    width: 140%;
    height: 140%;
    background: radial-gradient(
        ellipse at center,
        rgba(255, 235, 59, 0.8) 0%,
        rgba(255, 235, 59, 0.4) 30%,
        transparent 70%
    );
    transform: translate(-50%, -50%) rotate(-12deg);
    z-index: -1;
    pointer-events: none;
}

.new-badge::after {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    width: 100%;
    height: 100%;
    background: repeating-linear-gradient(
        45deg,
        transparent,
        transparent 2px,
        rgba(255, 255, 255, 0.1) 2px,
        rgba(255, 255, 255, 0.1) 4px
    );
    transform: translate(-50%, -50%);
    border-radius: 4px;
    pointer-events: none;
}

.new-text {
    font-family: "Arial Black", "Arial Bold", sans-serif;
    font-size: 14px;
    font-weight: 900;
    color: #fff;
    text-shadow:
        2px 2px 0 #000,
        -1px -1px 0 #000,
        1px -1px 0 #000,
        -1px 1px 0 #000,
        0px 2px 4px rgba(0, 0, 0, 0.8);
    letter-spacing: 1px;
    position: relative;
    z-index: 1;
}

.pagination-controls {
    margin-top: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    flex-wrap: wrap;
}

.page-size-select {
    min-width: 120px;
}

.pagination-info {
    color: var(--n-text-color);
    font-size: 0.9rem;
}

.new-badge-list {
    position: absolute;
    top: 8px;
    left: 8px;
    z-index: 10;
    background: linear-gradient(135deg, #ff0000 0%, #ff6b6b 50%, #ff0000 100%);
    padding: 3px 10px;
    transform: rotate(-8deg);
    box-shadow:
        0 3px 6px rgba(0, 0, 0, 0.3),
        0 0 0 2px #ffeb3b,
        0 0 0 4px #ff0000;
    border-radius: 3px;
    animation: pulse 1.5s ease-in-out infinite;
}

.new-badge-list .new-text {
    font-size: 11px;
}

@keyframes pulse {
    0%,
    100% {
        transform: rotate(12deg) scale(1);
    }
    50% {
        transform: rotate(12deg) scale(1.05);
    }
}

@keyframes wiggle {
    0% {
        transform: rotate(12deg) scale(0.8);
        opacity: 0;
    }
    50% {
        transform: rotate(15deg) scale(1.1);
    }
    100% {
        transform: rotate(12deg) scale(1);
        opacity: 1;
    }
}

.only-new-checkbox {
    display: flex;
    align-items: center;
    height: 34px; /* Match height of other controls */
}

.only-new-checkbox :deep(.n-checkbox-box) {
    margin-top: 0;
}

.only-new-checkbox :deep(.n-checkbox__label) {
    display: flex;
    align-items: center;
    height: 100%;
}

@media (max-width: 768px) {
    .browse-page {
        padding: 12px;
    }

    .parts-grid {
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 12px;
    }

    .part-image {
        height: 150px;
    }

    .no-image {
        height: 150px;
    }

    .new-badge {
        padding: 3px 8px;
        font-size: 10px;
    }

    .new-text {
        font-size: 11px;
    }
}
</style>
