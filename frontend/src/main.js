import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";

// Import views
import Home from "./views/Home.vue";
import Sites from "./views/Sites.vue";
import Parts from "./views/Parts.vue";

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
      path: "/sites",
      name: "Sites",
      component: Sites,
    },
    {
      path: "/parts",
      name: "Parts",
      component: Parts,
    },
  ],
});

// Create and mount app
const app = createApp(App);
app.use(router);
app.mount("#app");
