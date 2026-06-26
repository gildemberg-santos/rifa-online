import { createRouter, createWebHistory } from "vue-router"
import Home from "../pages/Home.vue"

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: Home,
      meta: { requiresAuth: true },
    },
    {
      path: "/raffles/:id",
      name: "raffle-detail",
      component: () => import("../pages/RaffleDetail.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/raffles/:id/checkout",
      name: "checkout",
      component: () => import("../pages/Checkout.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/raffles/:id/result",
      name: "raffle-result",
      component: () => import("../pages/RaffleResult.vue"),
      meta: { requiresAuth: true },
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
    {
      path: "/profile",
      name: "profile",
      component: () => import("../pages/Profile.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/subscription",
      name: "subscription",
      component: () => import("../pages/Subscription.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/my-purchases",
      name: "my-purchases",
      component: () => import("../pages/MyPurchases.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/admin",
      name: "admin",
      component: () => import("../pages/Admin.vue"),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: "/termos-de-uso",
      name: "terms-of-use",
      component: () => import("../pages/legal/TermsOfUse.vue"),
    },
    {
      path: "/politica-de-privacidade",
      name: "privacy-policy",
      component: () => import("../pages/legal/PrivacyPolicy.vue"),
    },
    {
      path: "/politica-de-cookies",
      name: "cookie-policy",
      component: () => import("../pages/legal/CookiePolicy.vue"),
    },
    {
      path: "/termo-do-organizador",
      name: "organizer-terms",
      component: () => import("../pages/legal/OrganizerTerms.vue"),
    },
    {
      path: "/contato",
      name: "contact",
      component: () => import("../pages/Contact.vue"),
    },
    {
      path: "/ajuda",
      name: "help",
      component: () => import("../pages/Help.vue"),
    },
    {
      path: "/ajuda/:slug",
      name: "help-article",
      component: () => import("../pages/HelpArticle.vue"),
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem("accessToken")
  if (to.meta.requiresAuth && !token) {
    next({ name: "login" })
    return
  }
  if (to.meta.requiresAdmin) {
    const userStr = localStorage.getItem("user")
    if (!userStr) {
      next({ name: "home" })
      return
    }
    try {
      const user = JSON.parse(userStr)
      if (user.role !== "ADMIN") {
        next({ name: "home" })
        return
      }
    } catch {
      next({ name: "home" })
      return
    }
  }
  next()
})

export default router
