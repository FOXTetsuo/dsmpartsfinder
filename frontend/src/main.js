import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";

// Import views
import Home from "./views/Home.vue";
import Parts from "./views/Parts.vue";
import Browse from "./views/Browse.vue";

// Create router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Home,
    },
    {
      path: "/parts",
      name: "Parts",
      component: Parts,
    },
    {
      path: "/browse",
      name: "Browse",
      component: Browse,
    },
  ],
});

// Create and mount app
const app = createApp(App);
app.use(router);
app.mount("#app");
