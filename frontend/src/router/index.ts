import { createRouter, createWebHashHistory } from "vue-router"
import { userAccountStore } from "../store/account"

const routes = [
  { path: "/", name: "Message", component: () => import("@/views/Message/Message.vue") },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
router.beforeEach((to, from, next) => {
  const store = userAccountStore()
  if (store.user != undefined) {
    // login
    if (to.path === "/login") {
      next({ path: "/" })
    } else {
      next()
    }
  } else {
    if (to.path === "/login") {
      next()
    } else {
      next("/login")
    }
  }
  // next()
})
export default router
