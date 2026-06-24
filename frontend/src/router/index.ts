import { createRouter, createWebHistory } from "vue-router"
import Home from "../pages/Home.vue"

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: Home,
    },
    {
      path: "/raffles/:id",
      name: "raffle-detail",
      component: () => import("../pages/RaffleDetail.vue"),
    },
    {
      path: "/raffles/:id/checkout",
      name: "checkout",
      component: () => import("../pages/Checkout.vue"),
    },
    {
      path: "/raffles/:id/result",
      name: "raffle-result",
      component: () => import("../pages/RaffleResult.vue"),
    },
    {
      path: "/payment/success",
      name: "payment-success",
      component: () => import("../pages/PaymentSuccess.vue"),
    },
    {
      path: "/payment/pending",
      name: "payment-pending",
      component: () => import("../pages/PaymentPending.vue"),
    },
    {
      path: "/login",
      name: "login",
      component: () => import("../pages/Login.vue"),
    },
    {
      path: "/register",
      name: "register",
      component: () => import("../pages/Register.vue"),
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: () => import("../pages/Dashboard.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/dashboard/raffles/new",
      name: "create-raffle",
      component: () => import("../pages/CreateRaffle.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/dashboard/raffles/:id/edit",
      name: "edit-raffle",
      component: () => import("../pages/EditRaffle.vue"),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem("accessToken")
  if (to.meta.requiresAuth && !token) {
    next({ name: "login" })
  } else {
    next()
  }
})

export default router
