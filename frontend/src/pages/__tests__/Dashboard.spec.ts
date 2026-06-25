import { describe, it, expect, beforeEach } from "vitest"
import { mount } from "@vue/test-utils"
import { setActivePinia, createPinia } from "pinia"
import { useAuthStore } from "../../stores/auth"
import Dashboard from "../Dashboard.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("Dashboard", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    const auth = useAuthStore()
    auth.accessToken = "fake-token"
    auth.user = {
      id: "1",
      name: "Test User",
      email: "test@example.com",
      role: "USER",
      subscriptionStatus: "ACTIVE",
      subscriptionIsTrial: false,
    }
  })

  it("renders dashboard", () => {
    const wrapper = mount(Dashboard, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Dashboard")
    expect(wrapper.text()).toContain("Gerencie suas rifas")
    expect(wrapper.text()).toContain("Criar Rifa")
  })
})
