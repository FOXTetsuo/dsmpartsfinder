import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";

// Import views
import Home from "./views/Home.vue";
import Parts from "./views/Parts.vue";
import Browse from "./views/Browse.vue";

// Create routes array
const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/browse",
    name: "Browse",
    component: Browse,
  },
];

// Conditionally add Parts route in development
if (!import.meta.env.PROD) {
  routes.push({
    path: "/parts",
    name: "Parts",
    component: Parts,
  });
}

// Create router
const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Create and mount app
const app = createApp(App);
app.use(router);
app.mount("#app");
