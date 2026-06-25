import { describe, it, expect, beforeEach } from "vitest"
import { mount } from "@vue/test-utils"
import { setActivePinia, createPinia } from "pinia"
import { useAuthStore } from "../../stores/auth"
import Navbar from "../Navbar.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("Navbar", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it("shows Entrar button when not authenticated", () => {
    const wrapper = mount(Navbar, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Entrar")
  })

  it("shows Dashboard/Perfil/Assinatura links when authenticated", () => {
    const auth = useAuthStore()
    auth.accessToken = "fake-token"
    auth.user = {
      id: "1",
      name: "Test",
      email: "test@test.com",
      role: "USER",
      subscriptionStatus: "ACTIVE",
      subscriptionIsTrial: false,
    }

    const wrapper = mount(Navbar, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Dashboard")
    expect(wrapper.text()).toContain("Perfil")
    expect(wrapper.text()).toContain("Assinatura")
  })

  it("shows Admin link only when user role is ADMIN", () => {
    const auth = useAuthStore()
    auth.accessToken = "fake-token"
    auth.user = {
      id: "1",
      name: "Admin",
      email: "admin@test.com",
      role: "ADMIN",
      subscriptionStatus: "ACTIVE",
      subscriptionIsTrial: false,
    }

    const wrapper = mount(Navbar, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Admin")
  })

  it("shows Sair button when authenticated", () => {
    const auth = useAuthStore()
    auth.accessToken = "fake-token"
    auth.user = {
      id: "1",
      name: "Test",
      email: "test@test.com",
      role: "USER",
      subscriptionStatus: "ACTIVE",
      subscriptionIsTrial: false,
    }

    const wrapper = mount(Navbar, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Sair")
  })
})
