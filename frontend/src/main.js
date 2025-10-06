import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";

// Import views
import Home from "./views/Home.vue";
import Sites from "./views/Sites.vue";

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
  ],
});

// Create and mount app
const app = createApp(App);
app.use(router);
app.mount("#app");
