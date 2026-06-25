import { describe, it, expect } from "vitest"
import { mount } from "@vue/test-utils"
import Home from "../Home.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("Home", () => {
  it("renders the home page content", () => {
    const wrapper = mount(Home, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
          RaffleCard: true,
        },
      },
    })
    expect(wrapper.text()).toContain("Rifa Online")
    expect(wrapper.text()).toContain("Participe de rifas e concorra a prêmios incríveis!")
    expect(wrapper.text()).toContain("Rifas Ativas")
  })
})
